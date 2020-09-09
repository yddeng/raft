package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yddeng/raft/kv"
	"net/http"
)

var (
	// raft 对外地址
	peers = map[string]string{
		"kv1": "127.0.0.1:7000",
		"kv2": "127.0.0.1:7001",
		"kv3": "127.0.0.1:7002",
	}
	list []string
)

func makeList(leader string) {
	list = []string{leader}
	for name := range peers {
		if name != leader {
			list = append(list, name)
		}
	}
}

func set(k, v string) {
	obj := kv.SetReq{Key: k, Value: v}
	byts, _ := json.Marshal(obj)

	for _, name := range list {
		url := fmt.Sprintf("http://%s/set", peers[name])
		reader := bytes.NewReader(byts)
		resp, err := http.Post(url, "application/json", reader)
		if err != nil {
			fmt.Println("set post err", err)
			continue
		}
		msg := kv.SetResp{}
		_ = json.NewDecoder(resp.Body).Decode(&msg)
		_ = resp.Body.Close()

		if !msg.IsLeader {
			continue
		}

		if msg.OK {
			makeList(name)
			fmt.Println("set ok", obj)
			return
		}
	}

	fmt.Println("set failed", obj)
}

func get(k string) string {
	obj := kv.GetReq{Key: k}
	byts, _ := json.Marshal(obj)

	for _, name := range list {
		url := fmt.Sprintf("http://%s/get", peers[name])
		reader := bytes.NewReader(byts)
		resp, err := http.Post(url, "application/json", reader)
		if err != nil {
			fmt.Println("get post err", err)
			continue
		}
		msg := kv.GetResp{}
		_ = json.NewDecoder(resp.Body).Decode(&msg)
		_ = resp.Body.Close()

		if msg.OK {
			makeList(name)
			fmt.Println("get ok", obj, msg)
			return msg.Value
		}
	}
	return ""
}

func main() {
	makeList("kv1")

	for {
		var op int
		fmt.Println("<<< op 1 -> set, 2 -> get >>>")
		fmt.Print(">>> ")
		fmt.Scan(&op)
		switch op {
		case 1:
			var k, v string
			fmt.Println("<<< set, input k,v >>>")
			fmt.Print(">>> ")
			fmt.Scan(&k, &v)
			set(k, v)
			fmt.Println()
		case 2:
			var k string
			fmt.Println("<<< get, input k >>>")
			fmt.Print(">>> ")
			fmt.Scan(&k)
			v := get(k)
			fmt.Printf("====> %s -> %s \n", k, v)
			fmt.Println()
		default:

		}
	}
}
