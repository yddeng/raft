package raft

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

var sString = map[State]string{
	Follower:  "follower",
	Candidate: "candidate",
	Leader:    "leader",
}

func (this State) String() string {
	return sString[this]
}

const NULL string = "null"

type ClientEnd struct {
	addr      string
	rpcClient *rpc.Client
}

type ApplyMsg struct {
	CommandValid bool // true为log，false为snapshot

	// 向application层提交日志
	Command      interface{}
	CommandIndex int
	CommandTerm  int

	// 向application层安装快照
	Snapshot          []byte
	LastIncludedIndex int
	LastIncludedTerm  int
}

type Raft struct {
	sync.Mutex
	peers map[string]*ClientEnd
	name  string

	state State

	// 所有服务器，持久化状态（lab-2A不要求持久化）
	currentTerm uint64      // 服务器最后知道的任期号(从0开始）
	votedFor    string      // 当前任期内收到选票的候选人ID (没有为NULL)
	log         []*LogEntry // 操作日志
	storage     Storage

	// 所有服务器，易失状态
	commitIndex uint64 // 已知的最大已提交日志索引(从0开始）
	lastApplied uint64 // 被状态机执行的最大日志索引(从0开始）

	// 仅Leader，易失状态（成为leader时重置）
	nextIndex  []int //	每个follower的log同步起点索引（初始为leader log的最后一项）
	matchIndex []int // 每个follower的log同步进度（初始为0），和nextIndex强关联

	// channel
	applyCh chan ApplyMsg // 应用层的提交队列

	// handler rpc
	voteCh      chan bool
	appendLogCh chan bool

	killCh chan bool
}

func Make(name string, peers map[string]string, applyCh chan ApplyMsg) (rf *Raft) {
	rf = new(Raft)
	rf.name = name
	rf.peers = make(map[string]*ClientEnd, len(peers))
	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = NULL
	rf.storage = NewMapStorage("state", name)
	rf.log = make([]*LogEntry, 0)
	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.nextIndex = make([]int, len(peers))
	rf.matchIndex = make([]int, len(peers))
	rf.applyCh = applyCh
	rf.voteCh = make(chan bool, 1)
	rf.appendLogCh = make(chan bool, 1)
	rf.killCh = make(chan bool, 1)

	for name, addr := range peers {
		rf.peers[name] = &ClientEnd{
			addr:      addr,
			rpcClient: nil,
		}
	}

	go rf.loop()
	return
}

func (rf *Raft) GetState() (uint64, State) {
	rf.Lock()
	defer rf.Unlock()
	return rf.currentTerm, rf.state
}

func sendCh(ch chan bool) {
	select {
	case <-ch:
	default:

	}
	ch <- true
}

// 最后的index
func (rf *Raft) getLastLogIndex() int {
	return len(rf.log) - 1
}

// 最后的term
func (rf *Raft) getLastLogTerm() uint64 {
	idx := rf.getLastLogIndex()
	if idx < 0 {
		return -1
	}
	return rf.log[idx].Term
}

func (rf *Raft) Submit(command Command) error {
	rf.Lock()
	defer rf.Unlock()

	dlog("submit %s : %v", rf.state.String(), command)
	// 只有leader才能写入
	if rf.state != Leader {
		return fmt.Errorf("state %s is not leader", rf.state.String())
	}

	log, err := newLogEntry(0, rf.currentTerm, command)
	if err != nil {
		return err
	}
	rf.log = append(rf.log, log)
	sendCh(rf.appendLogCh)
	//rf.persistToStorage()

	return nil
}

func (rf *Raft) loop() {
	defer dlog("server.loop.end")
	heartbeatTime := time.Millisecond * 100
	for {
		_, state := rf.GetState()
		dlog("server.loop.run %s", state)
		switch state {
		case Follower, Candidate:
			electionTime := time.Duration(rand.Intn(100)+300) * time.Millisecond
			select {
			case <-rf.appendLogCh:
			case <-rf.voteCh:
			case <-rf.killCh:
			case <-time.After(electionTime):
				rf.beCandidate()

			}
		case Leader:
			rf.appendLogLoop()
			time.Sleep(heartbeatTime)
		}

	}

}

func (rf *Raft) beCandidate() {
	rf.state = Candidate
	rf.currentTerm++
	rf.votedFor = rf.name
	go rf.electionLoop()
}

func (rf *Raft) beFollower(term uint64) {
	rf.state = Follower
	rf.votedFor = NULL
	rf.currentTerm = term
}

func (rf *Raft) beLeader() {
	if rf.state != Candidate {
		return
	}
	rf.state = Leader
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))
	// 初始化巍峨领导人上一条日志索引+1
	for i := 0; i < len(rf.nextIndex); i++ {
		rf.nextIndex[i] = rf.getLastLogIndex() + 1
	}
}

func (rf *Raft) electionLoop() {
	req := &RequestVoteReq{
		Term:          rf.currentTerm,
		LastLogIndex:  uint64(rf.getLastLogIndex()),
		LastLogTerm:   rf.getLastLogTerm(),
		CandidateName: rf.name,
	}

	// 投票展成计数
	var votes int32 = 1
	for name, _ := range rf.peers {
		if name == rf.name {
			continue
		}
		go func(name string) {
			resp := &RequestVoteResp{}
			err := rf.sendRequestVote(name, req, resp)
			if err == nil {
				if resp.Term > rf.currentTerm {
					rf.beFollower(resp.Term)
					sendCh(rf.voteCh)
					return
				}
				if rf.state != Candidate {
					return
				}
				if resp.VoteGranted {
					atomic.AddInt32(&votes, 1)
				}
				if atomic.LoadInt32(&votes) > int32(len(rf.peers)/2) {
					rf.beLeader()
					sendCh(rf.voteCh)
				}
			}
		}(name)
	}
}

func (rf *Raft) appendLogLoop() {

}
