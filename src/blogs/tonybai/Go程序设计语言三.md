# 多路复用
channel是原生值,这意味着他们也能通过channel发送。这个属性使得编写一个服务类多路复用器变得十分容易,因为客户端在提交请求时可一并提供用于回复应答的channel。

chanOfChans:=make(chan chan int)

## 或者更典型如:

type Reply struct{...}
type Request struct{
    arg1,arg2 someType
    replyc chan *Reply
}

## 多路复用服务器

type request struct{
    a,b int 
    replyc chan int
}

type binOp func(a,b int)int 
func run(op binOp,req *request){
    req.replyc<-op(req.a,req.b)
}

func server(op binOp,service<-chan *request){
    for{
        req:=<-service //请求到达这里
        go run(op,req) // 不等op
    }
}

## 启动服务器
使用"返回channel的函数"惯用法来为一个新服务器创建一个channel:
func startServer(op binOp) chan<- *request{
    service:=make(chan *request)
    go server(op,service)
    return service
}

adderChan:=startServer(
    func(a,b int)int{return a+b }
)

## 客户端
在教程中有个例子更为详尽,但这里是一个变体:
func(r *request)String()string{
    return fmt.Spintf("%d+%d=%d",r.a,r.b,<-r.replyc)
}

req1:=&request{7,8,make(chan int)}
req2:=&request{17,18,make(chan int)}

请求已经就绪,发送它们:
adderChan<-req1
adderChan<-req2

可以以任何顺序获得结果;r.replyc多路分解:

fmt.Println(req2,req1)

## 停掉
在多路复用的例子中,服务将永远运行下去。要将其干净地停掉,可通过一个channel发送信号。下面这个server具有相同的功能,但多了一个quit channel:

func server(op binOp,service<-chan *request,quit<-chan bool){
    for{
        select{
            case req:=<-service:
                go run(op,req) // don`t wait for it
            case <- quit:
                return
        }
    }
}

## 启动服务器:
其他代码都相似,只是多了一个channel:
func startSerever(op binOp)(service chan<-*request,quit chan<-bool){
    service=make(chan *request)
    quit=make(chan bool)
    go server(op,service,quit)
    return service,quit
}

adderChan,quitChan:=startServer(
    func(a,b int)int{return a+b }
)

## 停掉:客户端
只有当准备停掉服务端的时候,客户端才会受到影响:

req1 := &request{7, 8, make(chan int)}
req2 := &request{17, 18, make(chan int)}
adderChan <- req1
adderChan <- req2
fmt.Println(req2, req1)

所有都完成后,向服务器发送信号,让其退出:
quitChan<-true


# 链？？？？？？？？？
package main

import (
	"flag"
	"fmt"
)

var nGoroutine = flag.Int("n", 100000, "how many")

func f(left, right chan int) { left <- 1 + <-right }
func main() {
	flag.Parse()
	leftmost := make(chan int)
	var left, right chan int = nil, leftmost
	for i := 0; i < *nGoroutine; i++ {
		left, right = right, make(chan int)
		go f(left, right)
	}
	right <- 0      // bang!
	x := <-leftmost // 等待完成
	fmt.Println(x)  // 100000
}

# 例子:channel作为缓存
var freeList=make(chan *Buffer,100)
var serverChan=make(chan *Buffer)

func server(){
    for{
        b:=<-serverChan // 等待做work
        process(b) // 在缓冲中处理请求
        select{
            case freeList <-b: // 如果有空间,重用缓存
            default:           // 否则,丢弃它
        }
    }
}

func client(){
    for{
        var b *Buffer
        select{
            case b=<-freeList: // 如果就绪,抓取一个
            default: b=new(Buffer) // 否则,分配一个
        }
        load(b) // 读取下一个请求放入b中
        serverChan<-b // 将请求发送给server.
    }
}
