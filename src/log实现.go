前一篇文章我们看到了Golang标准库中log模块的使用，那么它是如何实现的呢？下面我从log.Logger开始逐步分析其实现。 其源码可以参考官方地址
1.Logger结构
首先来看下类型Logger的定义：

type Logger struct {
    mu     sync.Mutex // ensures atomic writes; protects the following fields
    prefix string     // prefix to write at beginning of each line
    flag   int        // properties
    out    io.Writer  // destination for output
    buf    []byte     // for accumulating text to write
}
主要有5个成员，其中3个我们比较熟悉，分别是表示Log前缀的 "prefix"，表示Log头标签的 "flag" ，以及Log的输出目的地out。 buf是一个字节数组，主要用来存放即将刷入out的内容，相当于一个临时缓存，在对输出内容进行序列化时作为存储目的地。 mu是一个mutex主要用来作线程安全的实习，当有多个goroutine同时往一个目的刷内容的时候，通过mutex保证每次写入是一条完整的信息。
2.std及整体结构
在前一篇文章中我们提到了log模块提供了一套包级别的简单接口，使用该接口可以直接将日志内容打印到标准错误。那么该过程是怎么实现的呢？其实就是通过一个内置的Logger类型的变量 "std" 来实现的。该变量使用：

var std = New(os.Stderr, "", LstdFlags)
进行初始化，默认输出到系统的标准输出 "os.Stderr" ,前缀为空，使用日期加时间作为Log抬头。
当我们调用 log.Print的时候是怎么执行的呢？我们看其代码：

func Print(v ...interface{}) {
    std.Output(2, fmt.Sprint(v...))
}
这里实际就是调用了Logger对象的 Output方法，将日志内容按照fmt包中约定的格式转义后传给Output。Output定义如下 :
1
func (l *Logger) Output(calldepth int, s string) error
其中s为日志没有加前缀和Log抬头的具体内容，xxxxx 。该函数执行具体的将日志刷入到对应的位置。
3.核心函数的实现
Logger.Output是执行具体的将日志刷入到对应位置的方法。
该方法首先根据需要获得当前时间和调用该方法的文件及行号信息。然后调用formatHeader方法将Log的前缀和Log抬头先格式化好 放入Logger.buf中，然后再将Log的内容存入到Logger.buf中，最后调用Logger.out.Write方法将完整的日志写入到输出目的地中。
由于写入文件以及拼接buf的过程是线程非安全的，因此使用mutex保证每次写入的原子性。

l.mu.Lock()
defer l.mu.Unlock()
将buf的拼接和文件的写入放入这个后面，使得在多个goroutine使用同一个Logger对象是，不会弄乱buf，也不会杂糅的写入。
该方法的第一个参数最终会传递给runtime.Caller的skip，指的是跳过的栈的深度。这里我记住给2就可以了。这样就会得到我们调用log 是所处的位置。
在golang的注释中说锁住 runtime.Caller的过程比较重，这点我还是不很了解，只是从代码中看到其在这里把锁打开了。


if l.flag&(Lshortfile|Llongfile) != 0 {
    // release lock while getting caller info - it's expensive.
    l.mu.Unlock()
    var ok bool
    _, file, line, ok = runtime.Caller(calldepth)
    if !ok {
        file = "???"
        line = 0
    }
    l.mu.Lock()
}
在formatHeader里面首先将前缀直接复制到Logger.buf中,然后根据flag选择Log抬头的内容，这里用到了一个log模块实现的 itoa的方法，作用类似c的itoa,将一个整数转换成一个字符串。只是其转换后将结果直接追加到了buf的尾部。
纵观整个实现，最值得学习的就是线程安全的部分。在什么位置合适做怎样的同步操作。
4.对外接口的实现
在了解了核心格式化和输出结构后，在看其封装就非常简单了，几乎都是首先用Output进行日志的记录，然后在必要的时候 做os.exit或者panic的操作，这里看下Fatal的实现。

func (l *Logger) Fatal(v ...interface{}) {
    l.Output(2, fmt.Sprint(v...))
    os.Exit(1)
}
// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
    l.Output(2, fmt.Sprintf(format, v...))
    os.Exit(1)
}
// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
    l.Output(2, fmt.Sprintln(v...))
    os.Exit(1)
}
这里也验证了我们之前做的Panic的结果，先做输出日志操作。再进行panic。
5.Golang的log模块设计
Golang的log模块主要提供了三类接口 ：
Print ： 一般的消息输出
Fatal : 类似assert一般的强行退出
Panic ： 相当于OO里面常用的异常捕获
与其说log模块提供了三类日志接口，不如说log模块仅仅是对类C中的 printf、assert、try...catch...的简单封装。Golang的log模块 并没有对log进行分类、分级、过滤等其他类似log4j、log4c、zlog当中常见的概念。当然在使用中你可以通过添加prefix,来进行简单的 分级，或者改变Logger.out改变其输出位置。但这些并没有在API层面给出直观的接口。
Golang的log模块就像是其目前仅专注于为服务器编程一样，他的log模块也专注于服务器尤其是基础组件而服务。就像nginx、redis、lighttpd、keepalived自己为自己写了一个简单的日志模块而没有实现log4c那样庞大且复杂的日志模块一样。他的日志模块仅仅需要为 本服务按照需要的格式和方式提供接口将日志输出到目的地即可。
Golang的log模块可以进行一般的信息记录，assert时的信息输出，以及出现异常时的日志记录，通过对其Print的包装可以实现更复杂的 输出。因此这个log模块可谓是语言层面上非常基础的一层库，反应的是语言本身的特征而不是一个服务应该怎样怎样。

