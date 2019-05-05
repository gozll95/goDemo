# 网络爬虫

##一、网络爬虫和框架
网络爬虫应该根据使用者的意愿自动下载、分析、筛选、统计以及存储指定的网络内容。注意,这里的关键词是***"自动"***和***"根据意愿"***,"自动"的含义是:网络爬虫在启动后自己完成整个爬取过程而无需人工干预,并且还能够在过程结束之后自动停止。而"根据意愿"则是说,网络爬虫最大限度的允许使用者对其爬取过程进行定制。

需要特别处理的细节:
- 有效网络地址的发现和提取
- 有效网络地址的边界定义和检查
- 重复的网络地址的过滤

##二、功能需求和分析
概括来说,网络爬虫框架会反复执行如下步骤直至触碰到停止条件:
- (1)"下载器"下载与给定网络地址相对应的内容。其中，在下载"请求"的组装方面,网络爬虫框架为使用者尽量预留出定制接口。使用者可以使用这些接口自定义"请求"的组装方法。
- (2)"分析器"分析下载到的内容,并从中筛选出可用的部分(以下称为"条目")和需要访问的新网络地址。其中,在用于分析和筛选内容的规则和策略方面,应该由网络爬虫框架提供灵活的定制接口。换句话说,由于只有使用者自己才知道他们真正想要的是什么,所以应该允许他们对这些规则和策略进行深入的定制。网络爬虫框架仅需要规定好定制的方式即可。
- (3)"分析器"把筛选出的"条目"发送给"条目处理管道"。同时,它会把发现的新网络地址和其他一些信息组装成新的下载"请求",然后把这些请求发送给"下载器"。在此步骤中,我们会过滤掉一些不符合要求的网络地址,比如忽略超出有效边界的网络地址。

这里,我再次强调一下网络爬虫框架与网络爬虫实现的区别。作为一个框架,该程序在每个处理模块中给予使用者尽量多的定制方法,而不去涉及各个处理步骤的实现细节。另外,框架更多的考虑使用者自定义的处理步骤在执行期间可能发生的各种情况和问题,并注意对这些问题的处理方式。这样才能在易于扩展的同时保证框架的稳定性。

## 三、总体设计
网络爬虫框架的处理模块有3个:下载器、分析器和条目处理管道。再加上调度和协调这些处理模块运行的控制模块,我们就可以明晰该框架的模块划分了。我把这里提到的控制模块称为"调度器"。
- 下载器:接受请求类型的数据,并依据该请求获得HTTP请求;将HTTP请求发送至与指定的网络地址对应的远程服务器;在HTTP请求发送完毕之后,立即等待相应的HTTP相应的到来;在收到HTTP相应之后,将其封装成响应并作为输出返回给下载器的调用方。其中,HTTP客户端程序可以由网络爬虫框架的使用方自行定义。另外,若在该子流程执行期间发生了错误,应该立即以适当的方法告知适用方。对于其他模块来讲,也是这样。
- 分析器:接收相应类型的数据,并依据该响应获得HTTP相应;对该HTTP相应的内容进行检查,并根据给定的规则进行分析、筛选以及生成新的请求和条目;将生成的请求或条目作为输出返回给分析器的调用方。在分析器的职责中,我可以想到的能留给网络爬虫框架的使用方自定义的部分并不少。例如:对HTTP相应的前期检查、对内容的筛选,以及生成请求和条目的方式,等等。不过,我在后面回对这些可以自定义的部分进行一些取舍。
- 条目处理通道:接受条目类型的数据,并对其执行若干步骤的处理;条目处理管道中可以产出最终的数据;这个最终的数据可以在其中的某个处理步骤中被持久化(不论是本地存储还是发送给远程的存储服务器)以备后用。我们可以把这些处理步骤的具体实现留给网络爬虫框架的使用方自定义。这样,网络爬虫框架就可以真正的与条目处理的细节脱离开来。网络爬虫框架丝毫不关心这些条目怎么样被处理和持久化,它仅仅负责控制整体的处理流程。我把负责单个处理步骤的程序称为条目处理器。条目处理管道接受条目类型的数据,并把处理完成的条目返回给条目处理管道。条目处理管道会紧接着把该条目传递给下一个条目处理器,直至给定的条目处理列表中的每个条目处理器都处理过该条目为止。
- 调度器:调度器在启动时仅接收首次请求,并且不会产生任何输出。调度器的主要职责是调度各个处理模块的运行。其中包括维护各个处理模块的实例、在不同的处理模块之间传递数据(包括请求、响应和条目),以及监听所有这些被调度者的状态,等等。有了调度器的维护,各个处理模块得以保持其职责的简洁和专一。由于调度器是网络爬虫框架最重要的一个模块,所以还需要再编写一些工具来支撑它的功能。

这里需要说明***条目处理通道***,它是以***流式***为基础的,其设计灵感来自于我之前讲过的Linux系统中的管道。我们可以不断地向该管道发送条目,而该管道则会让其中的若干个条目处理器依次处理每一个条目。我们可以轻易的使用一些同步方法来保证条目处理管道的并发性,因此即使调度器只持有该管道的一个实例,要不会有任何问题。



## 四、详细设计

### 4.1 基本数据结构

在分析网络爬虫框架的需求时,提到过这样几类数据--请求、相应、条目,下面我们逐个讲解它们的声明和设计理念。

#### 1)请求
- 请求:用来承载向某个网络地址发起的HTTP请求,它由调度器或分析器生成并传递给下载器,下载器会根据它从远程服务器下载相应的内容。因此,它有一个net/http.Request类型的字段。***不过,为了减少不必要的零值的生成(http.Request是一个结构体类型,它的零值不是nil)和实例复制,我们把*http.Request作为该字段的类型。**

v1:
//数据请求的类型
type Request struct{
    //HTTP请求
    httpReq *http.Request
}

量化内容爬取程度的一个比较常用的方法,是计算每个下载的网络内容的深度。
v2:

//数据请求的类型
type Request struct{
    //HTTP请求
    httpReq *http.Request
    //请求的深度
    depth uint32
}


//用于创建一个新的请求实例
func NewRequest(httpReq *http.Request,depth uint32)*Request{
    return &Request{httpReq:httpReq,depth:depth}
}

// HTTPReq 用于获取HTTP请求。
func (req *Request) HTTPReq() *http.Request {
	return req.httpReq
}

// Depth 用于获取请求的深度。
func (req *Request) Depth() uint32 {
	return req.depth
}

我希望这个类型的值是不可变的。也就是说,在该类型的一个值创建和初始化之后,当前代码包之外的任何代码都不能更改它的任何字段值。
基于这样的需求,一般都会通过三个步骤来实现:
- 把该类型的所有字段的访问权限都设计为包级私有。也就是说,要保证这些字段的首字母均为小写。
- 编写一个创建和初始化该类型值的函数。Newxxx
- 编写必要的用来获取字段值的方法。

注意,NewRequest函数的结果类型是*Request,而不是Request,这样做的主要原因为xxxxx,更深层次的原因是,值在作为参数传递给函数或者作为结果由函数返回时会被复制一次。***指针值往往更能减少复制的开销***。

关于深度:一个请求的深度值=对它的父请求的深度值+1


#### 2)响应
//数据响应的类型
type Response struct{
    //HTTP响应
    httpResp *http.Response
    //响应的深度
    depth uint32
}

// NewResponse 用于创建一个新的响应实例。
func NewResponse(httpResp *http.Response, depth uint32) *Response {
	return &Response{httpResp: httpResp, depth: depth}
}

// HTTPResp 用于获取HTTP响应。
func (resp *Response) HTTPResp() *http.Response {
	return resp.httpResp
}

// Depth 用于获取响应深度。
func (resp *Response) Depth() uint32 {
	return resp.depth
}

#### 3)条目
//条目的类型
type Item map[string]interface{}

####4)其他类型
好了,我们需要的3个基本数据类型都在这里了。为了能够用一个类型从整体上标识这3个基本数据类型,我们又声明了Data接口类型:

//数据的接口类型
type Data interface{
    //用于判断数据是否有效
    Valid() bool
}

这个接口类型只有名为Valid的方法,可以通过调用该方法来判断数据的有效性。显然,Data接口类型的作用更多的是作为数据类型的一个标签,而不是定义某种类型的行为。为了让表示请求、响应或条目的类型都实现Data接口,又在当前的源码文件添加了这样几个方法:

// Valid 用于判断请求是否有效。
func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

// Valid 用于判断响应是否有效。
func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

// Valid 用于判断条目是否有效。
func (item Item) Valid() bool {
	return item != nil
}

