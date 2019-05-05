package main

import (
	"fmt"
	"golang.org/x/net/context"
	"time"
)

// 模拟一个最小执行时间的阻塞函数
func inc(a int) int {
	res := a + 1                // 虽然我只做了一次简单的 +1 的运算,
	time.Sleep(1 * time.Second) // 但是由于我的机器指令集中没有这条指令,
	// 所以在我执行了 1000000000 条机器指令, 续了 1s 之后, 我才终于得到结果。B)
	return res
}

// 向外部提供的阻塞接口
// 计算 a + b, 注意 a, b 均不能为负
// 如果计算被中断, 则返回 -1
func Add(ctx context.Context, a, b int) int {
	res := 0
	for i := 0; i < a; i++ {
		res = inc(res)
		select {
		case <-ctx.Done():
			return -1
		default:
		}
	}
	for i := 0; i < b; i++ {
		res = inc(res)
		select {
		case <-ctx.Done():
			return -1
		default:
		}
	}
	return res
}

func main() {
	{
		// 使用开放的 API 计算 a+b
		a := 1
		b := 2
		timeout := 2 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		res := Add(ctx, 1, 2)
		fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
	}
	{
		// 手动取消
		a := 1
		b := 2
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(2 * time.Second)
			cancel() // 在调用处主动取消
		}()
		res := Add(ctx, 1, 2)
		fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
	}
}

//url:https://studygolang.com/articles/9517

/*

简介:
golang中的创建一个新的goroutine,并不会返回像C语言类似的pid,所以我们不能从外部杀死某个goroutine,所以我就得让它自己结束，之前我们用channel+select的方式，来解决这个问题。但是有些场景实现起来比较麻烦，例如
由一个请求衍生出的各个goroutine之间需要满足一定的约束关系，以实现一些诸如有效期,中止routine树，传递请求全局
变量之类的功能。于是google就为我们提供了一个解决方案，开源了context包。使用context实现上下文功能约定需要在
你的方法的传入参数的第一个传入一个context.Context类型的变量。

源码剖析:
context.Context接口

context包的核心
//context包的方法是线程安全的，可以被多个goroutine使用

type Context interface{
    //当Context被cancled或是times out的时候,Done返回一个被closed的channel
    Done() <-chan struct{}

    //在Done的channel被closed后,Err代表被关闭的原因
    Err() error

    //如果存在,Deadline返回Context将要关闭的时间
    Deadline()(deadline time.Time,ok bool)

    //如果存在,Value返回与key相关的值，不存在返回nil
    Value(key interface{}) interface{}
}

我们不需要手动实现这个接口，context包已经给我们提供了两个,一个是Background(),一个是TODO(),这两个函数都会返回一个Context的实例。只是返回的两个实例都是空
Context。


主要结构
cancelCtx 结构体继承了 Context ，实现了 canceler 方法：

//*cancelCtx 和 *timerCtx 都实现了canceler接口，实现该接口的类型都可以被直接canceled
type canceler interface {
    cancel(removeFromParent bool, err error)
    Done() <-chan struct{}
}

type cancelCtx struct {
    Context
    done chan struct{} // closed by the first cancel call.
    mu       sync.Mutex
    children map[canceler]bool // set to nil by the first cancel call
    err      error             // 当其被cancel时将会把err设置为非nil
}

func (c *cancelCtx) Done() <-chan struct{} {
    return c.done
}

func (c *cancelCtx) Err() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.err
}

func (c *cancelCtx) String() string {
    return fmt.Sprintf("%v.WithCancel", c.Context)
}

//核心是关闭c.done
//同时会设置c.err = err, c.children = nil
//依次遍历c.children，每个child分别cancel
//如果设置了removeFromParent，则将c从其parent的children中删除
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
    if err == nil {
        panic("context: internal error: missing cancel error")
    }
    c.mu.Lock()
    if c.err != nil {
        c.mu.Unlock()
        return // already canceled
    }
    c.err = err
    close(c.done)
    for child := range c.children {
        // NOTE: acquiring the child's lock while holding parent's lock.
        child.cancel(false, err)
    }
    c.children = nil
    c.mu.Unlock()

    if removeFromParent {
        removeChild(c.Context, c) // 从此处可以看到 cancelCtx的Context项是一个类似于parent的概念
    }
}



timerCtx 结构继承 cancelCtx

type timerCtx struct {
    cancelCtx //此处的封装为了继承来自于cancelCtx的方法，cancelCtx.Context才是父亲节点的指针
    timer *time.Timer // Under cancelCtx.mu. 是一个计时器
    deadline time.Time
}
valueCtx 结构继承 cancelCtx

type valueCtx struct {
    Context
    key, val interface{}
}


主要方法
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithValue(parent Context, key interface{}, val interface{}) Context
WithCancel 对应的是 cancelCtx ,其中，返回一个 cancelCtx ，同时返回一个 CancelFunc，CancelFunc 是 context 包中定义的一个函数类型：type CancelFunc func()。调用这个 CancelFunc 时，关闭对应的c.done，也就是让他的后代goroutine退出。

WithDeadline 和 WithTimeout 对应的是 timerCtx ，WithDeadline 和 WithTimeout 是相似的，WithDeadline 是设置具体的 deadline 时间，到达 deadline 的时候，后代 goroutine 退出，而 WithTimeout 简单粗暴，直接 return WithDeadline(parent, time.Now().Add(timeout))。

WithValue 对应 valueCtx ，WithValue 是在 Context 中设置一个 map，拿到这个 Context 以及它的后代的 goroutine 都可以拿到 map 里的值。

详细 context 包源码解读:go源码解读



使用原则
使用 Context 的程序包需要遵循如下的原则来满足接口的一致性以及便于静态分析

不要把 Context 存在一个结构体当中，显式地传入函数。Context 变量需要作为第一个参数使用，一般命名为ctx

即使方法允许，也不要传入一个 nil 的 Context ，如果你不确定你要用什么 Context 的时候传一个 context.TODO

使用 context 的 Value 相关方法只应该用于在程序和接口中传递的和请求相关的元数据，不要用它来传递一些可选的参数

同样的 Context 可以用来传递到不同的 goroutine 中，Context 在多个goroutine 中是安全的
*/
