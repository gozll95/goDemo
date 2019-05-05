# 一、Goroutine调度器

## 传统的语言:
提到“调度”，我们首先想到的就是操作系统对进程、线程的调度。操作系统调度器会将系统中的多个线程按照一定算法调度到物理CPU上去运行。传统的编程语言比如C、C++等的并发实现实际上就是基于操作系统调度的，即***程序负责创建线程***(一般通过pthread等lib调用实现)，***操作系统负责调度***。这种传统支持并发的方式有诸多不足.(复杂/难于scaling)

Go的runtime负责对goroutine进行管理。
- 决定何时哪个goroutine将获得资源开始执行
- 哪个goroutine应该停止执行让出资源
- 哪个goroutine应该被唤醒恢复执行等


为此，Go采用了***用户层轻量级thread***或者说是***类coroutine***的概念来解决这些问题，Go将之称为***”goroutine“***。goroutine占用的资源非常小(Go 1.4将每个goroutine stack的size默认设置为2k)，goroutine调度的切换也不用陷入(trap)操作系统内核层完成，代价很低。因此，一个Go程序中可以创建成千上万个并发的goroutine。所有的Go代码都在goroutine中执行，哪怕是go的runtime也不例外。将这些goroutines按照一定算法放到“CPU”上执行的程序就称为***goroutine调度器***或***goroutine scheduler***。
不过，一个Go程序对于操作系统来说只是一个***用户层程序***，对于操作系统而言，它的眼中只有thread，它甚至不知道有什么叫Goroutine的东西的存在。goroutine的调度全要靠Go自己完成，实现Go程序内goroutine之间“公平”的竞争“CPU”资源，这个任务就落到了Go runtime头上，要知道在一个Go程序中，除了用户代码，剩下的就是go runtime了。
于是Goroutine的调度问题就演变为go runtime如何将程序内的众多goroutine按照一定算法调度到“CPU”资源上运行了。在操作系统层面，Thread竞争的“CPU”资源是真实的物理CPU，但在Go程序层面，各个Goroutine要竞争的”CPU”资源是什么呢？Go程序是用户层程序，它本身整体是运行在一个或多个操作系统线程上的，因此goroutine们要竞争的所谓“CPU”资源就是***操作系统线程***。这样Go scheduler的任务就明确了：将goroutines按照一定算法放到不同的操作系统线程中去执行。这种在语言层面自带调度器的，我们称之为***原生支持并发***。


# Go调度器模型与演化过程
## G-M模型-->G-P-M模型

                                        Global队列
                                    --入队--G|G|G|G|----出队
                            G         G       G       G
                            G         G       G       G          P的Local队列
                            G         G       G       G
                            G         G       G       G      
            ---- 基本单元     |         |       |       |
            |    (Goroutine) |         |       |       |
            |                |         |       |       |
            |    
 Goroutine  |
 调度器      |        逻辑       P         P       P       P  <-------GOMAXPROCS个
            |    (Processor)
            |
            |    物理Processor   M  M   M   M   M     M    M
        ——— |
        |        (User-level Thread)
--------|--------------------------------------------
        |
        |
 操作系统 |
 调度器  |       Kernel Thread
        |
        |
        |
--------|--------------------------------------------
        |   
        |___    物理多核CPU  CPU-1|CPU-2|CPU-3|CPU-4        物理机器


有名人曾说过：“***计算机科学领域的任何问题都可以通过增加一个间接的中间层来解决***”，我觉得Dmitry Vyukov的***G-P-M***模型恰是这一理论的践行者。Dmitry Vyukov通过向G-M模型中增加了一个P，实现了Go scheduler的scalable。
P是一个“逻辑Proccessor”，每个G要想真正运行起来，首先需要被分配一个P（进入到P的local runq中，这里暂忽略global runq那个环节）。对于G来说，P就是运行它的“CPU”，可以说：G的眼里只有P。但从Go scheduler视角来看，真正的“CPU”是M，只有将P和M绑定才能让P的runq中G得以真实运行起来。这样的P与M的关系，就好比Linux操作系统调度层面用户线程(user thread)与核心线程(kernel thread)的对应关系那样(N x M)。

