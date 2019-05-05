
Go 1.8中值得关注的几个变化
二月 3, 2017 3 条评论
在已经过去的2016年，Go语言继在2009年之后再次成为编程语言界的明星- 问鼎TIOBE 2016年度语言。这与Go team、Go community和全世界的Gophers的努力是分不开的。按计划在这个2月份，Go team将正式发布Go 1.8版本(截至目前，Go的最新版本是Go 1.8rc3)。在这里我们一起来看一下在Go 1.8版本中都有哪些值得Gopher们关注的变化。
一、语言（Language）

Go 1.8版本依旧坚守Go Team之前的承诺，即Go1兼容性：使用Go 1.7及以前版本编写的Go代码，理论上都可以通过Go 1.8进行编译并运行。因此在臆想中的Go 2.0变成现实之前，每个Go Release版本在语言这方面的“改变”都会是十分微小的。
1、仅tags不同的两个struct可以相互做显式类型转换

在Go 1.8版本以前，两个struct即便字段个数相同且每个字段类型均一样，但如果某个字段的tag描述不一样，这两个struct相互间也不能做显式类型转换，比如：
//go18-examples/language/structtag.go
package main

import "fmt"

type XmlEventRegRequest struct {
    AppID     string `xml:"appid"`
    NeedReply int    `xml:"Reply,omitempty"`
}

type JsonEventRegRequest struct {
    AppID     string `json:"appid"`
    NeedReply int    `json:"reply,omitempty"`
}

func convert(in *XmlEventRegRequest) *JsonEventRegRequest {
    out := &JsonEventRegRequest{}
    *out = (JsonEventRegRequest)(*in)
    return out
}

func main() {
    in := XmlEventRegRequest{
        AppID:     "wx12345678",
        NeedReply: 1,
    }
    out := convert(&in)
    fmt.Println(out)
}

采用Go 1.7.4版本go compiler进行编译，我们会得到如下错误输出：
$go build structtag.go
# command-line-arguments
./structtag.go:17: cannot convert *in (type XmlEventRegRequest) to type JsonEventRegRequest
但在Go 1.8中，gc将忽略tag值的不同，使得显式类型转换成为可能：
$go run structtag.go
&{wx12345678 1}
改变虽小，但带来的便利却不小，否则针对上面代码中的convert，我们只能做逐一字段赋值了。
2、浮点常量的指数部分至少支持16bits长

在Go 1.8版本之前的The Go Programming Language Specificaton中，关于浮点数常量的指数部分的描述如下：
Represent floating-point constants, including the parts of a complex constant, with a mantissa of at least 256 bits and a signed exponent of at least 32 bits.
在Go 1.8版本中，文档中对于浮点数常量指数部分的长度的实现的条件放宽了，由支持最少32bit，放宽到最少支持16bits：
Represent floating-point constants, including the parts of a complex constant, with a mantissa of at least 256 bits and a signed binary exponent of at least 16 bits.
但Go 1.8版本go compiler实际仍然支持至少32bits的指数部分长度，因此这个改变对现存的所有Go源码不会造成影响。
二、标准库（Standard Library）

Go号称是一门”Batteries Included”编程语言。“Batteries Included”指的就是Go语言强大的标准库。使用Go标准库，你可以完成绝大部分你想要的功能，而无需再使用第三方库。Go语言的每次版本更新，都会在标准库环节增加强大的功能、提升性能或是提高使用上的便利性。每次版本更新，标准库也是改动最大的部分。这次也不例外，我们逐一来看。
1、便于slice sort的sort.Slice函数

在Go 1.8之前我们要对一个slice进行sort，需要定义出实现了下面接口的slice type：
//$GOROOT/src/sort.go
... ...
type Interface interface {
    // Len is the number of elements in the collection.
    Len() int
    // Less reports whether the element with
    // index i should sort before the element with index j.
    Less(i, j int) bool
    // Swap swaps the elements with indexes i and j.
    Swap(i, j int)
}

