# 隐藏状态

下面通过一个例子来演示一下channel如何用来隐藏状态:

## 例子:唯一的ID服务

package main

import (
	"fmt"
)

func newUniqueIDService() <-chan string {
	id := make(chan string)
	go func() {
		var counter int64 = 0
		for {
			id <- fmt.Sprintf("%x", counter)
			counter += 1
		}
	}()
	return id
}

func main() {
	id := newUniqueIDService()
	for i := 0; i < 10; i++ {
		fmt.Println(<-id)
	}
}

/*
0
1
2
3
4
5
6
7
8
9
*/


newUniqueIDService通过一个channel与main goroutine关联，main goroutine无需知道uniqueid实现的细节以及当前状态，只需通过channel获得最新id即可。

# 默认情况
我想这里John Graham-Cumming主要是想告诉我们select的default分支的实践用法。

## 1.select for non-blocking receive
idle:=make(chan []byte,5) //用一个带缓冲的channel构造一个简单的队列

select{
    case b=<-idle: // 尝试从idle队列中读取
    ... 
    default: // 队列空,分配一个新的buffer
        makes+=1
        b=make([]byte,size)
}

## 2.select for non-blocking send
idle:=make(chan []byte,5) //用一个带缓冲的channel构造一个简单的队列

select{
    case idle<-b: //尝试向队列中插入一个buffer
    // ... 
    default: //队列满?
}

# Nil Channels

## 1.nil channels阻塞

对一个没有初始化的channel进行读写操作都将发生阻塞

package main 

func main(){
    var c chan int
    <-c 
}

$go run testnilchannel.go
fatal error: all goroutines are asleep – deadlock!

package main

func main() {
        var c chan int
        c <- 1
}

$go run testnilchannel.go
fatal error: all goroutines are asleep – deadlock!

## 2.nil channel在select中很有用

看下面这个例子：

//testnilchannel_bad.go
package main

import "fmt"
import "time"

func main() {
        var c1, c2 chan int = make(chan int), make(chan int)
        go func() {
                time.Sleep(time.Second * 5)
                c1 <- 5
                close(c1)
        }()

        go func() {
                time.Sleep(time.Second * 7)
                c2 <- 7
                close(c2)
        }()

        for {
                select {
                case x := <-c1:
                        fmt.Println(x)
                case x := <-c2:
                        fmt.Println(x)
                }
        }
        fmt.Println("over")
}

我们原本期望程序交替输出5和7两个数字，但实际的输出结果却是：

5
0
0
0
… … 0死循环

再仔细分析代码，原来select每次按case顺序evaluate：
    – 前5s，select一直阻塞；
    – 第5s，c1返回一个5后被close了，“case x := <-c1”这个分支返回，select输出5，并重新select
    – 下一轮select又从“case x := <-c1”这个分支开始evaluate，由于c1被close，按照前面的知识，close的channel不会阻塞，我们会读出这个 channel对应类型的零值，这里就是0；select再次输出0；这时即便c2有值返回，程序也不会走到c2这个分支
    – 依次类推，程序无限循环的输出0

我们利用nil channel来改进这个程序，以实现我们的意图，代码如下：

//testnilchannel.go
package main

import "fmt"
import "time"

func main() {
        var c1, c2 chan int = make(chan int), make(chan int)
        go func() {
                time.Sleep(time.Second * 5)
                c1 <- 5
                close(c1)
        }()

        go func() {
                time.Sleep(time.Second * 7)
                c2 <- 7
                close(c2)
        }()

        for {
                select {
                case x, ok := <-c1:
                        if !ok {
                                c1 = nil
                        } else {
                                fmt.Println(x)
                        }
                case x, ok := <-c2:
                        if !ok {
                                c2 = nil
                        } else {
                                fmt.Println(x)
                        }
                }
                if c1 == nil && c2 == nil {
                        break
                }
        }
        fmt.Println("over")
}

$go run testnilchannel.go
5
7
over

可以看出：通过将已经关闭的channel置为nil，下次select将会阻塞在该channel上，使得select继续下面的分支evaluation。



