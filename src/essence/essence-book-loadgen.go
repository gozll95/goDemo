#book-loadgen

1.

请求
响应
包装之后的响应

"流程":
请求-->xx软件-->响应--->包装之后的响应---->chan 结果通道

要考虑到并发

将
请求-->xx软件-->响应--->包装之后的响应---->chan 结果通道
作为一次asnycDo

再考虑到每隔xx执行这一次asyncDo以及考虑到负载均衡器的自动停止

自动停止一般使用ctx+cancelFunc

那么一个负载均衡器开始的时候:
- 检查是否停止:
	- 停止了 -> 直接return
	- 没有 -> asyncDo ->使用ticker->case ticker/case 停止了

所以整个就是几个要点:
负载均衡器状态的变更
负载均衡器的async do的执行(使用go票池)

这里还要考虑到xx软件的超时导致结果的变化:
- 没超时 - > 结果1
- 超时了 - > 结果2

一般我们判断是否超时:
go func(){
	xxx执行xx
	c<-xxxxx
}()
select{
case <- timeout:
case <-c: 
}

但是由于我们将这个过程放置到了asyncDo里,所以我们尽量避免再在goroutine下开一个goroutine。
所以我们使用
timer.AfterFunc()+uint32来替换以上:

timer := time.AfterFunc(gen.timeoutNS, func() {
        if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
                return
        }
        result := &lib.CallResult{
                ID:     rawReq.ID,
                Req:    rawReq,
                Code:   lib.RET_CODE_WARNING_CALL_TIMEOUT,
                Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNS),
                Elapse: gen.timeoutNS,
        }
        gen.sendResult(result)
})

rawResp:=gen.callOne(&rawReq)
if !atomic.CompareAndSwapUint32(&callStatus,0,1){
    return
}
timer.Stop()

2."数据结构":


```
resultCh chan *lib.CallResult //调用结果通道
```

type CallResult struct{
    ID int64 //ID
    Req RawReq //原生请求
    Resp RawResp //原生响应
    Code RetCode //响应代码
    Msg string  //结果成因的简述
    Elapse time.Duration //耗时
}

// 用于表示原生请求的结构
type RawReq struct{
    ID int64
    Req []byte
}

// 用于表示原生响应的结构
type RawResp struct{
    ID int64
    Resp []byte
    Err error 
    Elapse time.Duration
}


// 声明代表载荷发生器状态的常量。
const (
        // STATUS_ORIGINAL 代表原始。
        STATUS_ORIGINAL uint32 = 0
        // STATUS_STARTING 代表正在启动。
        STATUS_STARTING uint32 = 1
        // STATUS_STARTED 代表已启动。
        STATUS_STARTED uint32 = 2
        // STATUS_STOPPING 代表正在停止。
        STATUS_STOPPING uint32 = 3
        // STATUS_STOPPED 代表已停止。
        STATUS_STOPPED uint32 = 4
)


// 用于表示调用器的接口
type Caller interface{
    //构建请求
    BuildReq()，
    //调用
    Call(req []byte,timeoutNS time.Duration)([]byte,error)
    //检查响应
    CheckResp(rawReq RawReq,rawResp RawResp)*CallResult
}

完整代码:
// myGenerator 代表载荷发生器的实现类型。
type myGenerator struct {
        caller      lib.Caller           // 调用器。
        timeoutNS   time.Duration        // 处理超时时间，单位：纳秒。
        lps         uint32               // 每秒载荷量。
        durationNS  time.Duration        // 负载持续时间，单位：纳秒。
        concurrency uint32               // 载荷并发量。
        tickets     lib.GoTickets        // Goroutine票池。
        ctx         context.Context      // 上下文。
        cancelFunc  context.CancelFunc   // 取消函数。
        callCount   int64                // 调用计数。
        status      uint32               // 状态。
        resultCh    chan *lib.CallResult // 调用结果通道。
}


//用于表示载荷发生器的接口
type Generator interface {
        // 启动载荷发生器。
        // 结果值代表是否已成功启动。
        Start() bool
        // 停止载荷发生器。
        // 结果值代表是否已成功停止。
        Stop() bool
        // 获取状态。
        Status() uint32
        // 获取调用计数。每次启动会重置该计数。
        CallCount() int64
}

3.:
status
ctx
cancelFunc
go票池
result chan *xx通道


4."测试":
一个TCP server
type TCPServer struct{
        listener net.Listener
        active uint32 //0-未激活;1-已激活
}

methods:
        - NewTCPServer ..  
        - Listener(xx地址)
                // 使用CAS初始化状态
                - init(xx地址)
                - go func(){
                        - for{
                                conn:=accept()
                                go handle with conn
                                        - read from conn 
                                        - handle 
                                        - write to conn
                        }
                }()

一个client:
        参数:TCPServer/addr
methods:
        - buildReq 
        - Call
                - conn,err:=net.DialTimeout(xx地址)
                - write(conn,req,DELIM)
                - return read(conn,DELIM)