标准库定义了一些应对常见类型slice的sort类型以及对应的函数：
StringSlice -> sort.Strings
IntSlice -> sort.Ints
Float64Slice -> sort.Float64s
但即便如此，对于用户定义的struct或其他自定义类型的slice进行排序仍需定义一个新type，比如下面这个例子中的TiboeIndexByRank：
//go18-examples/stdlib/sort/sortslice-before-go18.go
package main

import (
    "fmt"
    "sort"
)

type Lang struct {
    Name string
    Rank int
}

type TiboeIndexByRank []Lang

func (l TiboeIndexByRank) Len() int           { return len(l) }
func (l TiboeIndexByRank) Less(i, j int) bool { return l[i].Rank < l[j].Rank }
func (l TiboeIndexByRank) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func main() {
    langs := []Lang{
        {"rust", 2},
        {"go", 1},
        {"swift", 3},
    }
    sort.Sort(TiboeIndexByRank(langs))
    fmt.Printf("%v\n", langs)
}

$go run sortslice-before-go18.go
[{go 1} {rust 2} {swift 3}]

从上面的例子可以看到，我们要对[]Lang这个slice进行排序，我们就需要为之定义一个专门用于排序的类型：这里是TiboeIndexByRank，并让其实现sort.Interface接口。使用过sort包的gophers们可能都意识到了，我们在为新的slice type实现sort.Interface接口时，那三个方法的Body几乎每次都是一样的。为了使得gopher们在排序slice时编码更为简化和便捷，减少copy&paste，Go 1.8为slice type新增了三个函数：Slice、SliceStable和SliceIsSorted。我们重新用Go 1.8的sort.Slice函数实现上面例子中的排序需求，代码如下：
//go18-examples/stdlib/sort/sortslice-in-go18.go
package main

import (
    "fmt"
    "sort"
)

type Lang struct {
    Name string
    Rank int
}

func main() {
    langs := []Lang{
        {"rust", 2},
        {"go", 1},
        {"swift", 3},
    }
    sort.Slice(langs, func(i, j int) bool { return langs[i].Rank < langs[j].Rank })
    fmt.Printf("%v\n", langs)
}

$go run sortslice-in-go18.go
[{go 1} {rust 2} {swift 3}]
实现sort，需要三要素：Len、Swap和Less。在1.8之前，我们通过实现sort.Interface实现了这三个要素；而在1.8版本里，Slice函数通过reflect获取到swap和length，通过结合闭包实现的less参数让Less要素也具备了。我们从下面sort.Slice的源码可以看出这一点：
// $GOROOT/src/sort/sort.go
... ...
func Slice(slice interface{}, less func(i, j int) bool) {
    rv := reflect.ValueOf(slice)
    swap := reflect.Swapper(slice)
    length := rv.Len()
    quickSort_func(lessSwap{less, swap}, 0, length, maxDepth(length))
}

2、支持HTTP/2 Push

继在Go 1.6版本全面支持HTTP/2之后，Go 1.8又新增了对HTTP/2 Push的支持。HTTP/2是在HTTPS的基础上的下一代HTTP协议，虽然当前HTTPS的应用尚不是十分广泛。而HTTP/2 Push是HTTP/2的一个重要特性，无疑其提出的初衷也仍然是为了改善网络传输性能，提高Web服务的用户侧体验。这里我们可以借用知名网络提供商Cloudflare blog上的一幅示意图来诠释HTTP/2 Push究竟是什么：

