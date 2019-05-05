这个软件到底能跑多快?
在高负载的情况下,该软件是否还能保证正确性。或者说,载荷数量与软件正确性之间的关系是怎模样的?载荷量是一个笼统的词,可以是HTTP请求的数量,也可以是API调用的次数。
在保证一定的正确性的条件下,该软件的可伸缩性是怎样的?
在软件可以正常工作的情况下,负载与系统资源(包括CPU、内存和各种I/O资源等)使用率之间的关系是怎么样的?

现在编写的载荷器可以作为软件性能评测的辅助工具,它可以向被测软件发送指定量的载荷,并记录下了被测软件处理载荷的结果。这样，你就可统计出被测软件在给定场景下的性能数值了。

# 参数和结果
1.重要的参数
最重要的一点在于一个软件在给定运行环境下最多能被多少个同时使用。在进行性能测试的时候，我们需要明确在同一时刻(或在同一时间段内)向软件发送载荷的数量。
在这一方面有两个专业术语:
- QPS(Query Per Second):每秒查询量
- TPS(Transactions Per Second):每秒事务处理量

这两者都是体现服务器软件性能的指标，其含义都是在1s之内可以正常响应请求的数量的平均值。
不同的是,前者针对的是对服务器上数据的读取操作,后者针对的是对服务器上数据的写入或修改操作。


其次,软件在承受一定量载荷的情况下对系统资源的消耗也是值得特别关注的,这与软件的可靠性洗洗相关。打个比方来说,有两个服务器软件A和B。经性能评测，A的QPS是2000，B的OPS是2200。但是由于B对系统资源的消耗较大，以及对系统资源的释放不及时。导致其在接受每秒2000个载荷并持续了200个小时之后宕机了。但是A在接受相同的负载的情况下，可以无故障的运行了200个小时。
所以应该积极的了解软件在持续接受一定量的载荷情况下，能够无差错的运行多久(也称"平均无故障时间")。通过明确的设定持续发送载荷的时间("负载持续时间")，我们就可以评估这个时间段内软件性能的具体状态。

第三个需要了解的参数是评判软件正确性的重要标准。这个参数就是载荷的处理超时时间("处理超时时间"):从向软件发出请求到接受到软件返回的相应的最大值。超过这个最大耗时就会被认为不可接受的。

- 每秒载荷量 - 负载持续时间 -处理超时时间

2.输出的结果
请求(或者说是载荷)和响应的内容、响应的状态以及请求处理耗时


```
type a struct {
	timeourNA  time.Duration //响应超时时间,单位:ns
	Ips        uint32        //每秒载荷量
	durationNS time.Duration //负载持续时间,单位:ns
}
```

负载发生器的输出是一个结果列表。但是，这里不应该使用数组或切片作为收纳结果的容器。原因是，负载发生器需要并发的发送载荷，因此也并发的输出结果。
***已知,数组和切片都不是并发安全的,GO原生的数据类型中只有通道是并发安全的,它是收纳结果的最佳容器***
因此,我将这样一个通道类型的字段也加入载荷发生器的类型声明中:

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

前面提到,代表载荷发生器实现类型的resultCh字段是chan *lib.CallResult类型的,由于结构体类型的零值不是nil,因此如果这个通道的元素类型是lib.CallResult的话,就会给后面对其中元素值的零值判断带来小麻烦。我使用它的指针类型作为通道的元素类型,既可以消除麻烦,也可以减少元素值复制带来的开销。

concurrency uint32 //载荷并发量

一旦确定了并发量,就有了控制载荷发生器使用系统资源的依据。另外,我们关心Go程序使用的goroutine的个数;过少的goroutine数量会使程序的并发程度不够，从而导致程序不能充分的利用系统资源；而过多的goroutine数量则可能会使程序的性能不升反降,因为这对于Go运行时系统及其依托的操作系统来说都会造成额外的负担。那么怎么样合理的控制程序所启用的goroutine的数量呢？

这里可以用一个goroutine票池对此作出限定。此票池中的票的数量就由concurrency字段的值决定。goroutine票池以一个缓冲通道作为载体。我们定义这个goroutine票池的接口类型名为
GoTickets。