我们还需要一个额外的类型,这个类型是作为error接口类型的实现类型而存在的。它的主要作用是封装爬取过程中的错误,并以统一的方式生成字符串形式的描述。我们知道,只要某个类型的方法集合中包含了下面这个方法,就等于实现了error接口类型:

func Error() string

首先,声明一个名为CrawlerError的接口类型:

// CrawlerError 代表爬虫错误的接口类型。
type CrawlerError interface {
	// Type 用于获得错误的类型。
	Type() ErrorType
	// Error 用于获得错误提示信息。
	Error() string
}

由于CrawlerError类型的声明中也包含了Error方法,所以只要某个类型实现了它,就等于实现了error接口类型。先编写这样一个接口类型而不是直接编写出error接口类型的实现类型的原因有两个:
- 我们在编程过程中应该遵循面向接口编程的原则
- 为了扩展error接口类型。网络爬虫框架拥有多个处理模块,错误类型值可以表明该错误是哪一个处理模块产生的,这也是Type方法起到的作用。

// myCrawlerError 代表爬虫错误的实现类型。
type myCrawlerError struct {
	// errType 代表错误的类型。
	errType ErrorType
	// errMsg 代表错误的提示信息。
	errMsg string
	// fullErrMsg 代表完整的错误提示信息。
	fullErrMsg string
}

// NewCrawlerError 用于创建一个新的爬虫错误值。
func NewCrawlerError(errType ErrorType, errMsg string) CrawlerError {
	return &myCrawlerError{
		errType: errType,
		errMsg:  strings.TrimSpace(errMsg),
	}
}

func (ce *myCrawlerError) Type() ErrorType {
	return ce.errType
}

func (ce *myCrawlerError) Error() string {
	if ce.fullErrMsg == "" {
		ce.genFullErrMsg()
	}
	return ce.fullErrMsg
}

你可以能已经发现,Error方法中用到了myCrawlerError类型的fullErrMsg字段。并且,它还调用了一个名为genFullErrMsg的方法,该方法的实现类型如下:

// genFullErrMsg 用于生成错误提示信息，并给相应的字段赋值。
func (ce *myCrawlerError) genFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("crawler error: ")
	if ce.errType != "" {
		buffer.WriteString(string(ce.errType))
		buffer.WriteString(": ")
	}
	buffer.WriteString(ce.errMsg)
	ce.fullErrMsg = fmt.Sprintf("%s", buffer.String())
	return
}

这里看到,没有直接使用errMsg字段的值,而是以它为基础生成了一条完整的错误提示信息。在这条信息中,明确显示出它是一个网络爬虫的错误,也给出了错误的类型和详情。注意,这条错误提示信息缓存在fullErrMsg字段中。回顾该类型的Error方法的实现,只有当fullErrMsg字段的值为""时,才会调用genFullErrMsg方法,否则会直接把fullErrMsg字段的值作为Error方法的结果值返回。这也是为了避免频繁的拼接字符串给程序性能带来的负面影响。


### 4.2 接口的设计

#### 1)下载器
从***下载器***充当的角色来讲,它的功能只有两个:发送请求和接收响应。因此,我可以设计出这样一个方法声明:

//用于根据请求获取内容并返回响应
Download(req *Request)(*Response,error)

Download的签名完全体现了下载器应有的功能。但是作为处理模块,下载器还应该拥有一些方法以供统计、描述只用。不过正因为这些方法是所有处理模块都具备的,所以还要编写一个更加抽象的接口类型。


// Module 代表组件的基础接口类型。
// 该接口的实现类型必须是并发安全的！
type Module interface {
	// ID 用于获取当前组件的ID。
	ID() MID
	// Addr 用于获取当前组件的网络地址的字符串形式。
	Addr() string
	// Score 用于获取当前组件的评分。
	Score() uint64
	// 用于设置当前组件的评分。
	SetScore(score uint64)
	// ScoreCalculator 用于获取评分计算器。
	ScoreCalculator() CalculateScore
	// CallCount 用于获取当前组件被调用的计数。
	CalledCount() uint64
	// AcceptedCount 用于获取被当前组件接受的调用的计数。
	// 组件一般会由于超负荷或参数有误而拒绝调用。
	AcceptedCount() uint64
	// CompletedCount 用于获取当前组件已成功完成的调用的计数。
	CompletedCount() uint64
	// HandlingNumber 用于获取当前组件正在处理的调用的数量。
	HandlingNumber() uint64
	//Counts 用于一次性获取所有计数。
	Counts() Counts
	// Summary 用于获取组件摘要。
	Summary() SummaryStruct
}


MID是string的别名类型,它的值一般由3部分组成:标识组件类型的字母、代表生成顺序的序列号和用于定位组件的网络地址。网络地址是可选的,因为组件实例可以和网络爬虫的主程序处于同一个进程中

// midTemplate 代表组件ID的模板。
var midTemplate = "%s%d|%s"

说到标识组件类型的字母,就要首先介绍一下组件的类型。

// Type 代表组件的类型。
type Type string

// 当前认可的组件类型的常量。
const (
	// TYPE_DOWNLOADER 代表下载器。
	TYPE_DOWNLOADER Type = "downloader"
	// TYPE_ANALYZER 代表分析器。
	TYPE_ANALYZER Type = "analyzer"
	// TYPE_PIPELINE 代表条目处理管道。
	TYPE_PIPELINE Type = "pipeline"
)


// legalTypeLetterMap 代表合法的组件类型-字母的映射。
var legalTypeLetterMap = map[Type]string{
	TYPE_DOWNLOADER: "D",
	TYPE_ANALYZER:   "A",
	TYPE_PIPELINE:   "P",
}

组件ID中的序列号可以由网络爬虫框架的使用方提供。这就需要我们在框架内提供一个工具,以便于统一序列号的生成和获取。序列号原则上是不能重复的,也是顺序给出的。但是如果序列号超出了给定范围,就可以循环使用。据此,我编写了一个序列号生成器的接口类型:

// SNGenertor 代表序列号生成器的接口类型。
type SNGenertor interface {
	// Start 用于获取预设的最小序列号。
	Start() uint64
	// Max 用于获取预设的最大序列号。
	Max() uint64
	// Next 用于获取下一个序列号。
	Next() uint64
	// CycleCount 用于获取循环计数。
	CycleCount() uint64
	// Get 用于获得一个序列号并准备下一个序列号。
	Get() uint64
}

其中,最小序列号和最大序列号都可以由使用方在初始化序列号生成器的时给定。循环计数器代表了生成器生成的序列号在前两者指定的范围内循环的次数。


Module接口中的第3个至第5个方法是关于组件评分的,这又涉及组件注册方面的设计。按照我的设想,在网络爬虫程序真正启动之前,应该先向***组件注册器***注册足够的组件实例。只有如此,程序才能正常运转。***组件注册器***可以***注册***、***注销***以及***获取某类组件的实例***,并且还可以***清空所有组件实例***。

// Registrar 代表组件注册器的接口。
type Registrar interface {
	// Register 用于注册组件实例。
	Register(module Module) (bool, error)
	// Unregister 用于注销组件实例。
	Unregister(mid MID) (bool, error)
	// Get 用于获取一个指定类型的组件的实例。
	// 本函数应该基于负载均衡策略返回实例。
	Get(moduleType Type) (Module, error)
	// GetAllByType 用于获取指定类型的所有组件实例。
	GetAllByType(moduleType Type) (map[MID]Module, error)
	// GetAll 用于获取所有组件实例。
	GetAll() map[MID]Module
	// Clear 会清除所有的组件注册记录。
	Clear()
}

这个接口的Get方法用于获取一个特定类型的组件实例,它实现某种负载均衡策略使得同一类型的多个组件有相对平均的机会作为结果返回。这里所说的负载均衡策略就是基于组件评分。组件评分可以通过Module接口定义的Score方法获得。相对的,SetScore方法用于设置评分。这个评分的计算方法抽象为名为CalculateScore的函数类型

其声明如下:
//用于计算组件评分的函数类型
type CalculateScore func(counts Counts)uint64

Module接口之所以没有包含设置评分计算器的方法,是因为评分计算器在初始化组件实例时给定,并且之后不能变更。

Module接口的最后一个方法Summary,用于获取组件实例的摘要信息。注意,这个摘要信息并不是字符串形式的,而是SummaryStruct类型的。这种结构化的摘要信息对于控制模块和监控工具都更加友好,同时也有助于组装和嵌入。

// SummaryStruct 代表组件摘要结构的类型。
type SummaryStruct struct {
	ID        MID         `json:"id"`
	Called    uint64      `json:"called"`
	Accepted  uint64      `json:"accepted"`
	Completed uint64      `json:"completed"`
	Handling  uint64      `json:"handling"`
	Extra     interface{} `json:"extra,omitempty"`
}


