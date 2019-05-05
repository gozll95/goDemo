Golang的标准库提供了log的机制，但是该模块的功能较为简单（看似简单，其实他有他的设计思路）。不过比手写fmt. Printxxx还是强很多的。至少在输出的位置做了线程安全的保护。其官方手册见Golang log (天朝的墙大家懂的)。这里给出一个简单使用的例子：

package main
import (
    "log"
)
func main(){
    log.Fatal("Come with fatal,exit with 1 \n")
}
编译运行后，会看到程序打印了 Come with fatal,exit with 1 然后就退出了，如果用 echo $? 查看退出码，会发现是 “1”。
一般接口
Golang's log模块主要提供了3类接口。分别是 “Print 、Panic 、Fatal ”。当然是用前先包含log包。

import(
    "log"
)
为了方便是用，Golang和Python一样，在提供接口时，提供一个简单的包级别的使用接口。不同于Python，其输出默认定位到标准错误 可以通过SetOutput 进行修改。
对每一类接口其提供了3中调用方式，分别是 "Xxxx 、 Xxxxln 、Xxxxf" 比如对于Print就有:


log.Print
log.Printf
log.Println
log.Print ：表示其参数的调用方式和 fmt.Print 是类似的，即输出对象而不用给定特别的标志符号。
log.Printf ： 表示其参数的调用方式和 fmt.Printf 是类似的，即可以用C系列的格式化标志表示输出对象的类型，具体类型表示 可以参考fmt.Printf的文档
log.Println： 表示其调用方式和fmt.Println 类似，其和log.Print基本一致，仅仅是在输出的时候多输出一个换行
这里是以 “Print”来具体说明的，对于“Panic”和“Fatal”也是一样的。下面再以"Print"为例，看下调用方式：


package main
import (
    "log"
)
func main(){
    arr := []int {2,3}
    log.Print("Print array ",arr,"\n")
    log.Println("Println array",arr)
    log.Printf("Printf array with item [%d,%d]\n",arr[0],arr[1])
}
会得到如下结果：


2014/05/02 12:27:19 Print array [2 3]
2014/05/02 12:27:19 Println array [2 3]
2014/05/02 12:27:19 Printf array with item [2,3]
输出中的日期和时间是默认的格式，如果直接调用简单接口，其格式是固定的，可以通过 SetFlags 方法进行修改，同时这里输出 内容的（传个log.Print的内容）前面和时间的后面是空的，这也是默认的行为，我们可以通过添加前缀来表示其是一条"Warnning" 或者是一条"Debug"日志。通过使用 SetPrefix 可以设置该前缀。
SetOutputSetFlagsSetPrefix 这里关系不大，先不解释，留到后面介绍Logger类型中一并解释。
看完了 log.PrintXxx 接口，我们再来看下 log.FatalXxx 接口，我们以 log.Fatal 为例介绍其功能。如最开始看到的例子， 在调用 log.Fatal 接口后，会先将日志内容打印到标准输出，接着调用系统的 os.exit(1) 接口，退出程序返回状态为 “1”
比较复杂的是 log.PanicXxx ，看该函数的说明，其相当于再把日志内容刷到标准错误后调用 panic 函数(不清楚Golan的defer-recover-panic机制可以Golang Blog去学习一下)。这里举个常用的例子：


package main
import (
    "log"
    "fmt"
)
func main(){
    defer func(){
        if e:= recover();e!= nil {
            fmt.Println("Just comming recover")
            fmt.Println("e from recover is :",e)
            fmt.Println("After recover")
        }
    }()
    arr := []int {2,3}
    log.Panic("Print array ",arr,"\n")
}
结果为：
1
2
3
4
2014/05/03 13:52:42 Print array [2 3]
Just comming recover
e from recover is : Print array [2 3]
After recover
从结果我们可以看出，是先将日志刷入标准输出，然后通过defer里面的recover进行捕获panic的内容。
自定义Logger类型
理清了“Print 、Panic 、Fatal ”后我们就好介绍 log.Logger 类型了。该类型提供了一个New方法用来创建对象。
1
func New(out io.Writer, prefix string, flag int) *Logger
其初始化条件分别是日志写入的位置 out ，日志的前缀内容 prefix ，以及日志的内容flag。可以通过上面介绍的 SetOutputSetFlagsSetPrefix 依次对其进行设置。
输出位置out，是一个io.Writer对象，该对象可以是一个文件也可以是实现了该接口的对象。通常我们可以用这个来指定 其输出到哪个文件
prefix 我们在前面已经看到，就是在日志内容前面的内容。我们可以将其置为 "[Info]" 、 "[Warning]"等来帮助区分日志 级别。
flags 较为迷惑，其实际上就是一个选项，可选的值有：


Ldate         = 1 << iota     // the date: 2009/01/23 形如 2009/01/23 的日期
Ltime                         // the time: 01:23:23   形如 01:23:23   的时间
Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  形如01:23:23.123123   的时间
Llongfile                     // full file name and line number: /a/b/c/d.go:23 全路径文件名和行号
Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile 文件名和行号
LstdFlags     = Ldate | Ltime // 日期和时间
表示在日志内容开头，我们暂且称之为日志抬头，打印出相关内容。对于上面的默认格式就是 LstdFlags 打印出日期和时间。
该方法还定义了如上一些同名方法。

func (l *Logger) Print(v ...interface{})
func (l *Logger) Printf(format string, v ...interface{})
func (l *Logger) Println(v ...interface{})
func (l *Logger) Fatal(v ...interface{})
func (l *Logger) Fatalf(format string, v ...interface{})
func (l *Logger) Fatalln(v ...interface{})
func (l *Logger) Panic(v ...interface{})
func (l *Logger) Panicf(format string, v ...interface{})
func (l *Logger) Panicln(v ...interface{})
func (l *Logger) Flags() int
func (l *Logger) Prefix() string
func (l *Logger) SetFlags(flag int)
func (l *Logger) SetPrefix(prefix string)
其中 “Print 、Panic 、Fatal ” 系列函数和之前介绍的一样，Flags和Prefix分别可以获得log.Logger当前的日志抬头和前缀。 SetFlags ，SetPrefix 则可以用来设置日志抬头和前缀。
使用实例
最后我看有log模块将debug日志打印到文件的实例。



package main
import (
    "log"
    "os"
)
func main(){
    fileName := "xxx_debug.log"
    logFile,err  := os.Create(fileName)
    defer logFile.Close()
    if err != nil {
        log.Fatalln("open file error !")
    }
    debugLog := log.New(logFile,"[Debug]",log.Llongfile)
    debugLog.Println("A debug message here")
    debugLog.SetPrefix("[Info]")
    debugLog.Println("A Info Message here ")
    debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)
    debugLog.Println("A different prefix")
}
运行后打开日志文件我们可以看到相应的日志内容