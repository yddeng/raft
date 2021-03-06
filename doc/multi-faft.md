# Multi-Raft

## Multi-Raft需要解决的问题

单个Raft-Group在KV的场景下存在一些弊端:

1. 系统的存储容量受制于单机的存储容量（使用分布式存储除外）
2. 系统的性能受制于单机的性能（读写请求都由Leader节点处理）

Multi-Raft需要解决的一些核心问题：

1. 数据何如分片
2. 分片中的数据越来越大，需要分裂产生更多的分片，组成更多Raft-Group
3. 分片的调度，让负载在系统中更平均（分片副本的迁移，补全，Leader切换等等）
4. 一个节点上，所有的Raft-Group复用链接（否则Raft副本之间两两建链，链接爆炸了）
5. 如何处理stale的请求（例如Proposal和Apply的时候，当前的副本不是Leader、分裂了、被销毁了等等）
6. Snapshot如何管理（限制Snapshot，避免带宽、CPU、IO资源被过度占用）

### 数据如何分片
    
两种分片方式适用于不同的场景
    
1. 按照用户的Key做字典序，系统一开始只有1个分片，分片个数随着系统的数据量逐渐增大而不断分裂（这个实现和TiKV一致）
2. 按照Key的hash值，映射到一个uint64的范围，可以初始化的时候就可以设置N个分片，让系统一开始就可以支持较高的并发，后续随着数据量的增大继续分裂

分片后需要一个唯一的注册表，每个group分别管理的hash值范围。