package main 

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