注意一下,Extra字段,该字段的作用是为额外的组件信息的纳入提供支持。

讲完了Module接口的声明以及相关的各种类型定义和设计理念,让我们再回过头去接着设计下载器的接口。有了上述的一些列铺垫,组件实例的基本结构和方法以及对它们的管理规则都已经比较明确了。
下载器的接口声明反而变的简单了,如下:

// Downloader 代表下载器的接口类型。
// 该接口的实现类型必须是并发安全的！
type Downloader interface {
	Module
	// Download 会根据请求获取内容并返回响应。
	Download(req *Request) (*Response, error)
}


### 3)条目处理通道

条目处理管道的功能就是为条目的处理提供环境,并控制整体的处理流程,具体的处理步骤由网络爬虫框架的提供者提供。实现单一处理步骤的程序称为"条目处理器",它的类型同样由单一的函数类型代表,所以也可以称之为"条目处理函数"。这又会是一层***双层定制接口***。

// Pipeline 代表条目处理管道的接口类型。
// 该接口的实现类型必须是并发安全的！
type Pipeline interface {
	Module
	// ItemProcessors 会返回当前条目处理管道使用的条目处理函数的列表。
	ItemProcessors() []ProcessItem
	// Send 会向条目处理管道发送条目。
	// 条目需要依次经过若干条目处理函数的处理。
	Send(item Item) []error
	// FailFast方法会返回一个布尔值。该值表示当前条目处理管道是否是快速失败的。
	// 这里的快速失败是指：只要在处理某个条目时在某一个步骤上出错，
	// 那么条目处理管道就会忽略掉后续的所有处理步骤并报告错误。
	FailFast() bool
	// 设置是否快速失败。
	SetFailFast(failFast bool)
}

//用于处理条目的函数的类型
type ProcessItem func(item Item)(result Item,err error)

#### pay attation to
最后,一定要注意,与下载器和分析器一样,条目处理管道的实现也一定要是并发安全的。也就是说,它们的任何方法在同时调用时都不能产生竞态条件。这主要是因为调度器会在任何需要的时候从组件注册器中获取一个组件实例并使用。同一个组件实例可能会用来并发处理多个数据。组件实例不能成为调度器执行并发调度的阻碍。此外,与之有关的各种计数和摘要信息的读写操作要求组件本身具有并发安全性。


### 4) 调度器
调度器属于控制模块而非处理模块,它需要对各个模块的运作进行调度和控制。可以说,调度器是网络爬虫框架的心脏。因此,我需要由它来***启动***和***停止***爬取流程,另外,出于监控整个爬取流程的目的,还应该在这里提供***获取实时状态***和***摘要信息***的方法。

// Scheduler 代表调度器的接口类型。
type Scheduler interface {
	// Init 用于初始化调度器。
	// 参数requestArgs代表请求相关的参数。
	// 参数dataArgs代表数据相关的参数。
	// 参数moduleArgs代表组件相关的参数。
	Init(requestArgs RequestArgs,
		dataArgs DataArgs,
		moduleArgs ModuleArgs) (err error)
	// Start 用于启动调度器并执行爬取流程。
	// 参数firstHTTPReq即代表首次请求。调度器会以此为起始点开始执行爬取流程。
	Start(firstHTTPReq *http.Request) (err error)
	// Stop 用于停止调度器的运行。
	// 所有处理模块执行的流程都会被中止。
	Stop() (err error)
	// Status 用于获取调度器的状态。
	Status() Status
	// ErrorChan 用于获得错误通道。
	// 调度器以及各个处理模块运行过程中出现的所有错误都会被发送到该通道。
	// 若结果值为nil，则说明错误通道不可用或调度器已被停止。
	ErrorChan() <-chan error
	// Idle 用于判断所有处理模块是否都处于空闲状态。
	Idle() bool
	// Summary 用于获取摘要实例。
	Summary() SchedSummary
}

Scheduler接口的Init方法用于调度器的初始化。初始化调度器需要一些参数,这些参数分为3类:
- 请求相关参数 ***RequestArgs***
- 数据相关参数 ***DataArgs***
- 组件相关参数 ***ModuleArgs***

// RequestArgs 代表请求相关的参数容器的类型。
type RequestArgs struct {
	// AcceptedDomains 代表可以接受的URL的主域名的列表。
	// URL主域名不在列表中的请求都会被忽略，
	AcceptedDomains []string `json:"accepted_primary_domains"` //明确广度
	// maxDepth 代表了需要被爬取的最大深度。
	// 实际深度大于此值的请求都会被忽略。
	MaxDepth uint32 `json:"max_depth"`      //明确深度
}

DataArgs类型中包括的是与***数据缓冲池***相关的字段,这些字段的值用于初始化对应的数据缓冲池。调度器使用这些数据缓冲池传递数据。

具体来说,调度器使用的数据缓冲池有4个:
- 请求缓冲池:传输请求类型
- 响应缓冲池:传输响应类型
- 条目缓冲池:传输条目类型
- 错误缓冲池:传输错误类型

每个缓冲池需要2个参数:缓冲池中单一缓冲器的容量+缓冲池包含的缓冲器的最大数量。

这样算来,DataArgs类型中字段的总数就是8

// DataArgs 代表数据相关的参数容器的类型。
type DataArgs struct {
	// ReqBufferCap 代表请求缓冲器的容量。
	ReqBufferCap uint32 `json:"req_buffer_cap"`
	// ReqMaxBufferNumber 代表请求缓冲器的最大数量。
	ReqMaxBufferNumber uint32 `json:"req_max_buffer_number"`
	// RespBufferCap 代表响应缓冲器的容量。
	RespBufferCap uint32 `json:"resp_buffer_cap"`
	// RespMaxBufferNumber 代表响应缓冲器的最大数量。
	RespMaxBufferNumber uint32 `json:"resp_max_buffer_number"`
	// ItemBufferCap 代表条目缓冲器的容量。
	ItemBufferCap uint32 `json:"item_buffer_cap"`
	// ItemMaxBufferNumber 代表条目缓冲器的最大数量。
	ItemMaxBufferNumber uint32 `json:"item_max_buffer_number"`
	// ErrorBufferCap 代表错误缓冲器的容量。
	ErrorBufferCap uint32 `json:"error_buffer_cap"`
	// ErrorMaxBufferNumber 代表错误缓冲器的最大数量。
	ErrorMaxBufferNumber uint32 `json:"error_max_buffer_number"`
}


***一个缓冲池会包含若干个缓冲器***,两者都实现了并发安全+队列式的数据传输功能,但是前者是***可伸缩***的。

// ModuleArgs 代表组件相关的参数容器的类型。
type ModuleArgs struct {
	// Downloaders 代表下载器列表。
	Downloaders []module.Downloader
	// Analyzers 代表分析器列表。
	Analyzers []module.Analyzer
	// Pipelines 代表条目处理管道管道列表。
	Pipelines []module.Pipeline
}

有了这些参数,网络爬虫程序就可以正常启动了。不过,拿到这些参数时,需要做的第一件事就是必须检查它们的有效性。为了让这类参数都必须提供检查的方法,我编写了一个接口类型,并让上述3个类型都实现它:

// Args 代表参数容器的接口类型。
type Args interface {
	// Check 用于自检参数的有效性。
	// 若结果值为nil，则说明未发现问题，否则就意味着自检未通过。
	Check() error
}

对于RequestArg类型的值来说,若AcceptDomains字段的值为nil,说明参数无效。对于DataArgs类型的值来说,任何字段的值都不能为0,而对于ModuleArgs类型的值来说,3种组件的实例必须至少提供一个。

Scheduler接口的实现实例需要通过上述这些参数正确设置自己的状态,并未启动做好准备,一旦初始化成功,就可以调用它的Start方法以启动调度器。Start方法只接收一个参数--首次请求,一旦满足这最后一个必要条件,调度器就可以按照既定流程运转起来了。

***调度器的状态***:

// Status 代表调度器状态的类型。
type Status uint8

const (
	// SCHED_STATUS_UNINITIALIZED 代表未初始化的状态。
	SCHED_STATUS_UNINITIALIZED Status = 0
	// SCHED_STATUS_INITIALIZING 代表正在初始化的状态。
	SCHED_STATUS_INITIALIZING Status = 1
	// SCHED_STATUS_INITIALIZED 代表已初始化的状态。
	SCHED_STATUS_INITIALIZED Status = 2
	// SCHED_STATUS_STARTING 代表正在启动的状态。
	SCHED_STATUS_STARTING Status = 3
	// SCHED_STATUS_STARTED 代表已启动的状态。
	SCHED_STATUS_STARTED Status = 4
	// SCHED_STATUS_STOPPING 代表正在停止的状态。
	SCHED_STATUS_STOPPING Status = 5
	// SCHED_STATUS_STOPPED 代表已停止的状态。
	SCHED_STATUS_STOPPED Status = 6
)