tickets lib.GoTickets //goroutine票池

在载荷发生器运行的过程中,应该可以随时停止它,同时,根据durationNS字段的值,载荷发生器也应该能够自动停止(这可以通过传递停止信号的方式实现)。
stopSign chan struct{} //停止信号的传递通道
不过还有一种更好的选择:可以使用在Go 1.7发布时成为标准库一员的context代码包。context包中的context接口类型的一些函数，可以帮助我同时向多个goroutine通知载荷发生器需要停止。
->
ctx context.Context //上下文
cancelFunc context.CancelFunc  //取消函数

荷载发生器不止有一种状态,状态字段是数值类型的，并且足够短小，还可以用并发安全的方式操作。Go标准库里提供的原子操作方法支持的最短数值类型为uint32。载荷发生器的状态值没必要包含负值，所以这里选定uint32类型。

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

status uint32

最后,应该让载荷发生器的使用者可以根据具体需求对它进行适当的扩展和定制。为此，需要预先在其结构中添加一个字段，并以此为载荷发生器的扩展接口。

不过,在添加这个字段之前，应该搞清楚载荷发生器需要提供哪些扩展支持。首先，载荷发生器的核心功能，肯定是控制和协调载荷的生成和发送、响应的接收和验证，以最终结果的递交等一系列操作。既然由它来进行流程上的控制，那么一些具体的操作是否就可以由定制的组件来做呢?这样就可以把核心功能与扩展功能(或者说组件功能)区分开了。如此既可以保证核心功能的稳定,又可以提供较高的可扩展性。

我已经确定了一定要作为核心功能的部分,现在看看有哪些操作可以作为组件功能。显然,我不知道或者无法预测被测软件提供API的形式。况且,载荷发生器不应该对此有所约束，它们可以是任意的。因此,与调用被测软件API有关的功能应该作为组件功能，这涉及请求的发送操作和相应的接受操作。并且，既然要组件化调用被测软件API的功能，那么请求的生成操作和相应的检查操作，也肯定无法由载荷发生器本身来提供。

// 用于表示调用器的接口
type Caller interface{
    //构建请求
    BuildReq()RawReq
    //调用
    Call(req []byte,timeoutNS time.Duration)([]byte,error)
    //检查响应
    CheckResp(rawReq RawReq,rawResp RawResp)*CallResult
}

caller lib.Caller



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

#初始化
在Go中,一般会使用一类函数来创建和初始化较为复杂的结构体,这类函数的名称通常会以"New"作为前缀
eg: func NewMyGenerator()*myGenerator

根据***面向接口编程的原则***,我们不应该直接将myGenerator或*myGenerator作为上述函数的结果类型,因为着这样会是该函数及其调用方与具体实现紧密的绑定在一起。如果要修改该结构体类型或者完全换一套载荷发生器的实现，那么调用该函数的所有代码都会造成***散弹式の修改***。我们应该让这类函数返回一个***接口类型***。

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

有了这样的接口,创建和初始化载荷发生器的函数声明应该是:
func NewMyGenerator() lib.Generator

不过还需要加几个参数:

第二版声明:
func NewMyGenerator(
    caller lib.Caller,
    timeoutNS time.Duration,
    lps uint32,
    durationNS time.Duration,
    resultCh chan *lib.CallResult
)(lib.Generator,error)

并发量=单个载荷的响应超时时间/载荷的发送间隔时间

计算并发量的最大意义是:为约束并发运行的goroutine的数量提供依据。我会依据此数值确定载荷发生器的tickets字段所代表的那个goroutine票池的容量。这个容量也可以理解为goroutine票的总数量。
初始化票池:
tickets,err:=lib.NewGoTickets(gen.concurrency)

//用于表示goroutine票池的接口
type GoTickets interface{
    //获得一张票
    Take()
    //归还一张票
    Return()
    //票池是否已经被激活
    Active() bool
    //票的总数
    Total() uint32
    //剩余的票数
    Remainder()uint32
}

