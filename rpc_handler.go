package raft

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type ClientEnd struct {
	addr      string
	rpcClient *rpc.Client
	sync.Mutex
}

func (end *ClientEnd) call(method string, args interface{}, reply interface{}) error {
	var client *rpc.Client
	var err error
	end.Lock()
	client = end.rpcClient
	end.Unlock()
	if client == nil {
		client, err = rpc.DialHTTP("tcp", end.addr)
		if err != nil {
			return err
		}
		end.Lock()
		end.rpcClient = client
		end.Unlock()
	}

	if err := client.Call(method, args, reply); err != nil {
		return err
	}
	return nil
}

func (rf *Raft) startRPCServer() {
	addr := rf.peers[rf.name].addr

	server := rpc.NewServer()
	server.Register(rf)

	var err error
	var listener net.Listener
	if listener, err = net.Listen("tcp", addr); err != nil {
		log.Fatal(err)
	}
	if err = http.Serve(listener, server); err != nil {
		log.Fatal(err)
	}
}

func (rf *Raft) RequestVote(req *RequestVoteReq, resp *RequestVoteResp) (err error) {
	rf.Lock()
	defer rf.Unlock()

	resp.Term = rf.currentTerm
	resp.VoteGranted = true

	if req.GetTerm() < rf.currentTerm {
		// 如果term < currentTerm返回 false （5.2 节）
		resp.VoteGranted = false
		return
	}

	if req.GetTerm() > rf.currentTerm {
		// 如果接收到的 RPC 请求或响应中，任期号T > currentTerm，那么就令 currentTerm 等于 T，并切换状态为跟随者（5.1 节）
		rf.beFollower(req.Term)
	}

	// 如果 votedFor 为空或者为 candidateId，并且候选人的日志至少和自己一样新，那么就投票给他（5.2 节，5.4 节）
	// Raft 通过比较两份日志中最后一条日志条目的索引值和任期号定义谁的日志比较新。
	//   1.如果两份日志最后的条目的任期号不同，那么任期号大的日志更加新。2.如果两份日志最后的条目任期号相同，那么日志比较长的那个就更加新。

	lastLogTerm := rf.getLastLogTerm()
	lastLogIndex := rf.getLastLogIndex()
	if (rf.votedFor == NULL || rf.votedFor == req.GetCandidateName()) &&
		(req.GetLastLogTerm() > lastLogTerm ||
			(req.GetLastLogTerm() == lastLogTerm && req.GetLastLogIndex() >= lastLogIndex)) {
		rf.votedFor = req.CandidateName
		sendCh(rf.voteCh)
	} else {
		resp.VoteGranted = false
	}

	return
}

