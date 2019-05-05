# 一、原理
Appdash实现了Google dapper中的四个主要概念:

- [Span]

Span指的是一个服务调用的跨度,在实现中用SpanId标识。根服务调用者的Span为根span(root span),在根级别进行的下一级服务调用Span的Parent Span为root span。以此类推,服务调用链构成了一颗tree,整个tree构成了一个Trace。

Appdash中SpanId由三部分组成:TraceId/SpanId/parentSpanId,例如:34c31a18026f61df/aab2a63e86ac0166/592043d0a5871aaf。TraceId用于唯一标识一次Trace。traceid在申请RootSpanID时自动分配。

在上面原理图中,我们也可以看到一次Trace过程中SpanID的情况。图中调用链大致是:

frontservice:
    call serviceA
    call serviceB
                    call serviceB1
    ... ... 
    call serviceN

对应服务调用的Span的树形结构如下：
frontservice: SpanId = xxxxx/nnnn1，该span为root span：traceid=xxxxx, spanid=nnnn1，parent span id为空。
serviceA: SpanId = xxxxx/nnnn2/nnnn1，该span为child span：traceid=xxxxx, spanid=nnnn2，parent span id为root span id:nnnn1。
serviceB: SpanId = xxxxx/nnnn3/nnnn1，该span为child span：traceid=xxxxx, spanid=nnnn3，parent span id为root span id:nnnn1。
… …
serviceN: SpanId = xxxxx/nnnnm/nnnn1，该span为child span：traceid=xxxxx, spanid=nnnnm，parent span id为root span id:nnnn1。
serviceB1: SpanId = xxxxx/nnnn3-1/nnnn3，该span为serviceB的child span，traceid=xxxxx, spanid=nnnn3-1，parent span id为serviceB的spanid：nnnn3

- [Event]
个人理解在Appdash中Event是服务调用跟踪信息的wrapper。最终我们在Appdash UI上看到的信息,都是由event承载的并且发给Appdash Server的信息。在Appdash中,你可以显式使用event埋点,吐出跟踪信息,也可以使用Appdash封装好的包接口,比如httptrace.Transport等发送调用跟踪信息,这些包的底层实现也是基于event的。event在传输前会被encoding为Annotation的形式。

- [Recorder]
在Appdash中,Recorder是用来发送event给Appdash的Collector的,每个Recorder会与一个特定的span相关联。

- [Collector]
从Recorder那接收Annotation（即encoded event）。通常一个appdash server会运行一个Collector，监听某个跟踪信息收集端口，将收到的信息存储在Store中。


# 二、安装
appdash是开源的，通过go get即可得到源码并安装example：
 go get -u sourcegraph.com/sourcegraph/appdash/cmd/…
appdash自带一个example，在examples/cmd/webapp下面。执行webapp，你会看到如下结果：
$webapp
2015/06/17 13:14:55 Appdash web UI running on HTTP :8700
[negroni] listening on :8699
这是一个集appdash server, frontservice, fakebackendservice于一身的example，其大致结构如下图：

通过浏览器打开:localhost:8700页面，你会看到appdash server的UI，通过该UI你可以看到所有Trace的全貌。
访问http://localhost:8699/，你就触发了一次Trace。在appdash server ui下可以看到如下画面：


从页面上展示的信息可以看出，该webapp在处理用户request时共进行了三次服务调用，三次调用的耗时分别为：201ms，202ms， 218ms，共耗时632ms。
一个更复杂的例子在cmd/appdash下面，后面的应用实例也是根据这个改造出来的，这里就不细说了。

# 三、应用实例
这里根据cmd/appdash改造出一个应用appdash的例子，例子的结构如下图：

User --Requet :8080--->FrontService ---Event---:3001-------------->localCollector --->Store--->UI---:3000----> Operator
     <---Resp---------              --:6601--> Backend ServiceA    [    AppDash       Server    ]
                                    --:6602--> Backend ServiceB


例子大致分为三部分：
- appdash — 实现了一个appdash server， 该server带有一个collector，用于收集跟踪信息，收集后的信息存储在一个memstore中；appdash server提供ui，ui从memstore提取信息并展示在ui上供operator查看。
- backendservices — 实现两个模拟的后端服务，供frontservice调用。
- frontservice — 服务调用的起始端，当用户访问系统时触发一次跟踪。




