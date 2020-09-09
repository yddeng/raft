package raft

import "fmt"

//Helper function

func sendCh(ch chan struct{}) {
	select {
	case <-ch:
	default:
	}
	ch <- struct{}{}
}

func (rf *Raft) getLog(idx uint64) *LogEntry {
	idx = rf.index2LogPos(idx)
	if idx < 0 || uint64(len(rf.log)) <= idx {
		fmt.Println("---", idx, rf.lastIncludedIndex, len(rf.log))
		return nil
	}
	return rf.log[idx]
}

// 上一条下发给 follower 日志的索引
func (rf *Raft) getPrevLogIndex(name string) uint64 {
	return rf.nextIndex[name] - 1
}

// 上一条下发给 follower 日志的 term
func (rf *Raft) getPrevLogTerm(name string) uint64 {
	idx := rf.getPrevLogIndex(name)
	if idx <= rf.lastIncludedIndex {
		return rf.lastIncludedTerm
	}
	return rf.getLog(idx).Term
}

// 最后的index
func (rf *Raft) getLastLogIndex() uint64 {
	return uint64(len(rf.log)) + rf.lastIncludedIndex
}

// 最后的term
func (rf *Raft) getLastLogTerm() uint64 {
	idx := rf.getLastLogIndex()
	if idx <= rf.lastIncludedIndex {
		return rf.lastIncludedTerm
	}
	return rf.getLog(idx).Term
}

// 日志index转化成log数组下标
func (rf *Raft) index2LogPos(index uint64) (pos uint64) {
	return index - rf.lastIncludedIndex - 1
}