你也可以把goroutine票池看做是***POSIX标准中描述的多值信号量***,一个POSIX多值信号量代表了可用资源的数量，资源使用方获得或归还资源时会及时减少或者增大该多信号量的值，
以便其他使用方实时了解资源的使用情况。在该值被减至0之后，所有试图减少该值的程序都会为此而阻塞。而当该值重新增至一个正整数的时候，这些程序就都会被唤醒。lib.GoTickets接口
的Take方法和Return方法分别对应了多值信号量上的减一操作和增一操作。 -》***(这里应该想起了通道)***

首先编写lib.myGoTickets类型的基本结构:
//用于表示goroutine票池的实现
type myGoTickets struct{
    total uint32 //票的总数
    ticketCh chan struct{} //票的容器
    active bool // 票池是否已经被激活
}

票池初始化
func (gt *myGoTickets)init(total uint32) bool{
    if gt.active{
        return false
    }
    if total==0{
        return false
    }
    ch:=make(chan struct{},total)
    n:=int(total)
    for i:=0;i<n;i++{
        cn<-struct{}{}
    }
    gt.ticketCh=ch
    gt.total=total
    gt.active=true
    return true
}

在编写完这个init方法之后，lib.NewGoTickets就可以非常方便的创建并初始化一个lib.GoTickets类型的值了。
// NewGoTickets 会新建一个Goroutine票池。
func NewGoTickets(total uint32) (GoTickets, error) {
        gt := myGoTickets{}
        if !gt.init(total) {
                errMsg :=
                        fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
                return nil, errors.New(errMsg)
        }
        return &gt, nil
}

新建一个载荷发生器
// NewGenerator 会新建一个载荷发生器。
func NewGenerator(pset ParamSet) (lib.Generator, error) {

        logger.Infoln("New a load generator...")
        if err := pset.Check(); err != nil {
                return nil, err
        }
        gen := &myGenerator{
                caller:     pset.Caller,
                timeoutNS:  pset.TimeoutNS,
                lps:        pset.LPS,
                durationNS: pset.DurationNS,
                status:     lib.STATUS_ORIGINAL,
                resultCh:   pset.ResultCh,
        }
        if err := gen.init(); err != nil {
                return nil, err
        }
        return gen, nil
}

注意,NewGenerator函数的参数声明列表与之前展示的不太一样，这里把所有参数都内置到了一个名为ParamSet的结构体类型中。如此一来,在需要变动NewGenerator函数的参数时,就无需
改变它的声明了,变动只会改变ParamSet类型。另外，还为ParamSet类型添加了一个名为Check的指针方法，它会检查当前值中所有字段的有效性，一旦发现无效字段，它就会返回一个非nil的
error类型的值。这样做使得NewGeneraotr函数的调用方法可以先行检查这个待传入参数的集合的有效性。

另外,在上述代码中出现的标识符logger代表的是loadgen子包中声明的一个变量。
//日志记录器
var logger=log.DLogger()

顺便说一句,NewGenerator函数的最后调用了载荷发生器的init方法。该方法初始化了另外两个载荷发生器启动前必须的字段-concurrency和tickets。



# 四、启动和停止
如前文所述,***调用载荷发生器的Start方法就可以启动它了。然后,载荷发生器会按照给定的参数向被测软件发送一定量的载荷。在触达了指定的负载持续时间之后，载荷发生器会自动停止载荷发送曹组。***
***在从启动到停止的这个时间段内，载荷发生器还会将被测软件对各个载荷的响应，以及载荷发送操作的最终结果收集起来，并发送提供调用结果通道***。

这个流程看起来并不复杂，其中包含了很多细节。比较重要的是***有效控制载荷发送器的并发量,以及载荷发生器本身使用的goroutine的数量***

## 1.启动的准备
载荷发生器的lps字段指明了它每秒向被测软件发送载荷的数量。根据次值,可以很容易的得到发送间隔时间,表达式为1e9/lps。为了让发送间隔时间能够起到实质性的作用，需要使用缓冲通道和断续器。
//设定节流阀
var throttle <-chan time.Time
if gen.lps>0{
    interval:=time.Duration(1e9/gen.lps)
    logger.Infof("Setting throttle (%v)...",interval)
    throttle=time.Tick(interval)
}