只有已初始化的调度器才能被启动,只有已启动的调度器才能被停止。
另一方面,允许重新初始化操作使得调度器可被复用。调度器处于未初始化,已初始化或者停止状态,都可以重新初始化。

***ErrorChan方法***用于获得错误通道。注意:其结果类型是<-chan error,一个只允许接收操作的单向类型通道。调度器会把运行期间发生的绝大部分错误都封装成错误值并传入这个错误通道。调度器的使用方在启动它之后立即调用ErrorChan方法并不断地尝试从其结果值中获取错误值

//省略部分代码
sched := NewScheduler()
err = sched.Init(
    requestArgs,
    dataArgs,
    moduleArgs)
if err != nil {
    t.Fatalf("An error occurs when initializing scheduler: %s",
        err)
}
err = sched.Start(firstHTTPReq)
if err != nil {
    t.Fatalf("An error occurs when starting scheduler: %s",
        err)
}
// 观察错误。
go func() {
    errChan := sched.ErrorChan()
    for {
        err, ok := <-errChan
        if !ok {
            break
        }
        t.Errorf("An error occurs when running scheduler: %s", err)
    }
}()


***Idle方法***的作用是判断调度器当前是否空闲的。判断标准是调度器使用的所有组件都正处于空闲,并且4个缓冲池中也没有任何数据。这样的判断可以依靠组件和缓冲池提供的方法来实现。

最后,***Summary***方法会返回描述调度器当时的内部状态的摘要,与组件接口的Summary方法相同,这里返回的也不是字符串形式的摘要信息,而是返回了承载了调度器摘要信息的SchedSummary类型值。

// SchedSummary 代表调度器摘要的接口类型。
type SchedSummary interface {
	// Struct 用于获得摘要信息的结构化形式。
	Struct() SummaryStruct
	// String 用于获得摘要信息的字符串形式。
	String() string
}


// SummaryStruct 代表调度器摘要的结构。
type SummaryStruct struct {
	RequestArgs     RequestArgs             `json:"request_args"`
	DataArgs        DataArgs                `json:"data_args"`
	ModuleArgs      ModuleArgsSummary       `json:"module_args"`
	Status          string                  `json:"status"`
	Downloaders     []module.SummaryStruct  `json:"downloaders"`
	Analyzers       []module.SummaryStruct  `json:"analyzers"`
	Pipelines       []module.SummaryStruct  `json:"pipelines"`
	ReqBufferPool   BufferPoolSummaryStruct `json:"request_buffer_pool"`
	RespBufferPool  BufferPoolSummaryStruct `json:"response_buffer_pool"`
	ItemBufferPool  BufferPoolSummaryStruct `json:"item_buffer_pool"`
	ErrorBufferPool BufferPoolSummaryStruct `json:"error_buffer_pool"`
	NumURL          uint64                  `json:"url_number"`
}


### 5) 工具箱简述
缓冲池、缓冲器、和多重读取器

#### 5.1 缓冲池和缓冲器
缓冲池和缓冲器是一对程序实体,缓冲器是缓冲池的底层支持。缓冲器是缓冲池的底层支持,缓冲池是缓冲器的再封装。缓冲池利用它持有的缓冲器实现***数据存取***的功能,并可以根据情况***自动地增减它持有的缓冲器的数量***。

// Pool 代表数据缓冲池的接口类型。
type Pool interface {
	// BufferCap 用于获取池中缓冲器的统一容量。
	BufferCap() uint32
	// MaxBufferNumber 用于获取池中缓冲器的最大数量。
	MaxBufferNumber() uint32
	// BufferNumber 用于获取池中缓冲器的数量。
	BufferNumber() uint32
	// Total 用于获取缓冲池中数据的总数。
	Total() uint64
	// Put 用于向缓冲池放入数据。
	// 注意！本方法应该是阻塞的。
	// 若缓冲池已关闭则会直接返回非nil的错误值。
	Put(datum interface{}) error
	// Get 用于从缓冲池获取数据。
	// 注意！本方法应该是阻塞的。
	// 若缓冲池已关闭则会直接返回非nil的错误值。
	Get() (datum interface{}, err error)
	// Close 用于关闭缓冲池。
	// 若缓冲池之前已关闭则返回false，否则返回true。
	Close() bool
	// Closed 用于判断缓冲池是否已关闭。
	Closed() bool
}

Put方法和Get方法需要实现缓冲池最核心的功能——数据的存入和读出。对于这样的操作,在缓冲池关闭之后是不成功的。这时总是返回非nil的错误值。另外,这两个方法都是阻塞的。当缓冲池已满时,对Put方法的调用会产生阻塞。当缓冲池已空闲的时,对Get方法的调用总会产生阻塞。这遵从通道类型的行为模式。

#### pay attation to
如果缓冲池只持有固定数量的缓冲器,那么它的实现就会变得非常简单,基本只利用缓冲器的方法实现功能就可以了。不过这样的话,再封装一层就没有什么意义了。缓冲池这一层的核心功能恰恰就是***动态伸缩***。

对于一个固定容量的缓冲来说,缓冲器可以完全胜任,用不着缓冲池。并且,缓冲池只需要做到这种程度。这样足够简单。更高级的功能全部留给像缓冲池那样的高层类型去做。

***缓冲器***的接口是这样的:

// Buffer 代表FIFO的缓冲器的接口类型。
type Buffer interface {
	// Cap 用于获取本缓冲器的容量。
	Cap() uint32
	// Len 用于获取本缓冲器中的数据数量。
	Len() uint32
	// Put 用于向缓冲器放入数据。
	// 注意！本方法应该是非阻塞的。
	// 若缓冲器已关闭则会直接返回非nil的错误值。
	Put(datum interface{}) (bool, error)
	// Get 用于从缓冲器获取器。
	// 注意！本方法应该是非阻塞的。
	// 若缓冲器已关闭则会直接返回非nil的错误值。
	Get() (interface{}, error)
	// Close 用于关闭缓冲器。
	// 若缓冲器之前已关闭则返回false，否则返回true。
	Close() bool
	// Closed 用于判断缓冲器是否已关闭。
	Closed() bool
}

注意,这里的Put和Get方法与缓冲池的对应方法在行为上有一点不同,即前者是非阻塞的。当缓冲器已满时,Put方法的第一个结果值就会是false。当缓冲器已空时,Get方法的第一个结果值一定会是nil。这样做也是为了让缓冲器的实现保持足够简单。

你可能会有一个疑问,缓冲器的功能看似与通道类型就可以满足。为什么还需要再造一个类型出来呢?在讲通道类型的时候,强调过两个会引发运行时恐慌的操作:向一个已关闭的通道发送值和关闭一个已关闭的通道。实际上,缓冲器接口及其实现就是为了解决这两个问题而存在的。在Put方法中,我会检查当前缓冲器实例是否已关闭,并且保证只有在检查结果为是的时候才进行存入操作。在Close方法中,我仅会在当前缓存器实例未关闭的情况下进行相关操作。另外,我们无法知道一个通道是否已关闭,这也是导致上述第二个引发运行时恐慌的情况发生的最关键的原因。有了Closed方法,我们就可以知道缓冲器的关闭状态,问题也就迎刃而解了。

#### 5.2 多重读取器
如果你知道io.Reader接口并使用过它的实现类型(bytes.Reader、bufio.Reader等)的话,就肯定会知道通过这类读取器只能读取一遍它们持有的底层数据。当读完底层数据时,它们的Read方法总会把io.EOF变量的值作为错误值返回。另外,如果你使用net/http包中的程序实体编写过Web程序的话,还应该知道http.Response类型的Body字段是io.ReadCloser接口类型的,而且该接口的类型声明中嵌入了io.Reader接口。前者只是比后者多声明了一个名为Close的方法。相同的是,当HTTP响应从远程服务器返回并封装成*http.Response类型的值后,你只能通过它的Body字段的值读取HTTP响应体。

#### pay attation to
这些特性本身没有什么问题,但是在我对分析器的设计中,这样的读取器会造成一些小麻烦。还记得吗?一个分析器实例可以持有多个响应解析函数。由于Body字段值的上述特性,如果第一个函数通过它读取了HTTP响应体,那么之后的函数就再也读不到这个HTTP响应体了。响应函数解析函数的一个很重要的职责就是分析HTTP响应体并从中筛选出可用的部分。所以,如此一来,后面的函数就无法实现主要的功能了。

