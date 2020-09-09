## raft 的应用

一致性 kv 

1. server

server 模拟数据库, 响应 client 的请求。对外服务 set，get 请求。
```
type SetReq struct {
	Key   string
	Value string
}

type SetResp struct {
	OK       bool
	IsLeader bool
}

type GetReq struct {
	Key string
}

type GetResp struct {
	OK    bool
	Value string
}
```

两组配置，分别对应 raft 集群内 rpc 通讯 和 server 对外服务的地址。这里以后尝试统一一个地址。

raft 只有一个命令，用于修改本地数据。
```
type SetCmd struct {
	Key   string
	Value string
}
```

处理 raft 提交的应用层的逻辑

```
	for e := range applyMsgCh {
		if e.CommandValid {
			cmd := e.Command.(*SetCmd)
			fmt.Printf("raft %s event %v cmd %v\n", name, e, cmd)
			kvStorage.Set(cmd.Key, cmd.Value)
			// 每隔10条 一次快照
			if e.CommandIndex%10 == 0 {
				fmt.Printf("raft %s make snapshot\n", name)
				rfMtx.Lock()
				rf.MakeSnapshot(e.CommandIndex, kvStorage.packJson())
				rfMtx.Unlock()
			}
		} else {
			fmt.Printf("---- raft %s snapshot \n", name)
			kvStorage.unpackJson(e.Snapshot)
		}
	}
```


2. client

client 作为客户端，向 server 提交 set，get 请求。

处理 leader 的方式是：先默认一个节点为 leader 。从 leader 开始遍历发送请求（只有set需要提交的raft，需要验证leader身份）。
- 发送成功时该节点设置为 leader
- 失败继续向其他节点请求直到请求成功。若全部失败，表示 raft 集群还没有选出 leader，需稍后重试。

采用监听控制台输入的方式，操作命令

## 启动

1. server

由于配置为 3个节点，所以 server 至少启动两个。
```
go run server/server.go kv1
go run server/server.go kv2
go run server/server.go kv3
```

2. client

启动一个
```
go run client/client.go 
```