在真正使用节流阀之前,还有另外一个准备工作要做,即让载荷发生器能够运行一段时间之后自己停下来。它的代码ctx和cancelFunc可以做这件事，context包中有函数可以为他们赋值。
//初始化上下文和取消函数
gen.ctx,gen.cancelFunc=context.WithTimeout(
    context.Background(),gen.durationNs)

载荷发生器是可以重复使用的,所以每次启动的时候都必须重置它的callCount字段，这是启动前需要做的最后一步。
//初始化调用计数
gen.CallCount=0

//设置状态为已启动
atomic.StoreUint32(&gen.Status,lib.STATUS_STARTED) //***注意这里atomic设置状态***

***注意，这里改变状态时使用了原子操作。简单来说,原子操作就是一定会一次做完的操作。在操作过程中部允许任何中断，操作所在的goroutine和内核线程也绝不会被切换下CPU。

## 2.控制流程
在进入已启动状态后,载荷发生器才真正开始生成并发送载荷。包含了***载荷发送操作***和***载荷应接收操作***的调用操作是***异步执行***的。因为之后这样载荷发生器才能做到并发运行。

***//产生载荷并向承受方发送***
func(gen *myGenerator)genLoad(throttle <-chan time.Time){
    for{
        select{
            case <-gen.ctx.Done():
                gen.prepareToStop(gen.ctx.Err())
                return
            default:
        }
        gen.asyncCall()
        if gen.lps>0{
            select{
                case <-throttle:
                case <-gen.ctx.Done():
                    gen.prepareToStop(gen.ctx.Err())
                    return
            }
        }
    }
}

由于select信号在多个满足条件的case之间做伪随机选择时的不确定性，当节流阀的到期通知和上下文的"信号"同时到达时,后者代表的case不一定活被选中，这也是为了保险起见，所以在首尾都加了select,以使载荷发生器总能及时的停止。

***prepareToStop***方法用于为停止载荷发生器做准备,代码如下:
***//用于为停止载荷发生器做准备***
func (gen *myGenerator) prepareToStop(ctxError error) {
        logger.Infof("Prepare to stop load generator (cause: %s)...", ctxError)
        atomic.CompareAndSwapUint32(
                &gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING)
        logger.Infof("Closing result channel...")
        close(gen.resultCh)
        atomic.StoreUint32(&gen.status, lib.STATUS_STOPPED) //实现了原子的CAS操作(CAS又称为"比较并交换")
}



## 3.异步的调用
genLoad方法在循环中调用了gen.asyncCall方法。后者让控制流程和调用流程分离开来,真正实现了载荷发送操作的异步性和并发性。
一个调用过程分为5个操作步骤,即***生成载荷***、***发送载荷并接收响应***、***检查载荷相应***、***生成调用结果和发送调用结果***。
前3个步骤都会由适用方初始化载荷发生器时传入的那个调用器完成。

asyncCall方法在一开始就会启用一个专用的goroutine。因为对asyncCall方法的每一次调用都会意味着有一个专用的goroutine被启用。这里的
专用goroutine总数会由goroutine票池控制，后者会有载荷发生器的tickets字段代表。因此，在该方法中，我们需要适当的时候对goroutine票
池中的票进行"获得"和"归还"操作

//异步的调用承受方接口
func(gen *myGeneraotr)asyncCall(){
    gen.tickets.Take()
    go func(){
        defer func(){
            // 省略若干代码
            gen.tickets.Return()
        }()
        // 省略若干代码
    }()
}

在启用专用goroutine之前,从goroutine票池获得了一张goroutine票。当goroutine票池中无票可拿时,asyncCall方法所在的goroutine会被阻塞
于此。只有在存在多余的goroutine票时,专用goroutine才会被启用，从而当前的调用过程才会执行。

好了,现在异步调用的框架已经有了。下面来看看专用goroutine需要执行的语句。
### 1.1 生成载荷
因为有了调用方传入的调用器,所以这里的代码相当简单:
rawReq:=gen.Caller.BuildReq()

