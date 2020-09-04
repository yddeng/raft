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

func (rf *Raft) AppendEntries(req *AppendEntriesReq, resp *AppendEntriesResp) (err error) {
	rf.Lock()
	defer rf.Unlock()
	dlog("raft %s AppendEntries %v", rf.name, req)

	if req.GetTerm() > rf.currentTerm {
		// 如果接收到的 RPC 请求或响应中，任期号T > currentTerm，那么就令 currentTerm 等于 T，并切换状态为跟随者（5.1 节）
		rf.beFollower(req.Term)
	}

	resp.Term = rf.currentTerm
	resp.Success = true

	if req.GetTerm() < rf.currentTerm {
		// 如果 term < currentTerm 就返回 false （5.1 节）
		resp.Success = false
	} else {
		if rf.state != Follower {
			rf.beFollower(req.GetTerm())
		}
		sendCh(rf.appendLogCh)

		if req.GetPrevLogIndex() != -1 {
			if int64(len(rf.log)) <= req.GetPrevLogIndex() {
				resp.Success = false
				return
			}
			prevLog := rf.log[req.GetPrevLogIndex()]
			if prevLog == nil || prevLog.Term != req.GetPrevLogTerm() {
				//如果日志在 prevLogIndex 位置处的日志条目的任期号和 prevLogTerm 不匹配，则返回 false （5.3 节）
				resp.Success = false
			}

		}
		if req.GetLeaderCommit() > rf.commitIndex {
			//如果 leaderCommit > commitIndex，令 commitIndex 等于 leaderCommit 和 新日志条目索引值中较小的一个
			rf.commitIndex = min(req.GetLeaderCommit(), rf.getLastLogIndex())
		}
	}

	/* else if rf.log[req.GetPrevLogIndex()].Term != req.GetPrevLogTerm() {
		//如果日志在 prevLogIndex 位置处的日志条目的任期号和 prevLogTerm 不匹配，则返回 false （5.3 节）

	} else if  {
		// 如果已经存在的日志条目和新的产生冲突（索引值相同但是任期号不同），删除这一条和之后所有的 （5.3 节）

	}else if {
		// 附加日志中尚未存在的任何新条目
	}else if req.LeaderCommit > rf.commitIndex{
		//如果 leaderCommit > commitIndex，令 commitIndex 等于 leaderCommit 和 新日志条目索引值中较小的一个
		rf.commitIndex = min(req.GetLeaderCommit(),rf.getLastLogIndex())
	}*/

	return
}

func (rf *Raft) sendAppendEntries(name string, req *AppendEntriesReq, resp *AppendEntriesResp) error {
	return rf.peers[name].call("Raft.AppendEntries", req, resp)
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func (rf *Raft) RequestVote(req *RequestVoteReq, resp *RequestVoteResp) (err error) {
	rf.Lock()
	defer rf.Unlock()

	if req.GetTerm() > rf.currentTerm {
		// 如果接收到的 RPC 请求或响应中，任期号T > currentTerm，那么就令 currentTerm 等于 T，并切换状态为跟随者（5.1 节）
		rf.beFollower(req.Term)
	}

	resp.Term = rf.currentTerm
	resp.VoteGranted = true

	if req.GetTerm() < rf.currentTerm {
		// 如果term < currentTerm返回 false （5.2 节）
		resp.VoteGranted = false

	} else {
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
	}

	return
}

func (rf *Raft) sendRequestVote(name string, req *RequestVoteReq, resp *RequestVoteResp) error {
	return rf.peers[name].call("Raft.RequestVote", req, resp)
}