从上图中，我们可以看到：当Browser向Server发起Get page.html请求后，在同一条TCP Connection上，Server主动将style.css和image.png两个资源文件推送(Push)给了Browser。这是由于Server端启用了HTTP/2 Push机制，并预测判断Browser很可能会在接下来发起Get style.css和image.png两个资源的请求。这是一种典型的：“你可能会需要，但即使你不要，我也推给你”的处世哲学^0^。这种机制虽然在一定程度上能改善网络传输性能（减少Client发起Get的次数），但也可能造成带宽的浪费，因为这些主动推送给Browser的资源很可能是Browser所不需要的或是已经在Browser cache中存在的资源。
接下来，我们来看看Go 1.8是如何在net/http包中提供对HTTP/2 Push的支持的。由于HTTP/2是基于HTTPS的，因此我们先使用generate_cert.go生成程序所需的私钥和证书：
// 在go18-examples/stdlib/http2-push目录下，执行：

$go run $GOROOT/src/crypto/tls/generate_cert.go --host 127.0.0.1
2017/01/27 10:58:01 written cert.pem
2017/01/27 10:58:01 written key.pem
支持HTTP/2 Push的server端代码如下：
// go18-examples/stdlib/http2-push/server.go

package main

import (
    "fmt"
    "log"
    "net/http"
)

const mainJS = `document.write('Hello World!');`

func main() {
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        pusher, ok := w.(http.Pusher)
        if ok {
            // If it's a HTTP/2 Server.
            // Push is supported. Try pushing rather than waiting for the browser.
            if err := pusher.Push("/static/img/gopherizeme.png", nil); err != nil {
                log.Printf("Failed to push: %v", err)
            }
        }
        fmt.Fprintf(w, `<html>
<head>
<title>Hello Go 1.8</title>
</head>
<body>
    <img src="/static/img/gopherizeme.png"></img>
</body>
</html>
`)
    })
    log.Fatal(http.ListenAndServeTLS(":8080", "./cert.pem", "./key.pem", nil))
}
运行这段代码，打开Google Chrome浏览器，输入：https://127.0.0.1:8080，忽略浏览器的访问非受信网站的警告，继续浏览你就能看到下面的页面（这里打开了Chrome的“检查”功能）：

从示例图中的“检查”窗口，我们可以看到gopherizeme.png这个image资源就是Server主动推送给客户端的，这样浏览器在Get /后无需再发起一次Get /static/img/gopherizeme.png的请求了。
而这一切的背后，其实是HTTP/2的ResponseWriter实现了Go 1.8新增的http.Pusher interface：
// $GOROOT/src/net/http/http.go

// Pusher is the interface implemented by ResponseWriters that support
// HTTP/2 server push. For more background, see
// https://tools.ietf.org/html/rfc7540#section-8.2.
type Pusher interface {
    ... ...
    Push(target string, opts *PushOptions) error
}
3、支持HTTP Server优雅退出

Go 1.8中增加对HTTP Server优雅退出(gracefullly exit)的支持，对应的新增方法为：
func (srv *Server) Shutdown(ctx context.Context) error
和server.Close在调用时瞬间关闭所有active的Listeners和所有状态为New、Active或idle的connections不同，server.Shutdown首先关闭所有active Listeners和所有处于idle状态的Connections，然后无限等待那些处于active状态的connection变为idle状态后，关闭它们并server退出。如果有一个connection依然处于active状态，那么server将一直block在那里。因此Shutdown接受一个context参数，调用者可以通过context传入一个Shutdown等待的超时时间。一旦超时，Shutdown将直接返回。对于仍然处理active状态的Connection，就任其自生自灭（通常是进程退出后，自动关闭）。通过Shutdown的源码我们也可以看出大致的原理：
// $GOROOT/src/net/http/server.go
... ...
func (srv *Server) Shutdown(ctx context.Context) error {
    atomic.AddInt32(&srv.inShutdown, 1)
    defer atomic.AddInt32(&srv.inShutdown, -1)

    srv.mu.Lock()
    lnerr := srv.closeListenersLocked()
    srv.closeDoneChanLocked()
    srv.mu.Unlock()

    ticker := time.NewTicker(shutdownPollInterval)
    defer ticker.Stop()
    for {
        if srv.closeIdleConns() {
            return lnerr
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
        }
    }
}

我们来编写一个例子：
// go18-examples/stdlib/graceful/server.go