### 1.2 发送载荷并接收相应
关于这个步骤，我详细说明一下。首先是对载荷发生器的timeoutNS字段的使用。已知,该字段起到辅助载荷发生器实时判断被测软件处理单一载荷是否超时
的作用。我之前讲过,time包中的定时器可以用来设定某一个操作或任务的超时时间。要做到实时的超时判断，最好的方式就是与通道和select语句联用，不过
这就需要再启用一个goroutine来封装调用操作。如此一来，前面提到的那个goroutine票池就收效甚微。那么，可以在不额外启用goroutine的情况下实现
实时的超时判断吗?答案是:可以，但是这需要一些技巧。

具体来说,先要声明代表调用状态的变量,并保证仅其上实施原子操作。
//调用状态: 0-未调用或调用中; 1-调用完成; 2-调用超时
var callStatus uint32

然后,使用time包的AfterFunc函数设定超时以及后续处理:
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

在作为第二个参数的匿名函数中,先用atomic.CompareAndSwapUint32函数检查并设置callStatus的值，该函数会返回一个bool类型的值，用以表示比较并
交换是否成功。如果未成功，就说明载荷响应接收操作已先完成，忽略超时处理。

相对的,实现调用操作的代码也应该这样使用callStatus变量:
rawResp:=gen.callOne(&rawReq)
if !atomic.CompareAndSwapUint32(&callStatus,0,1){
    return
}
timer.Stop()

无论何时,一旦代表调用操作的callOne方法已经返回，就要检查并设置callStatus变量的值。如CAS操作不成功，就说明调用操作已经超时，前面传入
time.AfterFunc函数的那个匿名函数就已经先执行了，需要忽略后续的相应处理。当然，在相应处理之前要先停掉前面的定时器。

***相应处理***无非就是对原始响应进行再包装,然后再把包装后的响应发给调用结果通道。包装方式和响应发送方式都与那个匿名函数中的类似。不过，这里
要先检查是否存在调用错误:
var result *lib.CallResult
if rawResp.Err != nil {
        result = &lib.CallResult{
                ID:     rawResp.ID,
                Req:    rawReq,
                Code:   lib.RET_CODE_ERROR_CALL,
                Msg:    rawResp.Err.Error(),
                Elapse: rawResp.Elapse}
} else {
        result = gen.caller.CheckResp(rawReq, *rawResp)
        result.Elapse = rawResp.Elapse
}
gen.sendResult(result)

载荷发生器的sendResult方法会向调用结果通道发送一个调用结果值。为了确保调用结果发送的正确性,sendResult方法必须先检查载荷发生器的状态。
如果它的状态不是已启动，就不能执行发送操作了。不过，虽然终止了发送，但仍需要记录调用结果。另外，若调用通道已满，也不能执行发送操作。由于
该通道是载荷发生器的使用方传入的，因此无法保证没有这个情况发生。因此，这里需要把发送操作作为一条select语句中的一个case,并添加default
分支以确保不会发生阻塞。

***//用于发送调用结果***
// sendResult 用于发送调用结果。
func (gen *myGenerator) sendResult(result *lib.CallResult) bool {
        if atomic.LoadUint32(&gen.status) != lib.STATUS_STARTED {
                gen.printIgnoredResult(result, "stopped load generator")
                return false
        }
        select {
        case gen.resultCh <- result:
                return true
        default:
                gen.printIgnoredResult(result, "full result channel")
                return false
        }
}

在不能执行发送操作时,会通过调用printIgnoreResult方法记录未发送的结果。

// RetCode 表示结果代码的类型。
type RetCode int

// 保留 1 ~ 1000 给载荷承受方使用。
const (
        RET_CODE_SUCCESS              RetCode = 0    // 成功。
        RET_CODE_WARNING_CALL_TIMEOUT         = 1001 // 调用超时警告。
        RET_CODE_ERROR_CALL                   = 2001 // 调用错误。
        RET_CODE_ERROR_RESPONSE               = 2002 // 响应内容错误。
        RET_CODE_ERROR_CALEE                  = 2003 // 被调用方（被测软件）的内部错误。
        RET_CODE_FATAL_CALL                   = 3001 // 调用过程中发生了致命错误！
)

