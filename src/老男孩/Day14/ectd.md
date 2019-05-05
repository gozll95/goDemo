# ectd介绍
概念: 高可用的分布式key-value存储,可以用于配置共享和服务发现
类似项目: zk/consul
开发语言: Go
接口: 提供restful的http接口,使用简单
实现算法:基于raft算法的强一致性、高可用的服务存储目录

# ectd的应用场景
- 服务发现和服务注册
- 配置中心
- 分布式锁 ？？？？？？ 再次复习下redis吧
- master选举

context_trace

```
func add(ctx context.Context, a, b int) int {
	traceId := ctx.Value("trace_id").(string)
	fmt.Printf("add trace_id:%v\n", traceId)
	return a + b
}
func calc(ctx context.Context, a, b int) int {
	traceId := ctx.Value("trace_id").(string)
	fmt.Printf("calc trace_id:%v\n", traceId)
	return add(ctx, a, b)
}
func main() {
	ctx := context.WithValue(context.Background(), "trace_id", "123456")
	calc(ctx, 388, 200)
}
```

# ectd watch 来监控
ectd watch 使用 push

但是不靠谱(万一etcd watch有异常的话)

所以建议***推拉结合***的方式
etcd push. 我来pull


# 限流
1s 最多收集1000条
记录当前秒数

type SecondLimit struct{
    unixSecond int64  // 当前的秒数
    curCount int32 // 当前收集的日志数
    limit int32
}

当过了这个秒数,就重置