import (
    "context"
    "io"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
)

func main() {
    exit := make(chan os.Signal)
    signal.Notify(exit, os.Interrupt)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Println("Handle a new request:", *r)
        time.Sleep(10 * time.Second)
        log.Println("Handle the request ok!")
        io.WriteString(w, "Finished!")
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Printf("listen: %s\n", err)
        }
    }()

    <-exit // wait for SIGINT
    log.Println("Shutting down server...")

    // Wait no longer than 30 seconds before halting
    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
    err := srv.Shutdown(ctx)

    log.Println("Server gracefully stopped:", err)
}

在上述例子中，我们通过设置Linux Signal的处理函数来拦截Linux Interrupt信号并处理。我们通过context给Shutdown传入30s的超时参数，这样Shutdown在退出之前会给各个Active connections 30s的退出时间。下面分为几种情况run一下这个例子：
a)  当前无active connections
在这种情况下，我们run上述demo，ctrl + C后，上述demo直接退出：
$go run server.go
^C2017/02/02 15:13:16 Shutting down server...
2017/02/02 15:13:16 Server gracefully stopped: <nil>
b) 当前有未处理完的active connections，ctx 超时
为了模拟这一情况，我们修改一下参数。让每个request handler的sleep时间为30s，而Shutdown ctx的超时时间改为10s。我们再来运行这个demo，并通过curl命令连接该server(curl  -v http://localhost:8080)，待连接成功后，再立即ctrl+c停止Server，待约10s后，我们得到如下日志：
$go run server.go
2017/02/02 15:15:57 Handle a new request: {GET / HTTP/1.1 1 1 map[User-Agent:[curl/7.30.0] Accept:[*/*]] {} <nil> 0 [] false localhost:8080 map[] map[] <nil> map[] [::1]:52590 / <nil> <nil> <nil> 0xc420016700}
^C2017/02/02 15:15:59 Shutting down server...
2017/02/02 15:15:59 listen: http: Server closed
2017/02/02 15:16:09 Server gracefully stopped: context deadline exceeded
c) 当前有未处理完的active connections，ctx超时之前，这些connections处理ok了
我们将上述demo的参数还原，即request handler sleep 10s，而Shutdown ctx超时时间为30s，运行这个Demo后，通过curl命令连接该server，待连接成功后，再立即ctrl+c停止Server。等待约10s后，我们得到如下日志：
$go run server.go
2017/02/02 15:19:56 Handle a new request: {GET / HTTP/1.1 1 1 map[User-Agent:[curl/7.30.0] Accept:[*/*]] {} <nil> 0 [] false localhost:8080 map[] map[] <nil> map[] [::1]:52605 / <nil> <nil> <nil> 0xc420078500}
^C2017/02/02 15:19:59 Shutting down server...
2017/02/02 15:19:59 listen: http: Server closed
2017/02/02 15:20:06 Handle the request ok!
2017/02/02 15:20:06 Server gracefully stopped: <nil>
可以看出，当ctx超时之前，request处理ok，connection关闭。这时不再有active connection和idle connection了，Shutdown成功返回，server立即退出。
4、Mutex Contention Profiling

Go 1.8中runtime新增了对Mutex和RWMutex的profiling(剖析)支持。golang team成员，负责从go user角度去看待go team的work是否满足用户需求的Jaana B. Dogan在其个人站点上写了一篇介绍mutex profiling的文章，这里借用一下其中的Demo：
//go18-examples/stdlib/mutexprofile/mutexprofile.go

package main

import (
    "net/http"
    _ "net/http/pprof"
    "runtime"
    "sync"
)

func main() {
    var mu sync.Mutex
    var items = make(map[int]struct{})

    runtime.SetMutexProfileFraction(5)
    for i := 0; i < 1000*1000; i++ {
        go func(i int) {
            mu.Lock()
            defer mu.Unlock()
            items[i] = struct{}{}
        }(i)
    }

    http.ListenAndServe(":8888", nil)
}
运行该程序后，在浏览器中输入：http://localhost:8888/debug/pprof/mutex，你就可以看到有关该程序的mutex profile（耐心等待一小会儿，因为数据的采样需要一点点时间^0^）：
--- mutex:
cycles/second=2000012082
sampling period=5
378803564 776 @ 0x106c4d1 0x13112ab 0x1059991
构建该程序，然后通过下面命令：
go build mutexprofile.go
./mutexprofile
go tool pprof mutexprofile http://localhost:8888/debug/pprof/mutex?debug=1
可以进入pprof交互界面，这个是所有用过go pprof工具的gophers们所熟知的：
$go tool pprof mutexprofile http://localhost:8888/debug/pprof/mutex?debug=1
Fetching profile from http://localhost:8888/debug/pprof/mutex?debug=1
Saved profile in /Users/tony/pprof/pprof.mutexprofile.localhost:8888.contentions.delay.003.pb.gz
Entering interactive mode (type "help" for commands)
(pprof) list
Total: 12.98s
ROUTINE ======================== main.main.func1 in /Users/tony/Test/GoToolsProjects/src/github.com/bigwhite/experiments/go18-examples/stdlib/mutexprofile/mutexprofile.go
         0     12.98s (flat, cum)   100% of Total
         .          .     17:            mu.Lock()
         .          .     18:            defer mu.Unlock()
         .          .     19:            items[i] = struct{}{}
         .          .     20:        }(i)
         .          .     21:    }
         .     12.98s     22:
         .          .     23:    http.ListenAndServe(":8888", nil)
         .          .     24:}
