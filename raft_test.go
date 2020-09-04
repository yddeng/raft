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

func TestRaft_RequestVote(t *testing.T) {
	peers := map[string]string{
		"1": "127.0.0.1:6000",
		"2": "127.0.0.1:6001",
		"3": "127.0.0.1:6002",
	}

	for name := range peers {
		go func(name string) {
			if name == "3" {
				time.Sleep(2 * time.Minute) // 第3个node迟到加入
			}
			rf := Make(name, peers)
			for {
				time.Sleep(time.Second)
				rf.Submit(&testCommand{Msg: fmt.Sprintf("msg %s command", name)})
			}
		}(name)
	}

	time.Sleep(1 * time.Hour)
}

/*
QuestSyncToC isAll:true
Quests:<questID:10101 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:10102 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:10103 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:10104 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:10105 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:10131 questType:1 State:Running QuestConditions:<CondType:5 DoneTimes:0 Complete:false > >
Quests:<questID:10199 questType:1 State:Running QuestConditions:<CondType:4 DoneTimes:0 Complete:false > QuestConditions:<CondType:4 DoneTimes:0 Complete:false > QuestConditions:<CondType:4 DoneTimes:0 Complete:false > QuestConditions:<CondType:4 DoneTimes:0 Complete:false > QuestConditions:<CondType:4 DoneTimes:0 Complete:false > >
Quests:<questID:20102 questType:2 State:Running QuestConditions:<CondType:11 DoneTimes:0 Complete:false > >
Quests:<questID:20103 questType:2 State:Running QuestConditions:<CondType:10 DoneTimes:0 Complete:false > >
Quests:<questID:20104 questType:2 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:20101 questType:2 State:Running QuestConditions:<CondType:11 DoneTimes:0 Complete:false > >
Quests:<questID:30001 questType:3 State:Running QuestConditions:<CondType:10 DoneTimes:0 Complete:false > >
Quests:<questID:601 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:602 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:603 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:604 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:605 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >

QuestSyncToC isAll:true
Quests:<questID:30001 questType:3 State:Running QuestConditions:<CondType:10 DoneTimes:0 Complete:false > >
Quests:<questID:601 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:602 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:603 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:604 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:605 questType:6 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:10103 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:10104 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:10105 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:10131 questType:1 State:Running QuestConditions:<CondType:5 DoneTimes:0 Complete:false > >
Quests:<questID:10199 questType:1 State:Running QuestConditions:<CondType:4 DoneTimes:0 Complete:false > QuestConditions:<CondType:4 DoneTimes:0 Complete:false > QuestConditions:<CondType:4 DoneTimes:0 Complete:false > QuestConditions:<CondType:4 DoneTimes:0 Complete:false > QuestConditions:<CondType:4 DoneTimes:0 Complete:false > >
Quests:<questID:10101 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:10102 questType:1 State:Running QuestConditions:<CondType:1 DoneTimes:0 Complete:false > >
Quests:<questID:20102 questType:2 State:Running QuestConditions:<CondType:11 DoneTimes:0 Complete:false > >
Quests:<questID:20103 questType:2 State:Running QuestConditions:<CondType:10 DoneTimes:0 Complete:false > >
Quests:<questID:20104 questType:2 State:Running QuestConditions:<CondType:12 DoneTimes:0 Complete:false > >
Quests:<questID:20101 questType:2 State:Running QuestConditions:<CondType:11 DoneTimes:0 Complete:false > >
*/
