package main

import (
	"fmt"
	"github.com/yddeng/raft/kv"
	"os"
)

var (
	// raft 集群地址
	peers = map[string]string{
		"kv1": "127.0.0.1:6000",
		"kv2": "127.0.0.1:6001",
		"kv3": "127.0.0.1:6002",
	}
	// 对外服务的启动地址
	addrs = map[string]string{
		"kv1": "127.0.0.1:7000",
		"kv2": "127.0.0.1:7001",
		"kv3": "127.0.0.1:7002",
	}
)

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Sprintf("args : p , name "))
	}

	name := os.Args[1]
	kv.Start(peers, name, addrs[name])
}
