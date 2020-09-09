package raft

import "fmt"

var (
	ErrNotLeader = fmt.Errorf("raft is not leader")
)