-总结:GGG->local runq--->P---->M ---->

## 抢占式调度

G-P-M模型的实现算是Go scheduler的一大进步，但Scheduler仍然有一个头疼的问题，那就是不支持抢占式调度，导致***一旦某个G中出现死循环或永久循环的代码逻辑，那么G将永久占用分配给它的P和M，位于同一个P中的其他G将得不到调度***，出现***“饿死”***的情况。更为严重的是，当只有一个P时(GOMAXPROCS=1)时，整个Go程序中的其他G都将“饿死”。于是Dmitry Vyukov又提出了《Go Preemptive Scheduler Design》并在Go 1.2中实现了“抢占式”调度。

这个***抢占式调度***的原理则是在每个函数或方法的入口，加上一段额外的代码，让runtime有机会检查是否需要执行抢占调度。这种解决方案只能说局部解决了“饿死”问题，对于没有函数调用，纯算法循环计算的G，scheduler依然无法抢占。

## 其他优化
Go runtime已经实现了netpoller，这使得即便G发起网络I/O操作也不会导致M被阻塞（仅阻塞G），从而不会导致大量M被创建出来。但是对于regular file的I/O操作一旦阻塞，那么M将进入sleep状态，等待I/O返回后被唤醒；这种情况下P将与sleep的M分离，再选择一个idle的M。如果此时没有idle的M，则会新创建一个M，这就是为何大量I/O操作导致大量Thread被创建的原因。
Ian Lance Taylor在Go 1.9 dev周期中增加了一个Poller for os package的功能，这个功能可以像netpoller那样，在G操作支持pollable的fd时，仅阻塞G，而不阻塞M。不过该功能依然不能对regular file有效，regular file不是pollable的。不过，对于scheduler而言，这也算是一个进步了。



# 三、Go调度器原理的进一步理解
## G、P、M
于G、P、M的定义，大家可以参见$GOROOT/src/runtime/runtime2.go这个源文件。这三个struct都是大块儿头，每个struct定义都包含十几个甚至二、三十个字段。像scheduler这样的核心代码向来很复杂，考虑的因素也非常多，代码“耦合”成一坨。不过从复杂的代码中，我们依然可以看出来G、P、M的各自大致用途（当然雨痕老师的源码分析功不可没），这里简要说明一下：
* G:  表示goroutine，存储了goroutine的执行stack信息、goroutine状态以及goroutine的任务函数等；另外G对象是可以重用的。
* P:  表示逻辑processor，P的数量决定了系统内最大可并行的G的数量（前提：系统的物理cpu核数>=P的数量）；P的最大作用还是其拥有的各种G对象队列、链表、一些cache和状态。
* M:  M代表着真正的执行计算资源。在绑定有效的p后，进入schedule循环；而schedule循环的机制大致是从各种队列、p的本地队列中获取G，切换到G的执行栈上并执行G的函数，调用goexit做清理工作并回到m，如此反复。M并不保留G状态，这是G可以跨M调度的基础。

## G被抢占调度

和操作系统按时间片调度线程不同，Go并没有时间片的概念。如果某个G没有进行system call调用、没有进行I/O操作、没有阻塞在一个channel操作上，那么m是如何让G停下来并调度下一个runnable G的呢？答案是：G是被抢占调度的。

前面说过，除非极端的无限循环或死循环，否则只要G调用函数，Go runtime就有抢占G的机会。Go程序启动时，runtime会去启动一个名为***sysmon的m(一般称为监控线程)，该m无需绑定p即可运行***，该m在整个Go程序的运行过程中至关重要：

***启动时候启动sysmon的m(监控线程)***

sysmon每20us~10ms启动一次，按照《Go语言学习笔记》中的总结，sysmon主要完成如下工作：
* 释放闲置超过5分钟的span物理内存；
* 如果超过2分钟没有垃圾回收，强制执行；
* 将长时间未处理的netpoll结果添加到任务队列；
* 向长时间运行的G任务发出抢占调度；
* 收回因syscall长时间阻塞的P；