为了把调用器可能引发的运行时恐慌转变成错误,需要确保在asyncCall方法中的go函数的开始处有一条defer语句
defer func(){
    if p:=recover();p!=nil{
        err,ok:=interface{}(p).(error)
        var errMsg string
        if ok{
            errMsg=fmt.Spintf
        }else{
            errMsg=
        }
        logger.Errorln(errMsg)
        result:=&lib.CallResult{
            ID: -1,
            Code: lib.RET_CODE_FATAL_CALL,
            Msg: errMsg
        }
        gen.sendResult(result)
    }
    gen.tickets.Return() //***这里很重要***
}


至此,asyncCall方法已经拼接完整。其中调用的callOne方法的具体实现包括:
- 原子地递增载荷发生器的callCount字段值
- 检查参数
- 执行调用
- 记录调用时长
- 组装并返回原始响应值


***// callOne 会向载荷承受方发起一次调用。***
func (gen *myGenerator) callOne(rawReq *lib.RawReq) *lib.RawResp {
        atomic.AddInt64(&gen.callCount, 1)
        if rawReq == nil {
                return &lib.RawResp{ID: -1, Err: errors.New("Invalid raw request.")}
        }
        start := time.Now().UnixNano()
        resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNS)
        end := time.Now().UnixNano()
        elapsedTime := time.Duration(end - start)
        var rawResp lib.RawResp
        if err != nil {
                errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
                rawResp = lib.RawResp{
                        ID:     rawReq.ID,
                        Err:    errors.New(errMsg),
                        Elapse: elapsedTime}
        } else {
                rawResp = lib.RawResp{
                        ID:     rawReq.ID,
                        Resp:   resp,
                        Elapse: elapsedTime}
        }
        return &rawResp
}

### 1.4 启动和停止
再次回到启动载荷发生器的话题上来。在一切准备工作完成之后,载荷发生器的Start方法会调用前面展示genLoad方法以执行控制流程。genLoad方法一旦
从gen.ctx字段那里得到需要停止的"信号"，就会立即作出响应。此时,prepareToStop方法就会被调用，该方法会在关闭调用通道之后把字啊和发生器的
状态设置为已停止。在载荷发生器的Start方法的最后,有这样一段代码:
go func(){
    //生成载荷
    logger.Infoln("Generatine loads...")
    gen.genLoad(throttle)
    logger.Infof("Stopped.(call count: %d)",gen.callCount)
}()
return true


上述goroutine主要为了满足载荷发生器的自动停止,下面再来看载荷发生器的手动停止。可以通过调用载荷发生器的Stop方法来手动地停止它,Stop方法的实现如下:

func(gen *myGenerator)Stop()bool{
    if !atomic.CompareAndSwapUint32(&gen.status,lib.STATUS_STARTED,lib.STATUS_STOPPING){
        return false
    }
    gen.cancelFunc()
    for{
        if atomic.LoadUint32(&gen.status)==lib.STATUS_STOPPED{
            break
        }
        time.Sleep(time.Microsecond)
    }
    return true
}

Stop方法首先需要检查载荷负载器的状态，若状态不对就直接返回false。执行cancelFunc字段所代表的方法可以让ctx字段发出停止"信号"。而后,Stop方法
需要不断检查状态的变更。如果状态变为已停止,就说明prepareToStop方法已执行完毕，这时就可以返回了。

***特别提示***,channel的玩儿法有很多,需要根据具体情况选择和组合，不过,channel并不是银弹,它最大的使用场景是作为goroutine之间的数据传输通道,不要为了使用而使用。有时候，同步方法更加直截了当，并且更快，就像我在载荷发生器实现中多次用到的***原子操作***那样。


# 调用器和功能调试
这里是一个简单的调用器的demo
它提供了一个基于网络的API,该API的功能就是根据请求中的参数进行简单且有限的算术运算(针对整数的加减乘除运算),并将结果作为响应返回请求方。

//用于表示服务器请求的结构
type ServerReq struct{
    ID int64
    Operands []int
    Operator string
}

//用于表示服务器响应的结构
type ServerResp struct{
    ID int64
    Formula string
    Result int 
    Err error
}

## 1.调用器实现的基本结构
TCPComm 
//用于表示TCP通信器的结构
type TCPComm struct{
    addr string
}