ROUTINE ======================== runtime.goexit in /Users/tony/.bin/go18rc2/src/runtime/asm_amd64.s
         0     12.98s (flat, cum)   100% of Total
         .          .   2192:    RET
         .          .   2193:
         .          .   2194:// The top-most function running on a goroutine
         .          .   2195:// returns to goexit+PCQuantum.
         .          .   2196:TEXT runtime·goexit(SB),NOSPLIT,$0-0
         .     12.98s   2197:    BYTE    $0x90    // NOP
         .          .   2198:    CALL    runtime·goexit1(SB)    // does not return
         .          .   2199:    // traceback from goexit1 must hit code range of goexit
         .          .   2200:    BYTE    $0x90    // NOP
         .          .   2201:
         .          .   2202:TEXT runtime·prefetcht0(SB),NOSPLIT,$0-8
ROUTINE ======================== sync.(*Mutex).Unlock in /Users/tony/.bin/go18rc2/src/sync/mutex.go
    12.98s     12.98s (flat, cum)   100% of Total
         .          .    121:            return
         .          .    122:        }
         .          .    123:        // Grab the right to wake someone.
         .          .    124:        new = (old - 1<<mutexWaiterShift) | mutexWoken
         .          .    125:        if atomic.CompareAndSwapInt32(&m.state, old, new) {
    12.98s     12.98s    126:            runtime_Semrelease(&m.sema)
         .          .    127:            return
         .          .    128:        }
         .          .    129:        old = m.state
         .          .    130:    }
         .          .    131:}
(pprof) top10
1.29s of 1.29s total (  100%)
      flat  flat%   sum%        cum   cum%
     1.29s   100%   100%      1.29s   100%  sync.(*Mutex).Unlock
         0     0%   100%      1.29s   100%  main.main.func1
         0     0%   100%      1.29s   100%  runtime.goexit

go pprof的另外一个用法就是在go test时，mutexprofile同样支持这一点：
go test -mutexprofile=mutex.out
go tool pprof <test.binary> mutex.out
5、其他重要改动