可以看出，如果一个G任务运行10ms，sysmon就会认为其运行时间太久而发出抢占式调度的请求。一旦G的抢占标志位被设为true，那么待这个G下一次调用函数或方法时，runtime便可以将G抢占，并移出运行状态，放入P的local runq中，等待下一次被调度。

## 3、channel阻塞或network I/O情况下的调度

如果G被阻塞在某个channel操作或network I/O操作上时，G会被放置到某个wait队列中，而M会尝试运行下一个runnable的G；如果此时没有runnable的G供m运行，那么m将解绑P，并进入sleep状态。当I/O available或channel操作完成，在wait队列中的G会被唤醒，标记为runnable，放入到某P的队列中，绑定一个M继续执行。


## 4、 system call阻塞情况下的调度
如果G被阻塞在某个system call操作上，那么不光G会阻塞，执行该G的M也会解绑P(实质是被sysmon抢走了)，与G一起进入sleep状态。如果此时有idle的M，则P与其绑定继续执行其他G；如果没有idle M，但仍然有其他G要去执行，那么就会创建一个新M。
当阻塞在syscall上的G完成syscall调用后，G会去尝试获取一个可用的P，如果没有可用的P，那么G会被标记为runnable，之前的那个sleep的M将再次进入sleep。


# 四、调度器状态的查看方法
```
$GODEBUG=schedtrace=1000 godoc -http=:6060
SCHED 0ms: gomaxprocs=4 idleprocs=3 threads=3 spinningthreads=0 idlethreads=0 runqueue=0 [0 0 0 0]
SCHED 1001ms: gomaxprocs=4 idleprocs=0 threads=9 spinningthreads=0 idlethreads=3 runqueue=2 [8 14 5 2]
SCHED 2006ms: gomaxprocs=4 idleprocs=0 threads=25 spinningthreads=0 idlethreads=19 runqueue=12 [0 0 4 0]
SCHED 3006ms: gomaxprocs=4 idleprocs=0 threads=26 spinningthreads=0 idlethreads=8 runqueue=2 [0 1 1 0]
SCHED 4010ms: gomaxprocs=4 idleprocs=0 threads=26 spinningthreads=0 idlethreads=20 runqueue=12 [6 3 1 0]
SCHED 5010ms: gomaxprocs=4 idleprocs=0 threads=26 spinningthreads=1 idlethreads=20 runqueue=17 [0 0 0 0]
SCHED 6016ms: gomaxprocs=4 idleprocs=0 threads=26 spinningthreads=0 idlethreads=20 runqueue=1 [3 4 0 10]
... ...
```


GODEBUG这个Go运行时环境变量很是强大，通过给其传入不同的key1=value1,key2=value2… 组合，Go的runtime会输出不同的调试信息，比如在这里我们给GODEBUG传入了”schedtrace=1000″，其含义就是每1000ms，打印输出一次goroutine scheduler的状态，每次一行。每一行各字段含义如下：
以上面例子中最后一行为例：

```
SCHED 6016ms: gomaxprocs=4 idleprocs=0 threads=26 spinningthreads=0 idlethreads=20 runqueue=1 [3 4 0 10]

SCHED：调试信息输出标志字符串，代表本行是goroutine scheduler的输出；
6016ms：即从程序启动到输出这行日志的时间；
gomaxprocs: P的数量；
idleprocs: 处于idle状态的P的数量；通过gomaxprocs和idleprocs的差值，我们就可知道执行go代码的P的数量；
threads: os threads的数量，包含scheduler使用的m数量，加上runtime自用的类似sysmon这样的thread的数量；
spinningthreads: 处于自旋状态的os thread数量；
idlethread: 处于idle状态的os thread的数量；
runqueue=1： go scheduler全局队列中G的数量；
[3 4 0 10]: 分别为4个P的local queue中的G的数量。
```

# 关于debug
https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs