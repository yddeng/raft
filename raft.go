package raft

import (
	"fmt"
	"math/rand"
	"sort"
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

type Event struct {
	Cmd   Command
	Index int64
	Term  int64
}

type Raft struct {
	sync.Mutex
	name  string
	peers map[string]*ClientEnd
	//leaderName string
	state State

	// 所有服务器上持久存在的
	currentTerm int64       // 服务器最后一次知道的任期号（初始化为 0，持续递增）
	votedFor    string      // 当前任期内收到选票的候选人ID (没有为NULL)
	log         []*LogEntry // 日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号

	// 所有服务器上经常变的
	commitIndex int64 // 已知的最大的已经被提交的日志条目的索引值
	lastApplied int64 // 最后被应用到状态机的日志条目索引值（初始化为 0，持续递增）

	// 在领导人里经常改变的 （选举后重新初始化）
	nextIndex  map[string]int64 // 对于每一个服务器，需要发送给他的下一个日志条目的索引值（初始化为领导人最后索引值加一）
	matchIndex map[string]int64 // 每个follower的log同步进度（初始为0），和nextIndex强关联

	// follower 与 leader 的心跳超时
	heartbeatTimer *time.Timer

	// handler rpc
	voteCh         chan struct{}
	appendLogCh    chan struct{}
	commitNotifyCh chan struct{}

	// eventCh
	eventCh chan<- Event

	killCh chan struct{}
}

func Make(name string, peers map[string]string, eventCh chan<- Event) (rf *Raft) {
	rf = new(Raft)
	rf.name = name
	rf.peers = make(map[string]*ClientEnd, len(peers))
	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = NULL
	rf.log = make([]*LogEntry, 0)
	rf.commitIndex = -1
	rf.lastApplied = -1
	rf.nextIndex = make(map[string]int64, len(peers))
	rf.matchIndex = make(map[string]int64, len(peers))
	rf.heartbeatTimer = time.NewTimer(1000 * time.Millisecond)
	rf.voteCh = make(chan struct{}, 1)
	rf.appendLogCh = make(chan struct{}, 1)
	rf.commitNotifyCh = make(chan struct{}, 16)
	rf.eventCh = eventCh
	rf.killCh = make(chan struct{}, 1)

	for name, addr := range peers {
		rf.peers[name] = &ClientEnd{
			addr: addr,
		}
	}

	// 加载
	rf.loadFromDisk()
	// 启动rpc
	go rf.startRPCServer()
	// raft loop
	go rf.loop()
	// event loop
	go rf.eventChanLoop()
	return
}

func sendCh(ch chan struct{}) {
	select {
	case <-ch:
	default:
	}
	ch <- struct{}{}
}

// 最后的index
func (rf *Raft) getLastLogIndex() int64 {
	return int64(len(rf.log) - 1)
}

// 最后的term
func (rf *Raft) getLastLogTerm() int64 {
	idx := rf.getLastLogIndex()
	if idx < 0 {
		return -1
	}
	return rf.log[idx].Term
}

// 上一条下发给 follower 日志的索引
func (rf *Raft) getPrevLogIndex(name string) int64 {
	return rf.nextIndex[name] - 1
}

// 上一条下发给 follower 日志的 term
func (rf *Raft) getPrevLogTerm(name string) int64 {
	idx := rf.getPrevLogIndex(name)
	if idx < 0 {
		return -1
	}
	return rf.log[idx].Term
}

func (rf *Raft) loop() {
	defer dlog("raft %s server.loop.end", rf.name)
	heartbeatTime := time.Millisecond * 100
	for {
		rf.Lock()
		state := rf.state
		rf.Unlock()

		dlog("raft %s server.loop.run %s", rf.name, state)
		switch state {
		case Follower, Candidate:
			electionTime := time.Duration(rand.Intn(100)+300) * time.Millisecond
			rf.heartbeatTimer.Reset(electionTime)
			select {
			case <-rf.appendLogCh:
			case <-rf.voteCh:
			case <-rf.killCh:
				return
			case <-rf.heartbeatTimer.C:
				// 如果在超过选举超时时间的情况之前没有收到当前领导人（即该领导人的任期需与这个跟随者的当前任期相同）的心跳/附加日志，或者是给某个候选人投了票，就自己变成候选人（5.2 节)
				rf.Lock()
				rf.beCandidate()
				rf.Unlock()
			}
		case Leader:
			select {
			case <-rf.killCh:
				return
			default:
			}
			rf.startAppendLog()
			time.Sleep(heartbeatTime)
		}

	}

}

/*
	在转变成候选人后就立即开始选举过程
	1.自增当前的任期号（currentTerm）
	2.给自己投票
	3.重置选举超时计时器
	4.发送请求投票的 RPC 给其他所有服务器
*/
func (rf *Raft) beCandidate() {
	dlog("raft %s beCandidate", rf.name)
	rf.state = Candidate
	rf.currentTerm++
	rf.votedFor = rf.name
	go rf.startElection()
}

func (rf *Raft) beFollower(term int64) {
	dlog("raft %s beFollower term %d", rf.name, term)
	rf.state = Follower
	rf.votedFor = NULL
	rf.currentTerm = term
}

func (rf *Raft) beLeader() {
	if rf.state != Candidate {
		return
	}
	dlog("raft %s beLeader term %d", rf.name, rf.currentTerm)
	rf.state = Leader
	rf.nextIndex = make(map[string]int64, len(rf.peers))
	rf.matchIndex = make(map[string]int64, len(rf.peers))
	// 初始化为领导人上一条日志索引+1
	for name := range rf.peers {
		rf.nextIndex[name] = rf.getLastLogIndex() + 1
	}
}

func (rf *Raft) startElection() {
	req := &RequestVoteReq{
		Term:          rf.currentTerm,
		LastLogIndex:  rf.getLastLogIndex(),
		LastLogTerm:   rf.getLastLogTerm(),
		CandidateName: rf.name,
	}
	dlog("raft %s startElection %v ", rf.name, req)

	// 防重复触发
	var fired int32 = 0
	// 投票赞成计数
	var votes int32 = 1

	// 如果接收到大多数服务器的选票，那么就变成领导人
	for name, end := range rf.peers {
		if name == rf.name {
			continue
		}
		go func(name string, end *ClientEnd) {
			resp := &RequestVoteResp{}
			err := end.call("Raft.RequestVote", req, resp)
			if err == nil {
				dlog("raft %s->%s requestVoteResp req:%v resp:%v ", rf.name, name, req, resp)
				rf.Lock()
				defer rf.Unlock()
				if resp.GetTerm() > rf.currentTerm {
					rf.beFollower(resp.Term)
					sendCh(rf.voteCh)
					return
				}
				if rf.state != Candidate {
					return
				}
				if atomic.LoadInt32(&fired) == 0 && resp.GetVoteGranted() {
					if atomic.AddInt32(&votes, 1) >= int32(len(rf.peers)/2+1) {
						atomic.StoreInt32(&fired, 1)
						rf.beLeader()
						sendCh(rf.voteCh)
					}
				}
			}
		}(name, end)
	}
}

func (rf *Raft) startAppendLog() {
	rf.Lock()
	defer rf.Unlock()
	if rf.state != Leader {
		return
	}

	for name, end := range rf.peers {
		if name == rf.name {
			continue
		}

		idx := rf.nextIndex[name]
		req := &AppendEntriesReq{
			Term:         rf.currentTerm,
			LeaderName:   rf.name,
			PrevLogIndex: rf.getPrevLogIndex(name),
			PrevLogTerm:  rf.getPrevLogTerm(name),
			LeaderCommit: rf.commitIndex,
			Entries:      rf.log[idx:],
		}

		go func(name string, end *ClientEnd, req *AppendEntriesReq) {
			resp := &AppendEntriesResp{}
			if err := end.call("Raft.AppendEntries", req, resp); err != nil {
				dlog("sendAppendEntries %s", err)
				return
			}
			dlog("raft %s->%s appendEntriesResp req:%v resp:%v ", rf.name, name, req, resp)
			rf.Lock()
			defer rf.Unlock()

			if rf.state != Leader || req.GetTerm() != rf.currentTerm {
				return
			}
			//
			if resp.GetTerm() > rf.currentTerm {
				rf.beFollower(resp.GetTerm())
				return
			}

			if resp.GetSuccess() {
				rf.nextIndex[name] = req.GetPrevLogIndex() + int64(len(req.GetEntries())) + 1
				rf.matchIndex[name] = rf.nextIndex[name] - 1
				// 更新commitIndex，如果被大多数 follower 复制，则提交。
				// 排序每一个 follower 的 matchIdx，取位置中间值，就是最大的
				idxs := make([]int, 0, len(rf.peers))
				for name := range rf.peers {
					if name == rf.name {
						idxs = append(idxs, int(rf.getLastLogIndex()))
						continue
					}
					idxs = append(idxs, int(rf.matchIndex[name]))
				}
				sort.Ints(idxs)
				newCommitIdx := int64(idxs[len(rf.peers)/2])
				if newCommitIdx > rf.commitIndex {
					rf.commitIndex = newCommitIdx
					rf.commitNotifyCh <- struct{}{}
				}
			} else {
				// 在被跟随者拒绝之后，领导人就会减小 nextIndex 值并进行重试。最终 nextIndex 会在某个位置使得领导人和跟随者的日志达成一致。
				if resp.GetConflictTerm() != -1 {
					lastIndexOfTerm := -1
					conflictTerm := resp.GetConflictTerm()
					for i := len(rf.log) - 1; i >= 0; i-- {
						if rf.log[i].Term == conflictTerm {
							lastIndexOfTerm = i
							break
						}
					}
					if lastIndexOfTerm >= 0 {
						rf.nextIndex[name] = int64(lastIndexOfTerm) + 1
					} else {
						rf.nextIndex[name] = resp.GetConflictIndex()
					}
				} else {
					rf.nextIndex[name] = resp.GetConflictIndex()
				}
			}

		}(name, end, req)
	}
}

func (rf *Raft) eventChanLoop() {
	for {
		select {
		case <-rf.killCh:
			return
		case <-rf.commitNotifyCh:
			rf.Lock()
			savedTerm := rf.currentTerm
			savedLastApplied := rf.lastApplied
			var entries []*LogEntry
			if rf.commitIndex > rf.lastApplied {
				entries = rf.log[rf.lastApplied+1 : rf.commitIndex+1]
				rf.lastApplied = rf.commitIndex
			}
			rf.Unlock()
			dlog("commitChanLoop entries=%v, savedLastApplied=%d", entries, savedLastApplied)

			for i, entry := range entries {
				command, _ := newCommand(entry.GetCommandName(), entry.GetCommand())
				rf.eventCh <- Event{
					Cmd:   command,
					Index: savedLastApplied + int64(i) + 1,
					Term:  savedTerm,
				}
			}

		}
	}
	dlog("raft %s commitChanLoop done", rf.name)
}

func (rf *Raft) Submit(command Command) error {
	rf.Lock()
	defer rf.Unlock()

	dlog("submit %s %s : %v", rf.name, rf.state.String(), command)
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
	rf.saveToDisk()
	return nil
}

func init() {
	rand.Seed(time.Now().Unix())
}
