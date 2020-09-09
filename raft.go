package raft

import (
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

type ApplyMsg struct {
	CommandValid bool // true为log，false为snapshot

	// 向应用层提交日志
	Command      Command
	CommandIndex uint64
	CommandTerm  uint64

	// 向应用层安装快照
	Snapshot []byte
}

type Raft struct {
	sync.Mutex
	name  string
	peers map[string]*ClientEnd
	//leaderName string
	state State

	// 所有服务器上持久存在的
	currentTerm       uint64      // 服务器最后一次知道的任期号（初始化为 0，持续递增）
	votedFor          string      // 当前任期内收到选票的候选人ID (没有为NULL)
	log               []*LogEntry // 日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号 // first index is 1
	lastIncludedIndex uint64      // 快照中包含的最后日志条目的索引值
	lastIncludedTerm  uint64      // 快照中包含的最后日志条目的任期号

	// 所有服务器上经常变的
	commitIndex uint64 // 已知的最大的已经被提交的日志条目的索引值 (初始化为 0，持续递增）
	lastApplied uint64 // 最后提交到应用层的日志条目索引值（初始化为 0，持续递增）

	// 在领导人里经常改变的 （选举后重新初始化）
	nextIndex  map[string]uint64 // 对于每一个服务器，需要发送给他的下一个日志条目的索引值（初始化为领导人最后索引值加一）
	matchIndex map[string]uint64 // 每个follower的log同步进度（初始为0），和nextIndex强关联

	// follower 与 leader 的心跳超时
	heartbeatTimer *time.Timer

	// handler rpc
	voteCh         chan struct{}
	appendLogCh    chan struct{}
	commitNotifyCh chan struct{}

	// applyMsgCh
	applyMsgCh chan<- ApplyMsg

	killCh chan struct{}
}

func (rf *Raft) Stop() {
	close(rf.killCh)
}

func (rf *Raft) beCandidate() {
	dlog("raft %s beCandidate", rf.name)
	rf.state = Candidate
	rf.currentTerm++
	rf.votedFor = rf.name
	rf.saveStateToDisk()
	go rf.startElection()
}

func (rf *Raft) beFollower(term uint64) {
	dlog("raft %s beFollower term %d", rf.name, term)
	rf.state = Follower
	rf.votedFor = NULL
	rf.currentTerm = term
	rf.saveStateToDisk()
}

func (rf *Raft) beLeader() {
	if rf.state != Candidate {
		return
	}
	dlog("raft %s beLeader term %d", rf.name, rf.currentTerm)
	rf.state = Leader
	rf.nextIndex = make(map[string]uint64, len(rf.peers))
	rf.matchIndex = make(map[string]uint64, len(rf.peers))
	// 初始化为领导人上一条日志索引+1
	for name := range rf.peers {
		rf.nextIndex[name] = rf.getLastLogIndex() + 1
	}
}

func (rf *Raft) loop() {
	defer dlog("raft %s server.loop.end", rf.name)
	broadcastTime := time.Millisecond * 100
	for {
		rf.Lock()
		state := rf.state
		rf.Unlock()

		dlog("raft %s server.loop.run %s", rf.name, state)
		switch state {
		case Follower, Candidate:
			electionTimeout := time.Duration(rand.Intn(100)+300) * time.Millisecond
			rf.heartbeatTimer.Reset(electionTimeout)
			select {
			case <-rf.appendLogCh:
			case <-rf.voteCh:
			case <-rf.killCh:
				return
			case <-rf.heartbeatTimer.C:
				// 如果在超过选举超时时间的情况之前没有收到当前领导人（即该领导人的任期需与这个跟随者的当前任期相同）的心跳/附加日志，
				// 或者是给某个候选人投了票，就自己变成候选人（5.2 节)
				dlog("----- heartbeat timeout")
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
			time.Sleep(broadcastTime)
		}
	}
}

// 开始选举
func (rf *Raft) startElection() {
	rf.Lock()
	defer rf.Unlock()

	req := &RequestVoteReq{
		Term:          rf.currentTerm,
		LastLogIndex:  rf.getLastLogIndex(),
		LastLogTerm:   rf.getLastLogTerm(),
		CandidateName: rf.name,
	}
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
			if err := end.call("Raft.RequestVote", req, resp); err != nil {
				dlog("sendRequestVote err:%s", err)
				return
			}

			dlog("raft %s->%s requestVoteResp req:%v resp:%v ", rf.name, name, req, resp)
			rf.Lock()
			defer rf.Unlock()
			if resp.GetTerm() > rf.currentTerm {
				rf.beFollower(resp.Term)
				//sendCh(rf.voteCh)
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

		if rf.nextIndex[name] <= rf.lastIncludedIndex {
			// 如果nextIndex在leader的snapshot内，那么直接同步snapshot
			rf.sendInstallSnapshot(name, end)
		} else {
			// 否则同步日志
			rf.sendAppendLog(name, end)
		}

	}
}

func (rf *Raft) sendAppendLog(name string, end *ClientEnd) {
	logPos := rf.index2LogPos(rf.nextIndex[name])
	req := &AppendEntriesReq{
		Term:         rf.currentTerm,
		LeaderName:   rf.name,
		PrevLogIndex: rf.getPrevLogIndex(name),
		PrevLogTerm:  rf.getPrevLogTerm(name),
		LeaderCommit: rf.commitIndex,
		Entries:      append([]*LogEntry{}, rf.log[logPos:]...),
	}

	go func() {
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
			rf.nextIndex[name] = req.GetPrevLogIndex() + uint64(len(req.GetEntries())) + 1
			rf.matchIndex[name] = rf.nextIndex[name] - 1
			// 更新commitIndex，如果被大多数 follower 复制，则提交。
			rf.updateCommitIndex()
		} else {
			// 在被跟随者拒绝之后，领导人就会减小 nextIndex 值并进行重试。最终 nextIndex 会在某个位置使得领导人和跟随者的日志达成一致。
			if resp.GetConflictTerm() != -1 {
				lastIndexOfTerm := -1
				conflictTerm := uint64(resp.GetConflictTerm())
				for i := len(rf.log) - 1; i >= 0; i-- {
					if rf.log[i].Term == conflictTerm {
						lastIndexOfTerm = i
						break
					}
				}
				if lastIndexOfTerm >= 0 {
					rf.nextIndex[name] = uint64(lastIndexOfTerm) + 1
				} else {
					rf.nextIndex[name] = uint64(resp.GetConflictIndex())
				}
			} else {
				rf.nextIndex[name] = uint64(resp.GetConflictIndex())
			}
		}

	}()
}

func (rf *Raft) sendInstallSnapshot(name string, end *ClientEnd) {
	req := &InstallSnapshotReq{
		Term:              rf.currentTerm,
		LastIncludedIndex: rf.lastIncludedIndex,
		LastIncludedTerm:  rf.lastIncludedTerm,
		LeaderName:        rf.name,
		Data:              rf.readSnapshotFromDisk(),
		Offset:            0,
		Done:              true,
	}
	go func() {
		resp := &InstallSnapshotResp{}
		if err := end.call("Raft.InstallSnapshot", req, resp); err != nil {
			dlog("sendInstallSnapshot %s", err)
			return
		}

		rf.Lock()
		defer rf.Unlock()
		if rf.state != Leader || rf.currentTerm != req.GetTerm() {
			return
		}
		if resp.GetTerm() > rf.currentTerm {
			rf.beFollower(resp.GetTerm())
			return
		}

		rf.matchIndex[name] = rf.lastIncludedIndex
		rf.nextIndex[name] = rf.lastIncludedIndex + 1
		rf.updateCommitIndex()
	}()
}

func (rf *Raft) updateCommitIndex() {
	// 排序每一个 follower 的 matchIdx，取位置中间值，就是最大的
	idxs := make([]uint64, 0, len(rf.peers))
	for name := range rf.peers {
		if name == rf.name {
			idxs = append(idxs, rf.getLastLogIndex())
			continue
		}
		idxs = append(idxs, rf.matchIndex[name])
	}
	sort.Slice(idxs, func(i, j int) bool { return idxs[i] < idxs[j] })
	newCommitIdx := idxs[len(rf.peers)/2]
	if newCommitIdx > rf.commitIndex {
		rf.commitIndex = newCommitIdx
		sendCh(rf.commitNotifyCh)
	}
}

func (rf *Raft) eventChanLoop() {
	defer dlog("raft %s commitChanLoop done", rf.name)
	for {
		select {
		case <-rf.killCh:
			return
		case <-rf.commitNotifyCh:
			rf.Lock()
			for rf.commitIndex > rf.lastApplied {
				if rf.lastApplied < rf.lastIncludedIndex {
					rf.lastApplied = rf.lastIncludedIndex
				}
				rf.lastApplied += 1
				entry := rf.getLog(rf.lastApplied)
				command, _ := newCommand(entry.GetCommandName(), entry.GetCommand())
				// 堵塞后，进程不能进行
				rf.applyMsgCh <- ApplyMsg{
					CommandValid: true,
					Command:      command,
					CommandIndex: entry.GetIndex(),
					CommandTerm:  entry.GetTerm(),
				}
				dlog("commitChanLoop entries=%v", entry)
			}
			rf.Unlock()
		}
	}
}

func (rf *Raft) Submit(command Command) (isLeader bool, err error) {
	rf.Lock()
	dlog("submit %s %s : %v", rf.name, rf.state.String(), command)
	// 只有leader才能写入
	if rf.state != Leader {
		rf.Unlock()
		return false, ErrNotLeader
	}

	log, err := newLogEntry(rf.getLastLogIndex()+1, rf.currentTerm, command)
	if err != nil {
		rf.Unlock()
		return true, err
	}
	rf.log = append(rf.log, log)
	sendCh(rf.appendLogCh)
	rf.saveStateToDisk()
	rf.Unlock()

	// 立即执行一次
	rf.startAppendLog()
	return true, nil
}

func (rf *Raft) MakeSnapshot(commitIdx uint64, snapshot []byte) {
	rf.Lock()
	defer rf.Unlock()
	dlog("makeSnapshot raft %s state %s : commitIdx %d", rf.name, rf.state.String(), commitIdx)

	if commitIdx <= rf.lastIncludedIndex {
		return
	}

	//update last included index & term
	lopPos := rf.index2LogPos(commitIdx)
	rf.lastIncludedIndex = commitIdx
	rf.lastIncludedTerm = rf.log[lopPos].Term
	rf.log = append(make([]*LogEntry, 0), rf.log[lopPos+1:]...)
	rf.saveStateToDisk()
	rf.saveSnapshotToDisk(snapshot)
}

func Make(name string, peers map[string]string, applyMsgCh chan<- ApplyMsg) (rf *Raft) {
	rf = new(Raft)
	rf.name = name
	rf.peers = make(map[string]*ClientEnd, len(peers))
	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = NULL
	rf.log = make([]*LogEntry, 0)
	rf.lastIncludedIndex = 0
	rf.lastIncludedTerm = 0
	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.nextIndex = make(map[string]uint64, len(peers))
	rf.matchIndex = make(map[string]uint64, len(peers))
	rf.heartbeatTimer = time.NewTimer(1000 * time.Millisecond)
	rf.voteCh = make(chan struct{}, 1)
	rf.appendLogCh = make(chan struct{}, 1)
	rf.commitNotifyCh = make(chan struct{}, 1)
	rf.applyMsgCh = applyMsgCh
	rf.killCh = make(chan struct{}, 1)

	for name, addr := range peers {
		rf.peers[name] = &ClientEnd{
			addr: addr,
		}
	}

	// lead state
	rf.loadStateFromDisk()
	// lead snapshot
	rf.installSnapshot()

	// 启动rpc
	go rf.startRPCServer()
	// raft loop
	go rf.loop()
	// event loop
	go rf.eventChanLoop()
	return
}

func init() {
	rand.Seed(time.Now().Unix())
}