你也许会想到,分析器可以先读出HTTP响应体并赋予给一个[]byte类型的变量,然后把它作为参数直接传给多个响应解析函数。这是可行的,但是我认为这样做会让代码变得丑陋,因为这个值在内容方面与ParseResponse函数类型的第一个参数有所重叠。更为关键的是,这回改变ParseResponse函数类型的声明,这并不值得。

我的做法是,设计一个***可以多次提供基于同一底层数据(可以是[]byte类型的)io.ReadCloser类型值的类型***。我把这个类型命名为MultipleReader,意为***多重读取器***。


// MultipleReader 代表多重读取器的接口。
type MultipleReader interface {
	// Reader 用于获取一个可关闭读取器的实例。
	// 后者会持有本多重读取器中的数据。
	Reader() io.ReadCloser
}


在创建这个类型的值时,我们可以把HTTP相应的Body字段作为参数传入。作为产出,我们可以通过它的Reader方法多次获取基于同一个HTTP响应体的读取器。这些读取器除了基于同一底层数据之外毫不相干。这样一来,我们就可以让多个响应解析函数的分析筛选操作完全独立、互不影响了。

之所以让这个Reader方法返回io.ReadCloser类型的值,是因为我们要用这个值替换HTTP响应原有的Body字段值,这样做是为了让这一改进对响应解析函数透明。也就是说,不让响应解析函数感知到分析器中所作的改变。




# 6.5 工具的实现
## 6.5.1 缓冲器

缓冲器的基本结构如下:
// myBuffer 代表缓冲器接口的实现类型。
type myBuffer struct {
	// ch 代表存放数据的通道。
	ch chan interface{}
	// closed 代表缓冲器的关闭状态：0-未关闭；1-已关闭。
	closed uint32
	// closingLock 代表为了消除因关闭缓冲器而产生的竞态条件的读写锁。
	closingLock sync.RWMutex
}

显然,缓冲器的实现就是对通道类型的简单封装,只不过增加了两个字段用于解决前面所说的问题。字段closed用于标识缓冲器的状态。缓冲器自创建之后只有两种状态:未关闭和关闭。注意:我们需要用***原子操作***访问该字段的值。closingLock字段代表了读写锁。如果你在程序中并发的进行向通道发送值和关闭该通道的操作的话,会产生竞态条件。

// NewBuffer 用于创建一个缓冲器。
// 参数size代表缓冲器的容量。
func NewBuffer(size uint32) (Buffer, error) {
	if size == 0 {
		errMsg := fmt.Sprintf("illegal size for buffer: %d", size)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	return &myBuffer{
		ch: make(chan interface{}, size),
	}, nil
}

### Put方法的实现
func (buf *myBuffer) Put(datum interface{}) (ok bool, err error) {
	buf.closingLock.RLock()
	defer buf.closingLock.RUnlock()
	if buf.Closed() {
		return false, ErrClosedBuffer
	}
	select {
	case buf.ch <- datum:
		ok = true
	default:
		ok = false
	}
	return
}

select语句主要是为了让Put方法永远不会阻塞在发送操作上,在default分支中把结果变量ok的值设置为false,加之这时的结果变量err必为ni,就可以告知调用方放入数据的操作未成功,且原因并不是缓冲器已关闭,而是缓冲器已满。

###Get方法的实现
Get方法的实现要简单一些,因为从通道接收值的操作可以丝毫不受到通道关闭的影响,所以无需加锁。

func (buf *myBuffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-buf.ch:
		if !ok {
			return nil, ErrClosedBuffer
		}
		return datum, nil
	default:
		return nil, nil
	}
}


###Close方法的实现
再说Close方法,在关闭通道之前,先要避免重复操作。因为重复关闭一个通道也会引发运行时恐慌。***避免措施就是先检查closed字段的值。当然,必须使用原子操作***。

func (buf *myBuffer) Close() bool {
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) {
		buf.closingLock.Lock()
		close(buf.ch)
		buf.closingLock.Unlock()
		return true
	} 
	return false
}

###Closed方法的实现
在Closed方法中***读取closed字段的值***,也一定要使用***原子操作***

func (buf *myBuffer) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}

#######重点:千万不要假设读取共享资源就是并发安全的,除非资源本身做出了这种保证。


## 6.5.2 缓冲池

// myPool 代表数据缓冲池接口的实现类型。
type myPool struct {
	// bufferCap 代表缓冲器的统一容量。
	bufferCap uint32
	// maxBufferNumber 代表缓冲器的最大数量。
	maxBufferNumber uint32
	// bufferNumber 代表缓冲器的实际数量。
	bufferNumber uint32
	// total 代表池中数据的总数。
	total uint64
	// bufCh 代表存放缓冲器的通道。
	bufCh chan Buffer
	// closed 代表缓冲池的关闭状态：0-未关闭；1-已关闭。
	closed uint32
	// lock 代表保护内部共享资源的读写锁。
	rwlock sync.RWMutex
}

前两个字段用于记录创建缓冲池时的参数,它们在缓冲池运行期间用到 bufferNumber和total字段用于记录缓冲数据的实时情况。

### pay attation to:
注意:bufCh字段的类型是chan Buffer,一个元素类型为Buffer的通道类型。这与缓冲器同样是通道类型的ch字段联合起来看,就是一个***双层通道***的设计。***在放入或获取数据时,我会先从bufCh拿到一个缓冲器,再向该缓冲器放入数据或从该缓冲器获取数据,然后再把它发送回bufCh***。这样的设计有如下几点好处:
- bufCh中的每个缓冲器一次只会被一个goroutine中的程序(以下简称并发程序)拿到。并且,在放回bufCh之前,它对其他并发程序都是不可见的。一个缓冲器每次只会被并发程序放入或取走一个数据。即使同一个程序连续调用多次Put方法或Get方法,也会这样。缓冲器不至于一下被填满或取空。
- 更进一步看,bufCh是FIFO的。当把先前拿出的缓冲器归还给bufCh时,该缓冲器总会被放在队尾。也就是说,池中缓冲器的操作频率可以降到最低,这也有利于池中数据的均匀分布。
- 在从bufCh拿到缓冲器后,我可以判断是否需要缩减缓冲器的数量。如果需要并且该缓冲器已空,就可以直接把它关掉,并且不还给bufCh。另一方面,如果在放入数据时发现所有缓冲器都已满并且在一段时间内就没有空位,就可以新建一个缓冲器并放入bufCh。总之,这让缓冲池***自动伸缩功能***的实现变得简单了。
- 最后也是最重要的是,bufCh本身提供了对并发安全的保障。

### 需要研究:
你可能会想到,基于标准库的container包中的List或Ring类型也可以编写并发安全的缓冲器队列。确实可以,不过,用它们来实现会让你不得不编写更多的代码。因为原本一些现成的操作和功能都需要我们自己去实现,尤其是在保证并发安全性方面。并且,这样的缓冲器队列的运行效率可不一定高。

注意:上述设计会导致缓冲池中的数据不是FIFO的。不过,对于网络爬虫框架以及调度器来说,这并不会造成问题。

再看最后一个字段rwlock，之所以不叫它closingLock,是因为它不仅仅为了消除缓冲器的那个与关闭通道有关的竞态条件而存在。你可以思考一下,怎样并发的向bufCh放入新的缓冲器,同时避免池中的缓冲器数量超过最大值。

对bufferNumber和total字段的访问需要使用***原子操作***。

#### Put方法
Put方法有两个主要的功能:
- 向缓冲池中放入数据
- 当发现所有的缓冲器都已满一段时间后,新建一个缓冲器并将其放入缓冲池。当然,如果当前缓冲池持有的缓冲器已达最大数量,就不能这么做了。所以,这里我们首先需要建立一个***发现和触发追加缓冲器操作的机制***。我规定当对池中所有缓冲器的操作的失败次数都达到5次时,就追加一个缓冲器入池。


func (pool *myPool) Put(datum interface{}) (err error) {
	if pool.Closed() {
		return ErrClosedBufferPool
	}
	var count uint32
	maxCount := pool.BufferNumber() * 5
	var ok bool
	for buf := range pool.bufCh {
		ok, err = pool.putData(buf, datum, &count, maxCount)
		if ok || err != nil {
			break
		}
	}
	return
}

实际上,放入操作的核心逻辑在myPool类型的putData方法中。Put方法本身做的主要是不断的取出池中的缓冲器,并持有一个统一的***"已满"***计数。请注意count和maxCount变量的初始值。

#### PutData方法

func (pool *myPool) putData(
	buf Buffer, datum interface{}, count *uint32, maxCount uint32) (ok bool, err error) {
	...省略代码
}

##### 第一段
putData为了及时响应缓冲池的关闭,需要在一开始就检***查缓冲池的状态***。并且在方法执行结束前还要检查一次,以便***及时释放资源***。

