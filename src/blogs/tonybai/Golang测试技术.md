# 一、基础测试技术

## 1.测试GO代码

Go语言内置测试框架

内置的测试框架通过testing包以及go test命令来提供测试功能。

... 

## 2.表驱动测试

Golang的struct字面值(struct literals)语法让我们可以轻松写出表驱动测试。

package strings_test

import (
        "strings"
        "testing"
)

func TestIndex(t *testing.T) {
        var tests = []struct {
                s   string
                sep string
                out int
        }{
                {"", "", 0},
                {"", "a", -1},
                {"fo", "foo", -1},
                {"foo", "foo", 0},
                {"oofofoofooo", "f", 2},
                // etc
        }
        for _, test := range tests {
                actual := strings.Index(test.s, test.sep)
                if actual != test.out {
                        t.Errorf("Index(%q,%q) = %v; want %v",
                             test.s, test.sep, actual, test.out)
                }
        }
}

$go test -v strings_test.go
=== RUN TestIndex
— PASS: TestIndex (0.00 seconds)
PASS
ok      command-line-arguments    0.007s

## 3、T结构

*testing.T参数用于错误报告：

t.Errorf("got bar = %v, want %v", got, want)
t.Fatalf("Frobnicate(%v) returned error: %v", arg, err)
t.Logf("iteration %v", i)

也可以用于enable并行测试(parallet test)：
***t.Parallel()***

控制一个测试是否运行：

if runtime.GOARCH == "arm" {
    ***t.Skip***("this doesn't work on ARM")
}

## 4.运行测试
我们用go test命令来运行特定包的测试。

默认执行当前路径下包的测试代码。
go test 
go test -v 

要运行工程下的所有测试,我们执行如下命令:
go test github.com/nf/... 

标准库的测试:
go test std

## 5、测试覆盖率

go tool命令可以报告测试覆盖率统计。

我们在testgo下执行go test -cover，结果如下：

go build _/Users/tony/Test/Go/testgo: no buildable Go source files in /Users/tony/Test/Go/testgo
FAIL    _/Users/tony/Test/Go/testgo [build failed]

显然通过cover参数选项计算测试覆盖率不仅需要测试代码，还要有被测对象（一般是函数）的源码文件。

我们将目录切换到$GOROOT/src/pkg/strings下，执行go test -cover：

$go test -v -cover
=== RUN TestReader
— PASS: TestReader (0.00 seconds)
… …
=== RUN: ExampleTrimPrefix
— PASS: ExampleTrimPrefix (1.75us)
PASS
coverage: 96.9% of statements
ok      strings    0.612s

go test可以生成覆盖率的profile文件，这个文件可以被go tool cover工具解析。

在$GOROOT/src/pkg/strings下面执行：

$ go test -coverprofile=cover.out

会再当前目录下生成cover.out文件。

查看cover.out文件，有两种方法：

a) cover -func=cover.out

$sudo go tool cover -func=cover.out
strings/reader.go:24:    Len                66.7%
strings/reader.go:31:    Read                100.0%
strings/reader.go:44:    ReadAt                100.0%
strings/reader.go:59:    ReadByte            100.0%
strings/reader.go:69:    UnreadByte            100.0%
… …
strings/strings.go:638:    Replace                100.0%
strings/strings.go:674:    EqualFold            100.0%
total:            (statements)            96.9%

b) 可视化查看

执行go tool cover -html=cover.out命令，会在/tmp目录下生成目录coverxxxxxxx，比如/tmp/cover404256298。目录下有一个 coverage.html文件。用浏览器打开coverage.html，即可以可视化的查看代码的测试覆盖情况。

 
关于go tool的cover命令，我的go version go1.3 darwin/amd64默认并不自带，需要通过go get下载。

$sudo GOPATH=/Users/tony/Test/GoToolsProjects go get code.google.com/p/go.tools/cmd/cover

下载后，cover安装在$GOROOT/pkg/tool/darwin_amd64下面。


# 二、高级测试技术

## 1.一个例子程序
outyet是一个web服务，用于宣告某个特定Go版本是否已经打标签发布了。其获取方法：

go get github.com/golang/example/outyet

注：
go get执行后，cd $GOPATH/src/github.com/golang/example/outyet下，执行go run main.go。然后用浏览器打开http://localhost:8080即可访问该Web服务了。

## 2.测试Http客户端和服务端
net/http/httptest包提供了许多帮助函数,用于测试那些发送或处理Http请求的代码。

## 3.httptest.Server
httptest.Server在本地回环网口的一个系统选择的端口上listen。它常用于端到端的HTTP测试。

type Server struct{
    URL string
    Listener net.Listener
    TLS *tls.Config
    Config *http.Server
}

func NewServer(handler http.Handler)

func(*Server)Close()error

## 4.httptest.Server实战
下面代码创建了一个临时Http Server,返回简单的Hello应答:

ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, client")
    }))
    defer ts.Close()

    res, err := http.Get(ts.URL)
    if err != nil {
        log.Fatal(err)
    }

    greeting, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s", greeting)


## 5.httptest.ResponseRecorder

httptest.ResponseRecorder是http.ResponseWriter的一个实现,用来记录变化,用在测试的后续检视中。
type ResponseRecorder struct{
    Code int
    HeaderMap http.Header
    Body *bytes.Buffer
    Flushed bool
}

## 6.httptest.ResponseRecorder实战
向一个HTTP handler中传入一个ResponseRecorder,通过它我们可以来检视生成的应答。
handler := func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "something failed", http.StatusInternalServerError)
    }

    req, err := http.NewRequest("GET", "http://example.com/foo", nil)
    if err != nil {
        log.Fatal(err)
    }

    w := httptest.NewRecorder()
    handler(w, req)

    fmt.Printf("%d – %s", w.Code, w.Body.String())


## 7.竞争检测(race detection)
当两个goroutine并发访问同一个变量,且至少一个goroutine对变量进行写操作时,就会发生数据竞争(data race)。

go test -race mypkg
go run -race mysrc.go
go build -race  mycmd
go install -race mypkg


## 13.子进程测试

有些时候,你需要测试的是一个进程的行为,而不仅仅是一个函数。如:
func Crasher(){
    fmt.Println("Going down in flames!")
    os.Exit(1)
}

为了测试上面的代码,我们将测试程序本身作为一个子进程进行测试:
func TestCrasher(t *testing.T){
    if os.Geetenv("BE_CRASHER")=="1"{
        Crasher()
        return
    }
    cmd:=exec.Command(os.Args[0],"-test.run=TestCrasher")
    cmd.Env=append(os.Environ(),"BE_CRASHER=1")
    err:=cmd.Run()
    if e,ok:=err.(*exec.ExitError);ok && !e.Success(){
        return
    }
    t.Fatalf("process ran with err %v, want exit status 1", err)
}

