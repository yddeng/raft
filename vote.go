package raft

func (rf *Raft) RequestVote(req *RequestVoteReq, resp *RequestVoteResp) {

	if req.Term > rf.currentTerm {
		rf.beFollower(req.Term)
		sendCh(rf.voteCh)
	}
	success := false

	if req.Term < rf.currentTerm {

	} else if rf.votedFor != NULL && rf.votedFor != req.CandidateName {

	} else if req.LastLogTerm < rf.getLastLogTerm() {

	} else if req.LastLogTerm == rf.getLastLogTerm() && req.LastLogIndex < rf.getLastLogIndex() {

	} else {
		rf.votedFor = req.CandidateName
		success = true
		rf.state = Follower
		sendCh(rf.voteCh)
	}

	resp.Term = rf.currentTerm
	resp.VoteGranted = success
}

func (rf *Raft) sendRequestVote(name string, req *RequestVoteReq, resp *RequestVoteResp) error {
	return rf.peers[name].rpcClient.Call("Raft.RequestVote", req, resp)
}