if pool.Closed() {
	return false, ErrClosedBufferPool
}
defer func() {
	pool.rwlock.RLock()
	if pool.Closed() {
		atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
		err = ErrClosedBufferPool
	} else {
		pool.bufCh <- buf
	}
	pool.rwlock.RUnlock()
}()


##### 第二段 
执行向拿到的缓冲器放入数据的操作,并在必要时增加***已满***计数:

	ok, err = buf.Put(datum)
	if ok {
		atomic.AddUint64(&pool.total, 1)
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已满而未放入数据就递增计数。
	(*count)++


请注意那两条return语句以及最后的(*count)++。在试图向缓冲器放入数据后,我们需要立即判断操作结果。如果ok的值是true,就说明放入成功,此时就可以在递增total字段的值后立即返回。如果err的值不为nil,就是说缓冲器已关闭,这时就不需要再执行后面的语句了。除了这两种情况,我们就需要递增count的值。因为这时说明缓冲器已满。

这里的count值递增操作与第三段代码息息相关,这涉及对追加缓冲器的操作的触发。
	// 如果尝试向缓冲器放入数据的失败次数达到阈值，
	// 并且池中缓冲器的数量未达到最大值，
	// 那么就尝试创建一个新的缓冲器，先放入数据再把它放入池。
	if *count >= maxCount &&
		pool.BufferNumber() < pool.MaxBufferNumber() {
		pool.rwlock.Lock()
		if pool.BufferNumber() < pool.MaxBufferNumber() {
			if pool.Closed() {
				pool.rwlock.Unlock()
				return
			}
			newBuf, _ := NewBuffer(pool.bufferCap)
			newBuf.Put(datum)
			pool.bufCh <- newBuf
			atomic.AddUint32(&pool.bufferNumber, 1)
			atomic.AddUint64(&pool.total, 1)
			ok = true
		}
		pool.rwlock.Unlock()
		*count = 0
	}
	return

在这段代码中,我用到了***双检锁***。如果第一次条件判断通过,就会立即再做一次条件判断。不过这之前,我会先锁定rwlock的写锁。这有两个作用:第一,防止向已关闭的缓冲池追加缓冲器。第二,防止缓冲器的数量超过最大值。在确保这两种情况不会发生后,我就会把一个已放入那个数据的缓冲器追加到缓冲池中。


#### Get方法
Get方法的总体流程与Put方法基本一致:

func (pool *myPool) Get() (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedBufferPool
	}
	var count uint32
	maxCount := pool.BufferNumber() * 10
	for buf := range pool.bufCh {
		datum, err = pool.getData(buf, &count, maxCount)
		if datum != nil || err != nil {
			break
		}
	}
	return
}

我把"已空"计数的上线maxCount设为缓冲器数量的10倍。也就是说,若在遍历所有缓冲器10次之后仍无法获取到数据。Get方法就会从缓冲池中去掉一个空的缓冲器。

#### getData方法
getData方法声明如下:

// getData 用于从给定的缓冲器获取数据，并在必要时把缓冲器归还给池。
func (pool *myPool) getData(
	buf Buffer, count *uint32, maxCount uint32) (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedBufferPool
	}
	defer func() {
		// 如果尝试从缓冲器获取数据的失败次数达到阈值，
		// 同时当前缓冲器已空且池中缓冲器的数量大于1，
		// 那么就直接关掉当前缓冲器，并不归还给池。
		if *count >= maxCount &&
			buf.Len() == 0 &&
			pool.BufferNumber() > 1 {
			buf.Close()
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			*count = 0
			return
		}
		pool.rwlock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = ErrClosedBufferPool
		} else {
			pool.bufCh <- buf
		}
		pool.rwlock.RUnlock()
	}()
	datum, err = buf.Get()
	if datum != nil {
		atomic.AddUint64(&pool.total, ^uint64(0))
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已空未取出数据就递增计数。
	(*count)++
	return
}

#### Close方法
func (pool *myPool) Close() bool {
	if !atomic.CompareAndSwapUint32(&pool.closed, 0, 1) {
		return false
	}
	pool.rwlock.Lock()
	defer pool.rwlock.Unlock()
	close(pool.bufCh)
	for buf := range pool.bufCh {
		buf.Close()
	}
	return true
}

#### Closed方法
func (pool *myPool) Closed() bool {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return true
	}
	return false
}


## 6.5.3 多重读取器

// myMultipleReader 代表多重读取器的实现类型。
type myMultipleReader struct {
	data []byte
}

// NewMultipleReader 用于新建并返回一个多重读取器的实例。
func NewMultipleReader(reader io.Reader) (MultipleReader, error) {
	var data []byte
	var err error
	if reader != nil {
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("multiple reader: couldn't create a new one: %s", err)
		}
	} else {
		data = []byte{}
	}
	return &myMultipleReader{
		data: data,
	}, nil
}

func (rr *myMultipleReader) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(rr.data))
}

***ioutil.NopCloser***:通常用这个函数包装无需关闭的读取器,这就是NopCloser的含义。

多重读取器的Reader方法总是返回一个新的可关闭读取器。因此,我可以利用它很多次读取底层数据,并可以用该方法的结果值替代原先HTTP响应的Body字段值很多次。这也就是"多重"的真正含义。


# 六、组件的实现
网络爬虫框架中的组件有3个:
- 下载器
- 分析器
- 条目处理通道

它们有很多共同点,比如:处理计数的记录、摘要信息的生成和评分以及计算方式的设定。
我应该在组件接口和实现类型之间再抽出一层,用以实现组件的这些通用功能。


## 6.1 内部基础接口
首先要做的是,先为组件通用功能定义一个内部接口,这里把它叫做***组件的内部基础接口***。

// myModule 代表组件内部基础接口的实现类型。
type myModule struct {
	// mid 代表组件ID。
	mid module.MID
	// addr 代表组件的网络地址。
	addr string
	// score 代表组件评分。
	score uint64
	// scoreCalculator 代表评分计算器。
	scoreCalculator module.CalculateScore
	// calledCount 代表调用计数。
	calledCount uint64
	// acceptedCount 代表接受计数。
	acceptedCount uint64
	// completedCount 代表成功完成计数。
	completedCount uint64
	// handlingNumber 代表实时处理数。
	handlingNumber uint64
}

Module接口中声明的更多的是获取内部状态的方法,比如:获取组件ID、组件地址、各种计数值,等等。而在ModuleInternal接口中,我添加的方法都是改变内部状态的方法。由于通常情况下外部不应该直接改变组件的内部状态,所以该接口的名字才以"Internal"为后缀,以起到提示的作用。并且,在module包中公开的程序实体并没有涉及该接口。ModuleInternal接口及其实现类型只是为了方便自行编写组件的人而准备的。我自己在编写组件时也用到了它们。

ModuleInternal接口是Module接口的扩展,前者的实现类型自然也是后者的实现类型。我把实现类型命名为myModule。
它的基本结构如下:

// myModule 代表组件内部基础接口的实现类型。
type myModule struct {
	// mid 代表组件ID。
	mid module.MID
	// addr 代表组件的网络地址。
	addr string
	// score 代表组件评分。
	score uint64
	// scoreCalculator 代表评分计算器。
	scoreCalculator module.CalculateScore
	// calledCount 代表调用计数。
	calledCount uint64
	// acceptedCount 代表接受计数。
	acceptedCount uint64
	// completedCount 代表成功完成计数。
	completedCount uint64
	// handlingNumber 代表实时处理数。
	handlingNumber uint64
}

// NewModuleInternal 用于创建一个组件内部基础类型的实例。
func NewModuleInternal(
	mid module.MID,
	scoreCalculator module.CalculateScore) (ModuleInternal, error) {
	parts, err := module.SplitMID(mid)
	if err != nil {
		return nil, errors.NewIllegalParameterError(
			fmt.Sprintf("illegal ID %q: %s", mid, err))
	}
	return &myModule{
		mid:             mid,
		addr:            parts[2],
		scoreCalculator: scoreCalculator,
	}, nil
}

myModule类型中的字段有几个是需要显式初始化的,包括:组件ID、组件的网络地址(下面简称组件地址)和组件评分计算器。参数mid提供了组件ID,同时也提供了组件地址。因为组件ID中可以包含组件地址。如果组件地址为空,就说明该组件与网络爬虫程序同处于在一个进程中。这时的addr字段自然就是""。

与之对应的,module包中还有一个GenMID函数,用它可以生成组件ID。调用GenMID函数时,需要给定一个序列号。你可以通过调用module包中的NewSNGenertor函数创建出一个***序列号生成器***。***强烈建议把序列号生成器的实例赋给一个全局变量。***


组件评分计算器理应由外部提供,并且一般会为同一类组件实例设置同一种组件评分计算器,而且一旦设置就不允许更改。所以,即使是ModuleInternal接口也没有提供改变它的方法。