func (rf *Raft) AppendEntries(req *AppendEntriesReq, resp *AppendEntriesResp) (err error) {
	rf.Lock()
	defer rf.Unlock()
	dlog("raft %s AppendEntries %v", rf.name, req)

	resp.Term = rf.currentTerm
	resp.Success = true
	resp.ConflictIndex = -1
	resp.ConflictTerm = -1

	if req.GetTerm() < rf.currentTerm {
		// 如果 term < currentTerm 就返回 false （5.1 节）
		resp.Success = false
		return
	}

	if req.GetTerm() > rf.currentTerm {
		// 如果接收到的 RPC 请求或响应中，任期号T > currentTerm，那么就令 currentTerm 等于 T，并切换状态为跟随者（5.1 节）
		rf.beFollower(req.GetTerm())
	}

	if rf.state != Follower {
		rf.beFollower(req.GetTerm())
	}
	sendCh(rf.appendLogCh)

	if req.GetPrevLogIndex() < rf.lastIncludedIndex {
		// 同步位置在快照内，直接从 0 开始同步,即同步快照过来
		resp.Success = false
		resp.ConflictIndex = 0
		return
	} else if req.GetPrevLogIndex() == rf.lastIncludedIndex {
		if req.GetPrevLogTerm() != rf.lastIncludedTerm {
			// 同步位置为快照最后索引，但是term 冲突，直接从 0 开始同步,即同步快照过来
			resp.Success = false
			resp.ConflictIndex = 0
			return
		}
	} else {
		// 日志在 prevLogIndex 位置处的日志条目的索引号和 prevLogIndex 不匹配
		if rf.getLastLogIndex() < req.GetPrevLogIndex() {
			resp.Success = false
			resp.ConflictIndex = rf.getLastLogIndex()
			return
		}

		// 日志在 prevLogIndex 位置处的日志条目的任期号和 prevLogTerm 不匹配
		if rf.getLog(req.GetPrevLogIndex()).Term != req.GetPrevLogTerm() {
			resp.Success = false
			resp.ConflictTerm = rf.getLog(req.GetPrevLogIndex()).Term
			// 找到冲突term的首次出现位置，最差就是PrevLogIndex
			for index := rf.lastIncludedIndex + 1; index <= req.GetPrevLogIndex(); index++ {
				if rf.getLog(req.GetPrevLogIndex()).Term == resp.ConflictTerm {
					resp.ConflictIndex = index
					break
				}
			}
			return
		}
	}

	// 如果已经存在的日志条目和新的产生冲突（索引值相同但是任期号不同），删除这一条和之后所有的 （5.3 节）
	// 采用极端方式，第一条就当作冲突。直接覆盖本地
	if len(req.GetEntries()) > 0 {
		index := req.GetPrevLogIndex() + 1
		logPos := rf.index2LogPos(index)
		rf.log = rf.log[:logPos]
		rf.log = append(rf.log, req.GetEntries()...)
		rf.saveStateToDisk()
	}

	if req.GetLeaderCommit() > rf.commitIndex {
		//如果 leaderCommit > commitIndex，令 commitIndex 等于 leaderCommit 和 新日志条目索引值中较小的一个
		rf.commitIndex = min(req.GetLeaderCommit(), rf.getLastLogIndex())
		sendCh(rf.commitNotifyCh)
	}

	return
}

func (rf *Raft) InstallSnapshot(req *InstallSnapshotReq, resp *InstallSnapshotResp) (err error) {
	rf.Lock()
	defer rf.Unlock()

	resp.Term = rf.currentTerm
	if req.GetTerm() < rf.currentTerm {
		// 如果term < currentTerm就立即回复
		return
	}

	if req.GetTerm() > rf.currentTerm {
		// 如果接收到的 RPC 请求或响应中，任期号T > currentTerm，那么就令 currentTerm 等于 T，并切换状态为跟随者（5.1 节）
		rf.beFollower(req.Term)
	}
	sendCh(rf.appendLogCh)

	// leader快照不如本地长，那么忽略这个快照
	if req.GetLastIncludedIndex() <= rf.lastIncludedIndex {
		return
	}

	if req.GetLastIncludedIndex() >= rf.getLastLogIndex() {
		// 快照比本地日志长，日志清空
		rf.log = make([]*LogEntry, 0)
	} else {
		// 快照外还有日志，判断是否需要截断
		curLog := rf.getLog(req.GetLastIncludedIndex())
		if curLog.Term != req.GetLastIncludedTerm() {
			// term冲突，丢弃整个日志
			rf.log = make([]*LogEntry, 0)
		} else {
			// 如果现存的日志条目与快照中最后包含的日志条目具有相同的索引值和任期号，则保留其后的日志条目
			logPos := rf.index2LogPos(req.GetLastIncludedIndex())
			rf.log = rf.log[logPos+1:]
			//logList := make([]*LogEntry, 0)
			//logList = append(logList, rf.log[logPos:]...)
			//rf.log = logList
			dlog("snapshot pos %d logs %v", logPos, rf.log)
		}
	}

	rf.lastIncludedIndex, rf.lastIncludedTerm = req.GetLastIncludedIndex(), req.GetLastIncludedTerm()
	rf.saveStateToDisk()
	rf.saveSnapshotToDisk(req.GetData())

	rf.commitIndex = max(rf.commitIndex, rf.lastIncludedIndex)
	rf.lastApplied = max(rf.lastApplied, rf.lastIncludedIndex)
	if rf.lastApplied > rf.lastIncludedIndex {
		return
	}

	rf.applyMsgCh <- ApplyMsg{
		CommandValid: false,
		Snapshot:     req.GetData(),
	}

	// 如果 done 是 false，则继续等待更多的数据
	// 保存快照文件，丢弃具有较小索引的任何现有或部分快照
	// 使用快照重置状态机（并加载快照的集群配置）
	return
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