Go 1.8标准库还有两个值得注意的改动，一个是：crypto/tls，另一个是database/sql。
在HTTPS逐渐成为主流的今天，各个编程语言对HTTPS连接的底层加密协议- TLS协议支持的成熟度日益被人们所关注。Go 1.8给广大Gophers们带来了一个更为成熟、性能更好、更为安全的TLS实现，同时也增加了对一些TLS领域最新协议规范的支持。无论你是实现TLS Server端，还是Client端，都将从中获益。
Go 1.8在crypto/tls中提供了基于ChaCha20-Poly1305的cipher suite，其中ChaCha20是一种stream cipher算法；而Poly1305则是一种code authenticator算法。它们共同组成一个TLS suite。使用这个suite，将使得你的web service或站点具有更好的mobile浏览性能，这是因为传统的AES算法实现在没有硬件支持的情况下cost更多。因此，如果你在使用tls时没有指定cipher suite，那么Go 1.8会根据硬件支持情况（是否有AES的硬件支持），来决定是使用ChaCha20还是AES算法。除此之外，crypto/tls还实现了更为安全和高效的X25519密钥交换算法等。
Go 1.4以来，database/sql包的变化很小，但对于该包的feature需求却在与日俱增。终于在Go 1.8这个dev cycle中，govendor的作者Daniel Theophanes在Brad Fitzpatrick的“指导”下，开始对database/sql进行“大规模”的改善。在Go 1.8中，借助于context.Context的帮助，database/sql增加了Cancelable Queries、SQL Database Type、Multiple Result Sets、Database ping、Named Parameters和Transaction Isolation等新Features。在GopherAcademy的Advent 2016系列文章中，我们可以看到Daniel Theophanes亲手撰写的文章，文章针对Go 1.8 database/sql包新增的features作了详细解释。
三、Go工具链（Go Toolchain）

在目前市面上的主流编程语言中，如果说Go的工具链在成熟度和完善度方面排第二，那没有语言敢称自己是第一吧^_^。Go 1.8在Go Toolchain上继续做着持续地改进，下面我们来逐一看看。
1、Plugins

Go在1.8版本中提供了对Plugin的初步支持，并且这种支持仅限于Linux。plugin这个术语在不同语言、不同情景上下文中有着不同的含义，那么什么是Go Plugin呢？
Go Plugin为Go程序提供了一种在运行时加载代码、执行代码以改变运行行为的能力，它实质上由两个部分组成：
* go build -buildmode=plugin xx.go 构建xx.so plugin文件
* 利用plugin包在运行时动态加载xx.so并执行xx.so中的代码
C程序员看到这里肯定会有似曾相识的赶脚，因为这和传统的动态共享库在概念上十分类似：
go build -buildmode=plugin xx.go 类似于 gcc -o xx.so -shared xx.c
go plugin包 类似于 linux上的dlopen/dlsym或windows上的LoadLibrary
我们来看一个例子！我们先来建立一个名为foo.so的go plugin：
//go18-examples/gotoolchain/plugins/foo.go

package main

import "fmt"

var V int
var v int

func init() {
        V = 17
        v = 23
        fmt.Println("init function in plugin foo")
}

func Foo(in string) string {
        return "Hello, " + in
}

func foo(in string) string {
        return "hello, " + in
}

通过go build命令将foo.go编译为foo.so：
# go build -buildmode=plugin foo.go
# ldd foo.so
    linux-vdso.so.1 =>  (0x00007ffe47f67000)
    libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f9d06f4b000)
    libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f9d06b82000)
    /lib64/ld-linux-x86-64.so.2 (0x000055c69cfcf000)

# nm foo.so|grep Foo
0000000000150010 t local.plugin/unnamed-69e21ef38d16a3fee5eb7b9e515c27a389067879.Foo
0000000000150010 T plugin/unnamed-69e21ef38d16a3fee5eb7b9e515c27a389067879.Foo
000000000036a0dc D type..namedata.Foo.

我们看到go plugin的.so文件就是一个标准的Linux动态共享库文件，我们可以通过nm命令查看.so中定义的各种符号。接下来，我们来load这个.so，并查找并调用相应符号：
//go18-examples/gotoolchain/plugins/main.go

