在golangweekly的第36期Go Newsletter中我发现一篇短文"How Goroutines Work" ，其作者在参考了诸多资料后，简短概要地总结了一下 Goroutine的工作原理，感觉十分适合刚入门的Gophers（深入理解Goroutine调度的话，可以参考Daniel Morsing的" The Go scheduler" )。这里粗译如下。

# 一、Go语言简介

如果你是Go语言新手，或如果你对"并发(Concurrency)不是并行(parallelism)"这句话毫无赶脚，那么请看一下Rob Pike大神关于这个主题的演讲吧，演讲共30分 钟，我敢保证你在这个演讲上花费30分钟是绝对值得的。

总结一下两者（Concurrency和Parallelism）的不同："当人们听到并发（Concurrency)这个词时，总是会想起并行 （Parallelism），它们之间有相关性，但却是两个明显不同的概念。在编程领域，并发（Concurrency）是独立的执行过程 (Process)的组合，而并行（Parallelism)则是计算（可能是相关联的）的同时执行。并发（Concurrency)是关于同时 应对很多事情(deal with lots of things)，而并行（Parallelism)则是同时做许多事情(do lots of things)"。(Rob Pike的“Concurrency is not parallelism")

Go语言支持我们编写并发(Concurrent)的程序。它提供了Goroutine以及更重要的在Goroutines之间通信的能力。这里 我们将聚焦在前者（译注：指并发）。

# 二、Goroutines和Threads

Goroutine是一个简单的模型：它是一个函数，与其他Goroutines并发执行且共享相同地址空间。Goroutines的通常用法是根据需要创建尽可 能的Groutines，成百上千甚至上万的。这种用法对于那些习惯了使用C++或Java的程序员来讲可能会有些奇怪。创建这么多 goroutines势必要付出不菲的代价？一个操作系统线程使用固定大小的内存作为它的执行栈，当线程数增多时，线程间切换的代价也是相当的 高。这也是每处理一个request就创建一个新线程的服务程序方案被诟病的原因。

不过Goroutine完全不同。它们由Go运行时初始化并调度，操作系统根本看不到Goroutine的存在。所有的goroutines都是 活着的，并且以多路复用的形式运行于操作系统为应用程序分配的少数几个线程上。创建一个Goroutine并不需要太多内存，只需要8K的栈空间 （在Go 1.3中这个Size发生了变化）。它们根据需要在堆上分配和释放内存以实现自身的增长。

Go运行时负责调度Goroutines。Goroutines的调度是协作式的，而线程不是。这意味着每次一个线程发生切换，你都需要保存/恢 复所有寄存器，包括16个通用寄存器、PC(程序计数器）、SP（栈指针）、段寄存器(segment register)、16个XMM寄存器、FP协处理器状态、X AVX寄存器以及所有MSR等。而当另一个Goroutine被调度时，只需要保存/恢复三个寄存器，分别是PC、SP和DX。Go调度器和任何现代操作 系统的调度器都是O(1)复杂度的，这意味着增加线程/goroutines的数量不会增加切换时间，但改变寄存器的代价是不可忽视的。

由于Goroutines的调度是协作式的，***一个持续循环的goroutine会导致运行于同一线程上的其他goroutines“饿死”。***在 Go 1.2中，这个问题或多或少可以通过在进入函数前间或地调用Go调度器来缓解一些，因此一个包含非内联函数调用的循环是可以被调度器抢占的。

# 三、Goroutine阻塞

只要阻塞存在，它在OS线程中就是不受欢迎的，因为你拥有的线程数量很少。如果你发现大量线程阻塞在网络操作或是Sleep操作上，那就是问题， 需要修正。正如前面提到的那样，Goroutine是廉价的。更关键地是，如果它们在网络输入操作、Sleep操作、Channel操作或 sync包的原语操作上阻塞了，也不会导致承载其多路复用的线程阻塞。如果一个goroutine在上述某个操作上阻塞，Go运行时会调度另外一 个goroutine。即使成千上万的Goroutine被创建了出来，如果它们阻塞在上述的某个操作上，也不会浪费系统资源。从操作系统的视角来看，你的程序的行为就像是一个事件驱动的C程序似的。

# 四、最后的想法

就是这样，Goroutines可以并发的运行。不过和其他语言一样，组织两个或更多goroutine同时访问共享资源是很重要的。最好采用Channel在不同Goroutine间传递数据。

最后，虽然你无法直接控制Go运行时创建的线程的数量，但可以通过调用runtime.GOMAXPROCS(n)方法设置变量GOMAXPROCS来设 定使用的处理器核的数量。提高使用的处理器核数未必能提升你的程序的性能，这取决于程序的设计。程序剖析诊断工具(profiling tool)可以用来检查你的程序使用处理器核数的真实情况。