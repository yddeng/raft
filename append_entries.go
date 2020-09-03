package raft

func (rf *Raft) AppendEntries(req *AppendEntriesReq, resp *AppendEntriesResp) {

}

func (rf *Raft) sendAppendEntries(name string, req *AppendEntriesReq, resp *AppendEntriesResp) error {
	return rf.peers[name].rpcClient.Call("Raft.AppendEntries", req, resp)
}
