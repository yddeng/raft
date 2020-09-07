package raft

import (
	"fmt"
	"testing"
	"time"
)

type testCommand struct {
	Msg string
}

func (this *testCommand) CommandName() string {
	return "testCommand"
}

func TestMakeRaft(t *testing.T) {
	debug = 1
	peers := map[string]string{
		"1": "127.0.0.1:6000",
		"2": "127.0.0.1:6001",
		"3": "127.0.0.1:6002",
	}

	RegisterCommand(&testCommand{})

	for name := range peers {
		go func(name string) {
			if name == "3" {
				time.Sleep(2 * time.Minute) // 第3个node迟到加入
			}
			eventCh := make(chan Event, 10)
			rf := Make(name, peers, eventCh)
			go func() {
				i := 0
				for e := range eventCh {
					i++
					fmt.Printf("---- raft %s i->%d event %v\n", name, i, e)
				}
			}()
			i := 0
			for {
				i++
				time.Sleep(time.Second)
				rf.Submit(&testCommand{Msg: fmt.Sprintf("msg %s->%d command", name, i)})
			}
		}(name)
	}

	time.Sleep(1 * time.Hour)
}
