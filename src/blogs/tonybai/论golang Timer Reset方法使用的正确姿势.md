# 一、关于Timer原理的一些说明

在网路编程方面,从用户视角看,golang表象是一种***"阻塞式"***网络编程范式,而支撑这种"阻塞式"范式的则是内置于go编译后的executable file中的runtime。runtime利用***网络IO多路复用机制***实现***多个进行网络通信的goroutine的合理调度***。goroutine中的执行函数则相当于你在传统C编程中传给epoll机制中的回调函数。golang一定程度上消除了在这方面"回调"这种"逆向思维"给你带来的心智负担,简化了网络编程的复杂性。

但是长时间"阻塞"显然不能满足大多数业务情景,因此还需要一定的超时机制。
- socket层面: net.Dial Timeout / SetReadDeadline、SetWriteDeadline、SetDeadline
- 应用层协议: eg: http    client: timeout    server:TimeoutHandler

这些timeout机制,有些是通过runtime的网络多路复用的timeout机制实现的,有些则是通过Timer实现的。

## 1.Timer的创建

```
场景1:

t:=time.AfterFunc(d,f)

场景2:

select{
    case m:=<-c:
        handle(m)
    case <-time.After(5*time.Minute):
        fmt.Println("timed out")
}

或:

t:=time.NewTimer(5*time.Minute)
select{
    case m:=<-c:
        handle(m)
    case <-t.C:
        fmt.Println("timed out")
}
```

从这两个场景中，我们可以看到Timer三种创建姿势：

```
t:= time.NewTimer(d)
t:= time.AfterFunc(d, f)
c:= time.After(d)
```

虽然姿势不同，但背后的原理则是相通的。
Timer有三个要素：

```
* 定时时间：也就是那个d
* 触发动作：也就是那个f
* 时间channel： 也就是t.C
```

对于AfterFunc这种创建方式而言，Timer就是在超时(timer expire)后，执行函数f，此种情况下：时间channel无用。

```
func AfterFunc(d Duration, f func()) *Timer {
    t := &Timer{
        r: runtimeTimer{
            when: when(d),
            f:    goFunc,
            arg:  f,
        },
    }
    startTimer(&t.r)
    return t
}

func goFunc(arg interface{}, seq uintptr) {
    go arg.(func())()
}
```


注意：从AfterFunc源码可以看到，外面传入的f参数并非直接赋值给了内部的f，而是作为wrapper function：goFunc的arg传入的。而goFunc则是启动了一个新的goroutine来执行那个外部传入的f。这是因为timer expire对应的事件处理函数的执行是在go runtime内唯一的timer events maintenance goroutine: timerproc中。为了不block timerproc的执行，必须启动一个新的goroutine。

```
//$GOROOT/src/runtime/time.go
func timerproc() {
    timers.gp = getg()
    for {
        lock(&timers.lock)
        ... ...
            f := t.f
            arg := t.arg
            seq := t.seq
            unlock(&timers.lock)
            if raceenabled {
                raceacquire(unsafe.Pointer(t))
            }
            f(arg, seq)
            lock(&timers.lock)
        }
        ... ...
        unlock(&timers.lock)
   }
}
```


而对于NewTimer和After这两种创建方法，则是Timer在超时(timer expire)后，执行一个标准库中内置的函数：sendTime。sendTime将当前当前事件send到timer的时间Channel中，那么说这个动作不会阻塞到timerproc的执行么？答案肯定是不会的，其原因就在下面代码中：

```
// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.
func NewTimer(d Duration) *Timer {
	c := make(chan Time, 1)
	t := &Timer{
		C: c,
		r: runtimeTimer{
			when: when(d),
			f:    sendTime,
			arg:  c,
		},
	}
	startTimer(&t.r)
	return t
}

func sendTime(c interface{}, seq uintptr) {
	// Non-blocking send of time on c.
	// Used in NewTimer, it cannot block anyway (buffer).
	// Used in NewTicker, dropping sends on the floor is
	// the desired behavior when the reader gets behind,
	// because the sends are periodic.
	select {
	case c.(chan Time) <- Now():
	default:
	}
}
```


我们看到NewTimer中创建了一个buffered channel，size = 1。正常情况下，当timer expire，t.C无论是否有goroutine在read，sendTime都可以non-block的将当前时间发送到C中；同时，我们看到sendTime还加了双保险：通过一个select判断c buffer是否已满，一旦满了，直接退出，依然不会block，这种情况在reuse active timer时可能会遇到。


## 2.Timer的资源释放

很多Go初学者在使用Timer时都会担忧Timer的创建会占用系统资源，比如：

- 有人会认为：创建一个Timer后，runtime会创建一个单独的Goroutine去计时并在expire后发送当前时间到channel里。
- 还有人认为：创建一个timer后，runtime会申请一个os级别的定时器资源去完成计时工作。

实际情况并不是这样。恰好近期gopheracademy blog发布了一篇 《How Do They Do It: Timers in Go》，通过对timer源码的分析，讲述了timer的原理，大家可以看看。
go runtime实际上仅仅是启动了一个单独的goroutine，运行timerproc函数，维护了一个”***最小堆***”，***定期wake up后，读取堆顶的timer，执行timer对应的f函数，并移除该timer element。***
- 创建一个Timer实则就是在这个最小堆中添加一个element
- Stop一个timer，则是从堆中删除对应的element。

同时，从上面的两个Timer常见的使用场景中代码来看，我们并没有显式的去释放什么。从上一节我们可以看到，Timer在创建后可能占用的资源还包括：

* 0或一个Channel
* 0或一个Goroutine

这些资源都会在timer使用后被GC回收。

***我们要做的***:
- 及时调用timer的Stop方法从最小堆删除timer element(如果timer没有expire)
- resuse active timer


