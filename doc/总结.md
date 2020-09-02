#### 课程有 4 个 Lab 实验：

1. Lab1 MapReduce：熟悉分布式基础概念

    - 阅读论文并梳理 MR 模型的执行流程，实现单机版、分布式版的 word count，最后使用模型来生成倒排索引。

2. Lab2 Raft：分三个模块实现 Raft 一致性算法

    - 2A Leader Election：实现 leader 选举和心跳通信，处理好网络分区、多节点失效及 split vote
    - 2B Log Replication：实现日志复制和一致性检查，在少部分节点失效时依旧能 commit 日志，并解决好日志冲突。
    - 2C Network Unavailable：实现节点状态的持久化和重启读取，应对 RPC 请求乱序、延迟甚至丢失的网络环境、节点频繁崩溃和重启的情况。
3. Lab3 kvDB：基于 Raft 实现线性一致的分布式容错 kv 数据库

    - 3A：基于 Raft 库保证在并发请求下数据的线性一致性，处理 RPC 超时重试，实现请求去重等逻辑。
    - 3B：实现 Raft 的日志压缩及数据快照，实时监测日志大小并剪切，切换 RPC 加速日志回放。
4. Lab4 Sharded kv：实现 Raft 的 Membership Change 配置更新，构建 kvDB 集群。实现中 …

    - 每个 Lab 的环境代码都有注释提示，都有写好的单元测试，目标是每个测试跑 go test -race -count 多次都能 pass