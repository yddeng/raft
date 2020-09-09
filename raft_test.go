package raft

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
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
	//debug = 1
	peers := map[string]string{
		"1": "127.0.0.1:6000",
		"2": "127.0.0.1:6001",
		"3": "127.0.0.1:6002",
	}

	RegisterCommand(&testCommand{})

	for name := range peers {
		go func(name string) {
			if name == "3" {
				time.Sleep(1 * time.Minute) // 第3个node迟到加入
			}
			applyMsgCh := make(chan ApplyMsg, 10)
			rf := Make(name, peers, applyMsgCh)
			go func() {
				i := 0
				for e := range applyMsgCh {
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

type cmdCalculation struct {
	Operator string
	Num      int32
}

func (this *cmdCalculation) CommandName() string {
	return "cmdCalculation"
}

func TestMakeRaft_Calculation(t *testing.T) {
	//debug = 1
	peers := map[string]string{
		"1": "127.0.0.1:6000",
		"2": "127.0.0.1:6001",
		"3": "127.0.0.1:6002",
	}

	RegisterCommand(&cmdCalculation{})

	ops := []string{"+", "-"}

	f := func(num, n int32, op string) int32 {
		switch op {
		case "+":
			num += n
		case "-":
			num -= n
		}
		return num
	}

	for name := range peers {
		go func(name string) {
			if name == "3" {
				time.Sleep(20 * time.Second) // 第3个node迟到加入
			}
			num := int32(1)
			applyMsgCh := make(chan ApplyMsg, 10)
			rf := Make(name, peers, applyMsgCh)
			go func() {
				for e := range applyMsgCh {
					cmd := e.Command.(*cmdCalculation)
					num = f(num, cmd.Num, cmd.Operator)
					fmt.Printf("---- raft %s event %v cmd %v num %d\n", name, e, cmd, num)
				}
			}()

			for {
				time.Sleep(time.Second)
				v := rand.Int31n(10) + 1
				op := ops[rand.Intn(2)]
				rf.Submit(&cmdCalculation{
					Operator: op,
					Num:      v,
				})
			}
		}(name)
	}

	time.Sleep(1 * time.Hour)
}

func TestMakeRaft_Snapshot(t *testing.T) {
	debug = 1
	peers := map[string]string{
		"1": "127.0.0.1:6000",
		"2": "127.0.0.1:6001",
		"3": "127.0.0.1:6002",
	}

	RegisterCommand(&cmdCalculation{})

	ops := []string{"+", "-"}

	f := func(num, n int32, op string) int32 {
		switch op {
		case "+":
			num += n
		case "-":
			num -= n
		}
		return num
	}

	for name := range peers {
		go func(name string) {
			if name == "3" {
				time.Sleep(20 * time.Second) // 第3个node迟到加入
			}
			num := int32(1)
			applyMsgCh := make(chan ApplyMsg, 10)
			rf := Make(name, peers, applyMsgCh)
			go func() {
				for e := range applyMsgCh {
					if e.CommandValid {
						cmd := e.Command.(*cmdCalculation)
						atomic.StoreInt32(&num, f(atomic.LoadInt32(&num), cmd.Num, cmd.Operator))
						n := atomic.LoadInt32(&num)
						fmt.Printf("---- raft %s event %v cmd %v num %d\n", name, e, cmd, n)
						if e.CommandIndex%5 == 0 {
							// 每5次安装一次快照
							str := fmt.Sprintf("num=%d", n)
							rf.MakeSnapshot(e.CommandIndex, []byte(str))
							fmt.Printf("---- raft %s MakeSnapshot num %d str %s\n", name, n, str)
						}
					} else {
						str := string(e.Snapshot)
						n, err := strconv.Atoi(strings.Split(str, "=")[1])
						atomic.StoreInt32(&num, int32(n))
						fmt.Printf("---- raft %s read snapshot str %s num %d  err %v \n", name, str, n, err)
					}
				}
			}()

			for {
				time.Sleep(time.Second)
				v := rand.Int31n(10) + 1
				op := ops[rand.Intn(2)]
				rf.Submit(&cmdCalculation{
					Operator: op,
					Num:      v,
				})

			}
		}(name)
	}

	time.Sleep(1 * time.Hour)
}