package main

import (
        "fmt"
        "plugin"
        "time"
)

func init() {
        fmt.Println("init in main program")
}

func loadPlugin(i int) {
        fmt.Println("load plugin #", i)
        var err error
        fmt.Println("before opening the foo.so")

        p, err := plugin.Open("foo.so")
        if err != nil {
                fmt.Println("plugin Open error:", err)
                return
        }
        fmt.Println("after opening the foo.so")

        f, err := p.Lookup("Foo")
        if err != nil {
                fmt.Println("plugin Lookup symbol Foo error:", err)
        } else {
                fmt.Println(f.(func(string) string)("gophers"))
        }

        f, err = p.Lookup("foo")
        if err != nil {
                fmt.Println("plugin Lookup symbol foo error:", err)
        } else {
                fmt.Println(f.(func(string) string)("gophers"))
        }

        v, err := p.Lookup("V")
        if err != nil {
                fmt.Println("plugin Lookup symbol V error:", err)
        } else {
                fmt.Println(*v.(*int))
        }

        v, err = p.Lookup("v")
        if err != nil {
                fmt.Println("plugin Lookup symbol v error:", err)
        } else {
                fmt.Println(*v.(*int))
        }
        fmt.Println("load plugin #", i, "done")
}

func main() {
        var counter int = 1
        for {
                loadPlugin(counter)
                counter++
                time.Sleep(time.Second * 30)
        }
}

执行这个程序：
# go run main.go
init in main program
load plugin # 1
before opening the foo.so
init function in plugin foo
after opening the foo.so
Hello, gophers
plugin Lookup symbol foo error: plugin: symbol foo not found in plugin plugin/unnamed-69e21ef38d16a3fee5eb7b9e515c27a389067879
17
plugin Lookup symbol v error: plugin: symbol v not found in plugin plugin/unnamed-69e21ef38d16a3fee5eb7b9e515c27a389067879
load plugin # 1 done

load plugin # 2
before opening the foo.so
after opening the foo.so
Hello, gophers
plugin Lookup symbol foo error: plugin: symbol foo not found in plugin plugin/unnamed-69e21ef38d16a3fee5eb7b9e515c27a389067879
17
plugin Lookup symbol v error: plugin: symbol v not found in plugin plugin/unnamed-69e21ef38d16a3fee5eb7b9e515c27a389067879
load plugin # 2 done
... ...
我们来分析一下这个执行结果！
a) foo.go中的代码也包含在main package下，但只是当foo.so被第一次加载时，foo.go中的init函数才会被执行；
b) foo.go中的exported function和variable才能被Lookup到，如Foo、V；查找unexported的变量和函数符号将得到error信息，如：“symbol foo not found in plugin”；
c) Lookup返回的是plugin.Symbol类型的值，plugin.Symbol是一个指向plugin中变量或函数的指针；
d) foo.go中的init在后续重复加载中并不会被执行。
注意：plugin.Lookup是goroutine-safe的。
在golang-dev group上，有人曾问过：buildmode=c-shared和buildmode=plugin有何差别？Go team member给出的答案如下：
The difference is mainly on the program that loads the shared library.

