package kv

import (
	"encoding/json"
	"fmt"
	"github.com/yddeng/raft"
	"net/http"
	"sync"
)

var (
	rf        *raft.Raft
	kvStorage *KV
	rfMtx     sync.Mutex
)

type KV struct {
	Data map[string]string
	sync.Mutex
}

func (this *KV) Set(k, v string) {
	this.Lock()
	defer this.Unlock()
	this.Data[k] = v
}

func (this *KV) Get(k string) (value string, ok bool) {
	this.Lock()
	defer this.Unlock()
	value, ok = this.Data[k]
	return
}

func (this *KV) packJson() []byte {
	data, err := json.Marshal(this.Data)
	if err != nil {
		fmt.Println("packJson", err)
		return nil
	}
	return data
}

func (this *KV) unpackJson(data []byte) {
	err := json.Unmarshal(data, &this.Data)
	if err != nil {
		fmt.Println("unpackJson", err)
	}
}

func Start(peers map[string]string, name, addr string) {
	raft.RegisterCommand(&SetCmd{})

	http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("handler set")
		msg := SetReq{}
		err := json.NewDecoder(request.Body).Decode(&msg)
		defer request.Body.Close()
		if err != nil {
			fmt.Println("handler set decode ", err)
			return
		}

		rfMtx.Lock()
		isLeader, err := rf.Submit(&SetCmd{
			Key:   msg.Key,
			Value: msg.Value,
		})
		rfMtx.Unlock()

		resp := &SetResp{OK: true, IsLeader: isLeader}
		if !isLeader || err != nil {
			resp.OK = false
			fmt.Println("handler set err", err)
		}

		err = json.NewEncoder(writer).Encode(resp)
		if err != nil {
			fmt.Println("set resp err", err)
		}
	})

	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("handler get")
		msg := GetReq{}
		err := json.NewDecoder(request.Body).Decode(&msg)
		defer request.Body.Close()
		if err != nil {
			fmt.Println("handler get decode ", err)
			return
		}

		v, ok := kvStorage.Get(msg.Key)
		resp := &GetResp{OK: ok, Value: v}
		err = json.NewEncoder(writer).Encode(resp)
		if err != nil {
			fmt.Println("get resp err", err)
		}
	})

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Println("ListenAndServe ", err)
		}
	}()

	kvStorage = &KV{Data: map[string]string{}}

	applyMsgCh := make(chan raft.ApplyMsg, 10)
	rf = raft.Make(name, peers, applyMsgCh)

	// 读命令
	for e := range applyMsgCh {
		if e.CommandValid {
			cmd := e.Command.(*SetCmd)
			fmt.Printf("raft %s event %v cmd %v\n", name, e, cmd)
			kvStorage.Set(cmd.Key, cmd.Value)
		} else {
			fmt.Printf("---- raft %s snapshot \n", name)
			kvStorage.unpackJson(e.Snapshot)
		}
	}
}