有了上述的那些字段,实现ModuleInternal接口的方法就相当简单了,唯一要注意的就是充分利用原子操作***保证它们的并发安全***。

## 6.2 组件注册器

// myRegistrar 代表组件注册器的实现类型。
type myRegistrar struct {
	// moduleTypeMap 代表组件类型与对应组件实例的映射。
	moduleTypeMap map[Type]map[MID]Module
	// rwlock 代表组件注册专用读写锁。
	rwlock sync.RWMutex
}

在组件注册器的实现类型myRegistrar中只有两个字段,一个用于分类存储组件实例,一个用于读写保护。

### Register方法:

Registrar接口的Register方法只需要做两件事情:
- 检查参数
- 注册实例

func (registrar *myRegistrar) Register(module Module) (bool, error) {
	if module == nil {
		return false, errors.NewIllegalParameterError("nil module instance")
	}
	mid := module.ID()
	parts, err := SplitMID(mid)
	if err != nil {
		return false, err
	}
	moduleType := legalLetterTypeMap[parts[0]]
	if !CheckType(moduleType, module) {
		errMsg := fmt.Sprintf("incorrect module type: %s", moduleType)
		return false, errors.NewIllegalParameterError(errMsg)
	}
	// 省略部分代码
}

如果所有检查都通过了,那么Register方法就会把组件实例存储在moduleTypeMap中。当然,我肯定会在rwlock的保护之下操作moduleTypeMap。

### Unregister方法:

会把与给定的组件ID对应的组件实例从moduleTypeMap删除掉。在真正进行查找和删除前,它会先通过调用SplitMID函数检查那个组件ID的合法性。

### Get方法:

Get方法包含***负载均衡策略***,并返回***最"空闲"***的那个组件实例

// Get 用于获取一个指定类型的组件的实例。
// 本函数会基于负载均衡策略返回实例。
func (registrar *myRegistrar) Get(moduleType Type) (Module, error) {
	modules, err := registrar.GetAllByType(moduleType)
	if err != nil {
		return nil, err
	}
	minScore := uint64(0)
	var selectedModule Module
	for _, module := range modules {
		SetScore(module)
		if err != nil {
			return nil, err
		}
		score := module.Score()
		if minScore == 0 || score < minScore {
			selectedModule = module
			minScore = score
		}
	}
	return selectedModule, nil
}


### New
// NewRegistrar 用于创建一个组件注册器的实例。
func NewRegistrar() Registrar {
	return &myRegistrar{
		moduleTypeMap: map[Type]map[MID]Module{},
	}
}


## 6.3 下载器
与ModuleInternal接口一样,下载器接口Downloader也内嵌了Module接口,它额外声明了一个Download方法。有了ModuleInternale接口及其实现类型,实现下载器时只需关注它的特色功能,其他的都交给内嵌的stub.ModuleInternal就可以了。

下载器的实现类型名为myDownloader

// myDownloader 代表下载器的实现类型。
type myDownloader struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// httpClient 代表下载用的HTTP客户端。
	httpClient http.Client
}

可以看到,我匿名的嵌入了一个stub.ModuleInternal类型的字段,***这种只有类型而没有名称的字段称为"匿名字段"***。如此一来,myDownloader类型的方法集合中就包含了stub.ModuleInternal类型的所有方法。因而,*myDownloader类型已经实现了Module接口。

// New 用于创建一个下载器实例。
func New(
	mid module.MID,
	client *http.Client,
	scoreCalculator module.CalculateScore) (module.Downloader, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, genParameterError("nil http client")
	}
	return &myDownloader{
		ModuleInternal: moduleBase,
		httpClient:     *client,
	}, nil
}


### Download方法
func (downloader *myDownloader) Download(req *module.Request) (*module.Response, error) {
	downloader.ModuleInternal.IncrHandlingNumber()
	defer downloader.ModuleInternal.DecrHandlingNumber()
	downloader.ModuleInternal.IncrCalledCount()
	if req == nil {
		return nil, genParameterError("nil request")
	}
	httpReq := req.HTTPReq()
	if httpReq == nil {
		return nil, genParameterError("nil HTTP request")
	}
	downloader.ModuleInternal.IncrAcceptedCount()
	logger.Infof("Do the request (URL: %s, depth: %d)... \n", httpReq.URL, req.Depth())
	httpResp, err := downloader.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	downloader.ModuleInternal.IncrCompletedCount()
	return module.NewResponse(httpResp, req.Depth()), nil
}

这个方法的功能实现起来很简单,不过要注意对那4个组件计数的操作。在方法的开始处,要递增实时处理数,并利用defer语句保证方法执行结束时递减这个计数。同时,还要递增调用计数。在所有参数检查都通过后,要递增接受计数以表明该方法接受了这次调用。一旦目标服务器发回了HTTP响应并且未发生错误,就可以递增成功完成计数了。

## 6.4 分析器
分析器的接口包含两个额外的方法
- RespParsers:前者会返回当前分析器使用的HTTP响应解析函数(以下简称解析函数)的列表
- Analyze


// 分析器的实现类型。
type myAnalyzer struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// respParsers 代表响应解析器列表。
	respParsers []module.ParseResponse
}

// New 用于创建一个分析器实例。
func New(
	mid module.MID,
	respParsers []module.ParseResponse,
	scoreCalculator module.CalculateScore) (module.Analyzer, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if respParsers == nil {
		return nil, genParameterError("nil response parsers")
	}
	if len(respParsers) == 0 {
		return nil, genParameterError("empty response parser list")
	}
	var innerParsers []module.ParseResponse
	for i, parser := range respParsers {
		if parser == nil {
			return nil, genParameterError(fmt.Sprintf("nil response parser[%d]", i))
		}
		innerParsers = append(innerParsers, parser)
	}
	return &myAnalyzer{
		ModuleInternal: moduleBase,
		respParsers:    innerParsers,
	}, nil
}

该函数中的大部分代码都用于参数检查,对参数respParsers的检查要尤为仔细,因为她们一定是网络爬虫框架的使用方提供的,属于外来代码。

### Analyze方法
Analyze方法的功能是:先接收响应并检查,再把HTTP响应一次交给它持有的若干解析函数处理,最后汇总并返回从解析函数那里获得的数据列表和错误列表。

#### 1.检查响应

func (analyzer *myAnalyzer) Analyze(
	resp *module.Response) (dataList []module.Data, errorList []error) {
	analyzer.ModuleInternal.IncrHandlingNumber()
	defer analyzer.ModuleInternal.DecrHandlingNumber()
	analyzer.ModuleInternal.IncrCalledCount()
	if resp == nil {
		errorList = append(errorList,
			genParameterError("nil response"))
		return
	}
	httpResp := resp.HTTPResp()
	if httpResp == nil {
		errorList = append(errorList,
			genParameterError("nil HTTP response"))
		return
	}
	httpReq := httpResp.Request
	if httpReq == nil {
		errorList = append(errorList,
			genParameterError("nil HTTP request"))
		return
	}
	var reqURL = httpReq.URL
	if reqURL == nil {
		errorList = append(errorList,
			genParameterError("nil HTTP request URL"))
		return
	}
	analyzer.ModuleInternal.IncrAcceptedCount()
	respDepth := resp.Depth()
	logger.Infof("Parse the response (URL: %s, depth: %d)... \n",
		reqURL, respDepth)

    //省略部分代码
}


#### 2.解析处理
还记得我们前面讲的多重读取器吗?现在该用到它了:

func (analyzer *myAnalyzer) Analyze(
	resp *module.Response) (dataList []module.Data, errorList []error) {
	//省略部分代码

	// 解析HTTP响应。
	if httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	multipleReader, err := reader.NewMultipleReader(httpResp.Body)
	if err != nil {
		errorList = append(errorList, genError(err.Error()))
		return
	}
	dataList = []module.Data{}
	for _, respParser := range analyzer.respParsers {
		httpResp.Body = multipleReader.Reader()
		pDataList, pErrorList := respParser(httpResp, respDepth)
		if pDataList != nil {
			for _, pData := range pDataList {
				if pData == nil {
					continue
				}
				dataList = appendDataList(dataList, pData, respDepth)
			}
		}
		if pErrorList != nil {
			for _, pError := range pErrorList {
				if pError == nil {
					continue
				}
				errorList = append(errorList, pError)
			}
		}
	}
	if len(errorList) == 0 {
		analyzer.ModuleInternal.IncrCompletedCount()
	}
	return dataList, errorList	
}



## 6.5 条目处理管道
条目处理管道的接口拥有额外的ItemProcesssors、Send、FailFast和SetFailFast方法,因此其实现类型myPipeline的基本结构是这样的:

// myPipeline 代表条目处理管道的实现类型。
type myPipeline struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// itemProcessors 代表条目处理器的列表。
	itemProcessors []module.ProcessItem
	// failFast 代表处理是否需要快速失败。
	failFast bool
}

### Send方法:
func (pipeline *myPipeline) Send(item module.Item) []error {
	pipeline.ModuleInternal.IncrHandlingNumber()
	defer pipeline.ModuleInternal.DecrHandlingNumber()
	pipeline.ModuleInternal.IncrCalledCount()
	var errs []error
	if item == nil {
		err := genParameterError("nil item")
		errs = append(errs, err)
		return errs
	}
	pipeline.ModuleInternal.IncrAcceptedCount()
	logger.Infof("Process item %+v... \n", item)
	var currentItem = item
	for _, processor := range pipeline.itemProcessors {
		processedItem, err := processor(currentItem)
		if err != nil {
			errs = append(errs, err)
			if pipeline.failFast {
				break
			}
		}
		if processedItem != nil {
			currentItem = processedItem
		}
	}
	if len(errs) == 0 {
		pipeline.ModuleInternal.IncrCompletedCount()
	}
	return errs
}


## 6.7 调度器
调度器的主要职责是对各个处理模块进行调度,以使它们能够进行良好的协作并共同完成整个爬取流程。

### 基本结构

// myScheduler 代表调度器的实现类型。
type myScheduler struct {
	// maxDepth 代表爬取的最大深度。首次请求的深度为0。
	maxDepth uint32
	// acceptedDomainMap 代表可以接受的URL的主域名的字典。
	acceptedDomainMap cmap.ConcurrentMap
	// registrar 代表组件注册器。
	registrar module.Registrar
	// reqBufferPool 代表请求的缓冲池。
	reqBufferPool buffer.Pool
	// respBufferPool 代表响应的缓冲池。
	respBufferPool buffer.Pool
	// itemBufferPool 代表条目的缓冲池。
	itemBufferPool buffer.Pool
	// errorBufferPool 代表错误的缓冲池。
	errorBufferPool buffer.Pool
	// urlMap 代表已处理的URL的字典。
	urlMap cmap.ConcurrentMap
	// ctx 代表上下文，用于感知调度器的停止。
	ctx context.Context
	// cancelFunc 代表取消函数，用于停止调度器。
	cancelFunc context.CancelFunc
	// status 代表状态。
	status Status
	// statusLock 代表专用于状态的读写锁。
	statusLock sync.RWMutex
	// summary 代表摘要信息。
	summary SchedSummary
}

// NewScheduler 会创建一个调度器实例。
func NewScheduler() Scheduler {
	return &myScheduler{}
}

一切初始化调度器的工作都交给Init方法去做。

### 初始化

调度器接口中声明的第一个方法就是Init方法,它的功能是初始化当前调度器。

关于Init方法接受的那三个参数,前面已经提到多次。Init方法会对它们进行检查。不过在这之前,它必须***先检查调度器的当前的状态***。

func (sched *myScheduler) Init(
	requestArgs RequestArgs,
	dataArgs DataArgs,
	moduleArgs ModuleArgs) (err error) {
	// 检查状态。
	logger.Info("Check status for initialization...")
	var oldStatus Status
	oldStatus, err =
		sched.checkAndSetStatus(SCHED_STATUS_INITIALIZING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_INITIALIZED
		}
		sched.statusLock.Unlock()
	}()
	// 省略部分代码
}

这里有对状态的两次检查。第一次是在开始处,用于确认当前调度器的状态允许我们进行初始化,这次检查是由调度器的checkAndSetStatus方法执行。该方法会在检查通过后按照我们的意愿设置调度器的状态。

// checkAndSetStatus 用于状态的检查，并在条件满足时设置状态。
func (sched *myScheduler) checkAndSetStatus(
	wantedStatus Status) (oldStatus Status, err error) {
	sched.statusLock.Lock()
	defer sched.statusLock.Unlock()
	oldStatus = sched.status
	err = checkStatus(oldStatus, wantedStatus, nil)
	if err == nil {
		sched.status = wantedStatus
	}
	return
}

下面是其中调用的checkStatus方法声明的片段:

// checkStatus 用于状态的检查。
// 参数currentStatus代表当前的状态。
// 参数wantedStatus代表想要的状态。
// 检查规则：
//     1. 处于正在初始化、正在启动或正在停止状态时，不能从外部改变状态。
//     2. 想要的状态只能是正在初始化、正在启动或正在停止状态中的一个。
//     3. 处于未初始化状态时，不能变为正在启动或正在停止状态。
//     4. 处于已启动状态时，不能变为正在初始化或正在启动状态。
//     5. 只要未处于已启动状态就不能变为正在停止状态。
func checkStatus(
	currentStatus Status,
	wantedStatus Status,
	lock sync.Locker)(err error){
	//省略部分代码
	}

这个方法的注释详细描述了检查规则,这决定了调度器是否能够从当前状态转换到我们想要的状态。只要欲进行的转换违反了这些规则中的一条,该方法就会直接返回一个可以说明状况的错误值,而checkAndSetStatus方法会检查checkStatus方法返回的这个错误值。只有当该值为nil时,它才会对调度器状态进行设置。

实际上,在调度器实现类型的Start方法和Stop方法的开始处,也都有类似的代码,它们共同保证了调度器的动作与状态的协同。


如果当前状态允许初始化,那么Init方法就会开始做***参数检查***。这并不麻烦,因为那3个参数的类型本身都提供了检查自身的方法Check。

func (sched *myScheduler) Init(
	requestArgs RequestArgs,
	dataArgs DataArgs,
	moduleArgs ModuleArgs) (err error) {
	// 省略部分代码
	// 检查参数。
	logger.Info("Check request arguments...")
	if err = requestArgs.Check(); err != nil {
		return err
	}
	logger.Info("Check data arguments...")
	if err = dataArgs.Check(); err != nil {
		return err
	}
	logger.Info("Data arguments are valid.")
	logger.Info("Check module arguments...")
	if err = moduleArgs.Check(); err != nil {
		return err
	}
	logger.Info("Module arguments are valid.")	
	// 省略部分代码
}


在这之后,Init方法就要初始化调度器内部的字段了。关于这些字段的初始化方法,之前都陆续讲过,这里就不再展示了。

最后,我们来看一下用于***组件实例注册***的代码:

func (sched *myScheduler) Init(
	requestArgs RequestArgs,
	dataArgs DataArgs,
	moduleArgs ModuleArgs) (err error) {
	// 省略部分代码
	// 注册组件。
	logger.Info("Register modules...")
	if err = sched.registerModules(moduleArgs); err != nil {
		return err
	}
	logger.Info("Scheduler has been initialized.")
	return nil
}

综上所述,Init方法就做了4件事:
- 检查调度器状态
- 检查参数
- 初始化内部字段
- 注册组件实例


### 启动
调度器接口中用于启动调度器的方法是Start。它只接受一个参数,这个参数是*http.Request类型的,代表调度器在当次启动时需要处理的第一个基于HTTP/HTTPS协议的请求。

Start方法首先要做的是防止启动过程中发生的运行时恐慌。
它还需要检查调度器的状态和使用方提供的参数值,并把首次请求的主域名添加到可接受的主域名的字典。

func (sched *myScheduler) Start(firstHTTPReq *http.Request) (err error) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal scheduler error: %sched", p)
			logger.Fatal(errMsg)
			err = genError(errMsg)
		}
	}()
	logger.Info("Start scheduler...")
	// 检查状态。
	logger.Info("Check status for start...")
	var oldStatus Status
	oldStatus, err =
		sched.checkAndSetStatus(SCHED_STATUS_STARTING)
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_STARTED
		}
		sched.statusLock.Unlock()
	}()
	if err != nil {
		return
	}
	// 检查参数。
	logger.Info("Check first HTTP request...")
	if firstHTTPReq == nil {
		err = genParameterError("nil first HTTP request")
		return
	}
	logger.Info("The first HTTP request is valid.")
	// 获得首次请求的主域名，并将其添加到可接受的主域名的字典。
	logger.Info("Get the primary domain...")
	logger.Infof("-- Host: %s", firstHTTPReq.Host)
	var primaryDomain string
	primaryDomain, err = getPrimaryDomain(firstHTTPReq.Host)
	if err != nil {
		return
	}
	logger.Infof("-- Primary domain: %s", primaryDomain)
	sched.acceptedDomainMap.Put(primaryDomain, struct{}{})
    //省略部分代码  
}

你可以把Start方法和Init方法中检查调度器状态的代码对照起来看,并想象这是一个***状态机***在运转。