For c-shared, we can't assume anything about the host, so the c-shared dynamic library must be self-contained, but for plugin, we know the host program will be a Go program built with the same runtime version, so the toolchain can omit at least the runtime package from the dynamic library, and possibly more if it's certain that some packages are linked into the host program. (This optimization hasn't be implemented yet, but we need the distinction to enable this kind of optimization in the future.)
2、默认的GOPATH

Go team在Go 1.8以及后续版本会更加注重”Go语言的亲民性”，即进一步降低Go的入门使用门槛，让大家更加Happy的使用Go。对于一个Go初学者来说，一上来就进行GOPATH的设置很可能让其感到有些迷惑，甚至有挫折感，就像建立Java开发环境需要设置JAVA_HOME和CLASSPATH一样。Gophers们期望能做到Go的安装即可用。因此Go 1.8就在这方面做出了改进：支持默认的GOPATH。
在Linux/Mac系下，默认的GOPATH为$HOME/go，在Windows下，GOPATH默认路径为：%USERPROFILE%/go。你可以通过下面命令查看到这一结果：
$ go env
GOARCH="amd64"
GOBIN="/home/tonybai/.bin/go18rc3/bin"
GOEXE=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOOS="linux"
GOPATH="/home/tonybai/go"
GORACE=""
GOROOT="/home/tonybai/.bin/go18rc3"
GOTOOLDIR="/home/tonybai/.bin/go18rc3/pkg/tool/linux_amd64"
GCCGO="gccgo"
CC="gcc"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build313929093=/tmp/go-build -gno-record-gcc-switches"
CXX="g++"
CGO_ENABLED="1"
PKG_CONFIG="pkg-config"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"

BTW，在Linux/Mac下，默认的GOROOT为/usr/local/go，如果你的Go环境没有安装到这个路径下，在没有设置$GOROOT环境变量的情况下，当你执行go subcommand相关命令时，你会看到如下错误：
$go env
go: cannot find GOROOT directory: /usr/local/go
3、其他变化

Go 1.8删除了Go 1.7中增加的用于关闭ssa新后端的”-ssa=0” compiler flag，并且将ssa backend扩展到所有architecture中，对ssa后端也进一步做了优化。与此同时，为了将来进一步的性能优化打基础，Go 1.8还引入了一个新编译器前端，当然这对于普通Gopher的Go使用并没有什么影响。
Go 1.8还新增go bug子命令，该命令会自动使用默认浏览器打开new issue页面，并将采集到的issue提交者的系统信息填入issue模板，以帮助gopher提交符合要求的go issue，下面是go bug打开的issue page的图示：

四、性能变化（Performance Improvement）

无论是Gotoolchain、还是runtime（包括GC）的性能，一直都是Go team重点关注的领域。本次Go 1.8依旧给广大Gophers们带来了性能提升方面的惊喜。
首先，Go SSA后端扩展到所有architecture和新编译器前端的引入，将会给除X86-64之外架构上运行的Go代码带来约20-30%的运行性能提升。对于x86-64，虽然Go 1.7就已经开启了SSA，但Go 1.8对SSA做了进一步优化，x86-64上的Go代码依旧可能会得到10%以内的性能提升。
其次，Go 1.8持续对Go compiler和linker做性能优化，和1.7相比，平均编译链接的性能提升幅度在15%左右。虽然依旧没有达到Go 1.4的性能水准。不过，优化依旧在持续进行中，目标的达成是可期的。
再次，GC在低延迟方面的优化给了我们最大的惊喜。在Go 1.8中，由于消除了GC的“stop-the-world stack re-scanning”，使得GC STW(stop-the-world)的时间通常低于100微秒，甚至经常低于10微秒。当然这或多或少是以牺牲“吞吐”作为代价的。因此在Go 1.9中，GC的改进将持续进行，会在吞吐和低延迟上做一个很好的平衡。
最后，defer的性能消耗在Go 1.8中下降了一半，与此下降幅度相同的还有通过cgo在go中调用C代码的性能消耗。
五、小结兼参考资料

Go 1.8的变化不仅仅是以上这些，更多变化以及详细的描述请参考下面参考资料中的“Go 1.8 Release Notes”：
* Go 1.8
* Go 1.8 Release Notes
* Go 1.8 Release Notes(before release)
* Announcing Support for HTTP/2 Server Push
* What’s coming in Go 1.8
以上demo中的代码在这里可以找到。
© 2017, bigwhite. 版权所有. 
Related posts:
1. Go 1.7中值得关注的几个变化
2. Go 1.6中值得关注的几个变化
3. 使用Golang开发微信公众平台-发送客服消息
4. Go语言错误处理
5. 使用Golang开发微信公众平台-接收加密消息