字段addr "127.0.0.1:8080"

## 2.BuildReq方法
// BuildReq 会构建一个请求。
func (comm *TCPComm) BuildReq() loadgenlib.RawReq {
        id := time.Now().UnixNano()
        sreq := ServerReq{
                ID: id,
                Operands: []int{
                        int(rand.Int31n(1000) + 1),
                        int(rand.Int31n(1000) + 1)},
                Operator: func() string {
                        return operators[rand.Int31n(100)%4]
                }(),
        }
        bytes, err := json.Marshal(sreq)
        if err != nil {
                panic(err)
        }
        rawReq := loadgenlib.RawReq{ID: id, Req: bytes}
        return rawReq
}


## 3.Call方法
下面就来看看*TCPComm的Call方法的具体实现:
//发起一次通信
func(comm *TCPComm)Call(req []byte,timeoutNS time.Duration)([]byte,error){
    conn,err:=net.DialTimeout("tcp",comm.addr,timeoutNS)
    if err!=nil{
        return nil,err
    }
    _,err=write(conn,req,DELIM)
    if err!=nil{
        return nil,err
    }
    return read(conn,DELIM)
}

需要注意的是,已知,基于TCP协议的通信是使用字节流来传递上层给予的消息的。它会根据具体情况为消息分段,但却无法感知消息的边界。因此,需要显式的为
请求数据添加结束符,而传给write方法和read方法的参数DELIM就代表了这个结束符,这两个方法会用它来切分出单个的请求或响应。

## 4.CheckResp方法
// CheckResp 会检查响应内容。
func (comm *TCPComm) CheckResp(
        rawReq loadgenlib.RawReq, rawResp loadgenlib.RawResp) *loadgenlib.CallResult {
        var commResult loadgenlib.CallResult
        commResult.ID = rawResp.ID
        commResult.Req = rawReq
        commResult.Resp = rawResp
        var sreq ServerReq
        err := json.Unmarshal(rawReq.Req, &sreq)
        if err != nil {
                commResult.Code = loadgenlib.RET_CODE_FATAL_CALL
                commResult.Msg =
                        fmt.Sprintf("Incorrectly formatted Req: %s!\n", string(rawReq.Req))
                return &commResult
        }
        var sresp ServerResp
        err = json.Unmarshal(rawResp.Resp, &sresp)
        if err != nil {
                commResult.Code = loadgenlib.RET_CODE_ERROR_RESPONSE
                commResult.Msg =
                        fmt.Sprintf("Incorrectly formatted Resp: %s!\n", string(rawResp.Resp))
                return &commResult
        }
        if sresp.ID != sreq.ID {
                commResult.Code = loadgenlib.RET_CODE_ERROR_RESPONSE
                commResult.Msg =
                        fmt.Sprintf("Inconsistent raw id! (%d != %d)\n", rawReq.ID, rawResp.ID)
                return &commResult
        }
        if sresp.Err != nil {
                commResult.Code = loadgenlib.RET_CODE_ERROR_CALEE
                commResult.Msg =
                        fmt.Sprintf("Abnormal server: %s!\n", sresp.Err)
                return &commResult
        }
        if sresp.Result != op(sreq.Operands, sreq.Operator) {
                commResult.Code = loadgenlib.RET_CODE_ERROR_RESPONSE
                commResult.Msg =
                        fmt.Sprintf(
                                "Incorrect result: %s!\n",
                                genFormula(sreq.Operands, sreq.Operator, sresp.Result, false))
                return &commResult
        }
        commResult.Code = loadgenlib.RET_CODE_SUCCESS
        commResult.Msg = fmt.Sprintf("Success. (%s)", sresp.Formula)
        return &commResult
}


# 测试代码:/
package loadgen

import (
        "testing"
        "time"

        loadgenlib "gopcp.v2/chapter4/loadgen/lib"
        helper "gopcp.v2/chapter4/loadgen/testhelper"
)

// printDetail 代表是否打印详细结果。
var printDetail = false

