//http://docs.ruanjiadeng.com/gopl-zh/ch11/ch11-05.html

$ go test -cpuprofile=cpu.out 
$ go test -blockprofile=block.out 
$ go test -memprofile=mem.out



对于一些非测试程序也很容易支持分析的特性, 具体的实现方式和程序是短时间运行的小工具还是长时间运行的服务会有很大不同, 因此Go的runtim运行时包提供了程序运行时控制分析特性的接口.

一旦我们已经收集到了用于分析的采样数据, 我们就可以使用 pprof 据来分析这些数据. 这是Go工具箱自带的一个工具, 但并不是一个日常工具, 它对应 go tool pprof 命令. 该命令有许多特性和选项, 但是最重要的有两个, 就是生成这个概要文件的可执行程序和对于的分析日志文件.

为了提高分析效率和减少空间, 分析日志本身并不包含函数的名字; 它只包含函数对应的地址. 也就是说pprof需要和分析日志对于的可执行程序. 虽然 go test 命令通常会丢弃临时用的测试程序, 但是在启用分析的时候会将测试程序保存为 foo.test 文件, 其中 foo 部分对于测试包的名字.

下面的命令演示了如何生成一个CPU分析文件. 我们选择 net/http 包的一个基准测试. 通常是基于一个已经确定了是关键代码的部分进行基准测试. 基准测试会默认包含单元测试, 这里我们用 -run=NONE 禁止单元测试.

$ go test -run=NONE -bench=ClientServerParallelTLS64 \
    -cpuprofile=cpu.log net/http
 PASS
 BenchmarkClientServerParallelTLS64-8  1000
    3141325 ns/op  143010 B/op  1747 allocs/op 
ok       net/http       3.395s

$ go tool pprof -text -nodecount=10 ./http.test cpu.log
2570ms of 3590ms total (71.59%)
Dropped 129 nodes (cum <= 17.95ms)
Showing top 10 nodes out of 166 (cum >= 60ms)
    flat  flat%   sum%     cum   cum%
  1730ms 48.19% 48.19%  1750ms 48.75%  crypto/elliptic.p256ReduceDegree
   230ms  6.41% 54.60%   250ms  6.96%  crypto/elliptic.p256Diff
   120ms  3.34% 57.94%   120ms  3.34%  math/big.addMulVVW
   110ms  3.06% 61.00%   110ms  3.06%  syscall.Syscall 
    90ms  2.51% 63.51%  1130ms 31.48%  crypto/elliptic.p256Square
    70ms  1.95% 65.46%   120ms  3.34%  runtime.scanobject
    60ms  1.67% 67.13%   830ms 23.12%  crypto/elliptic.p256Mul
    60ms  1.67% 68.80%   190ms  5.29%  math/big.nat.montgomery
    50ms  1.39% 70.19%    50ms  1.39%  crypto/elliptic.p256ReduceCarry
    50ms  1.39% 71.59%    60ms  1.67%  crypto/elliptic.p256Sum
参数 -text 标志参数用于指定输出格式, 在这里每行是一个函数, 根据使用CPU的时间来排序. 其中 -nodecount=10 标志参数限制了只输出前10行的结果. 对于严重的性能问题, 这个文本格式基本可以帮助查明原因了.

这个概要文件告诉我们, HTTPS基准测试中 crypto/elliptic.p256ReduceDegree 函数占用了将近一般的CPU资源. 相比之下, 如果一个概要文件中主要是runtime包的内存分配的函数, 那么减少内存消耗可能是一个值得尝试的优化策略.

对于一些更微妙的问题, 你可能需要使用 pprof 的图形显示功能. 这个需要安装 GraphViz 工具, 可以从 www.graphviz.org 下载. 参数 -web 用于生成一个有向图文件, 包含CPU的使用和最特点的函数等信息.

这一节我们只是简单看了下Go语言的分析据工具. 如果想了解更多, 可以阅读 Go官方博客的 ‘‘Proﬁling Go Programs’’ 一文.