func TestStart(t *testing.T) {

        // 初始化服务器。
        server := helper.NewTCPServer()
        defer server.Close()
        serverAddr := "127.0.0.1:8080"
        t.Logf("Startup TCP server(%s)...\n", serverAddr)
        err := server.Listen(serverAddr)
        if err != nil {
                t.Fatalf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
                t.FailNow()
        }

        // 初始化载荷发生器。
        pset := ParamSet{
                Caller:     helper.NewTCPComm(serverAddr),
                TimeoutNS:  50 * time.Millisecond,
                LPS:        uint32(1000),
                DurationNS: 10 * time.Second,
                ResultCh:   make(chan *loadgenlib.CallResult, 50),
        }
        t.Logf("Initialize load generator (timeoutNS=%v, lps=%d, durationNS=%v)...",
                pset.TimeoutNS, pset.LPS, pset.DurationNS)
        gen, err := NewGenerator(pset)
        if err != nil {
                t.Fatalf("Load generator initialization failing: %s\n",
                        err)
                t.FailNow()
        }

        // 开始！
        t.Log("Start load generator...")
        gen.Start()

        // 显示结果。
        countMap := make(map[loadgenlib.RetCode]int)
        for r := range pset.ResultCh {
                countMap[r.Code] = countMap[r.Code] + 1
                if printDetail {
                        t.Logf("Result: ID=%d, Code=%d, Msg=%s, Elapse=%v.\n",
                                r.ID, r.Code, r.Msg, r.Elapse)
                }
        }

        var total int
        t.Log("RetCode Count:")
        for k, v := range countMap {
                codePlain := loadgenlib.GetRetCodePlain(k)
                t.Logf("  Code plain: %s (%d), Count: %d.\n",
                        codePlain, k, v)
                total += v
        }

        t.Logf("Total: %d.\n", total)
        successCount := countMap[loadgenlib.RET_CODE_SUCCESS]
        tps := float64(successCount) / float64(pset.DurationNS/1e9)
        t.Logf("Loads per second: %d; Treatments per second: %f.\n", pset.LPS, tps)
}

func TestStop(t *testing.T) {

        // 初始化服务器。
        server := helper.NewTCPServer()
        defer server.Close()
        serverAddr := "127.0.0.1:8081"
        t.Logf("Startup TCP server(%s)...\n", serverAddr)
        err := server.Listen(serverAddr)
        if err != nil {
                t.Fatalf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
                t.FailNow()
        }

        // 初始化载荷发生器。
        pset := ParamSet{
                Caller:     helper.NewTCPComm(serverAddr),
                TimeoutNS:  50 * time.Millisecond,
                LPS:        uint32(1000),
                DurationNS: 10 * time.Second,
                ResultCh:   make(chan *loadgenlib.CallResult, 50),
        }
        t.Logf("Initialize load generator (timeoutNS=%v, lps=%d, durationNS=%v)...",
                pset.TimeoutNS, pset.LPS, pset.DurationNS)
        gen, err := NewGenerator(pset)
        if err != nil {
                t.Fatalf("Load generator initialization failing: %s.\n",
                        err)
                t.FailNow()
        }

        // 开始！
        t.Log("Start load generator...")
        gen.Start()
        timeoutNS := 2 * time.Second
        time.AfterFunc(timeoutNS, func() {
                gen.Stop()
        })

        // 显示调用结果。
        countMap := make(map[loadgenlib.RetCode]int)
        count := 0
        for r := range pset.ResultCh {
                countMap[r.Code] = countMap[r.Code] + 1
                if printDetail {
                        t.Logf("Result: ID=%d, Code=%d, Msg=%s, Elapse=%v.\n",
                                r.ID, r.Code, r.Msg, r.Elapse)
                }
                count++
        }

        var total int
        t.Log("RetCode Count:")
        for k, v := range countMap {
                codePlain := loadgenlib.GetRetCodePlain(k)
                t.Logf("  Code plain: %s (%d), Count: %d.\n",
                        codePlain, k, v)
                total += v
        }

        t.Logf("Total: %d.\n", total)
        successCount := countMap[loadgenlib.RET_CODE_SUCCESS]
        tps := float64(successCount) / float64(timeoutNS/1e9)
        t.Logf("Loads per second: %d; Treatments per second: %f.\n", pset.LPS, tps)
}
