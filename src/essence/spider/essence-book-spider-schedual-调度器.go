一."调度器"

1."调度器的model"

type Scheduler interface{
	Init(requestArgs RequestArgs,
		dataArgs DataArgs,
		moduleArgs ModuleArgs)(err error)
	Start(firstHTTPReq *http.Request)(err error)
	Stop()(err error)
	Status()Status
	ErrorChan()<-chan error
	Idle()bool 
	Summary()SchedSummary
}

type myScheduler struct{
	maxDepth uint32
	acceptedDomainMap cmap.ConcurrentMap
	registar model.Registrar
	reqBufferPool buffer.Pool
	respBufferPool buffer.Pool
	itemBufferPool buffer.Pool
	errorBufferPool buffer.Pool
	urlMap cmap.ConcurrentMap
	ctx context.Context
	cancelFunc context.CancelFunc
	status Status 
	statusLock sync.RWMutex
	summary SchedSummary
}

2."调度器的Init方法"
- Init(requestArgs RequestArgs,
		dataArgs DataArgs,
		moduleArgs ModuleArgs)(err error)
	- //检查状态
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
	- //检查参数
	- // 初始化内部字段
			- // 组件注册器
			"想着nil"
			- if sched.registrar==nil{
				sched.registrar=model.NewRegistrar()
			}else{
				sched.registrar.Clear()
			}
			- // 初始化,maxDepth
			sched.maxDepth=requesArgs.MaxDepth
			- // 初始化acceptedDomainMap
			- sched.acceptedDomainMap,_=cmap.NewConcurrentMap(1,nil)
				for _, domain := range requestArgs.AcceptedDomains {
					sched.acceptedDomainMap.Put(domain, struct{}{})
				}
			- //初始化urlMap
			- sched.urlMap,_=cmap.NewConcurrentMap(16,nil)
			- //初始化缓冲池
			- sched.initBufferPool(dataArgs)
					- //初始化请求缓冲池
					"时刻将函数当成考虑外的调用"
					- if sched.reqBufferPool!=nil && !sched.reqBufferPool.Closed(){
						sched.reqBufferPool.Close()
					}
						sched.reqBufferPool,_=buffer.NewPool(dataArgs.ReqBufferCap,dataArgs.ReqMaxBufferNumber)
					- //初始化响应缓冲池。
					- //初始化条目缓冲池。
					- //初始化错误缓冲池。
			- //重置context
			- sched.resetContext()
					- sched.ctx, sched.cancelFunc = context.WithCancel(context.Background())
			- //初始化summary
			- 	sched.summary = newSchedSummary(requestArgs, dataArgs, moduleArgs, sched)
	- // 注册组件
	- sched.registerModules(moduleArgs)
			- // registerModules 会注册所有给定的组件。
			- func (sched *myScheduler) registerModules(moduleArgs ModuleArgs) error 
					- //遍历
					- ok, err := sched.registrar.Register(d)
					- //遍历
					- ok, err := sched.registrar.Register(a)
					- //遍历
					- ok, err := sched.registrar.Register(p)


3."调度器的Start方法"
- Start(firstHTTPReq *http.Request)(err error)
- 	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal scheduler error: %sched", p)
			logger.Fatal(errMsg)
			err = genError(errMsg)
		}
	}()
- // 检查状态
- //检查参数
- // 开始调度数据和组件
	- //检查缓冲池是否已为调度器的启动准备就绪
	  // 如果某个缓冲池不可用，就直接返回错误值报告此情况。
	  // 如果某个缓冲池已关闭，就按照原先的参数重新初始化它。
	- err = sched.checkBufferPoolForStart()
	- // download会从请求缓冲池取出请求并下载
	  // 然后把的得到的相应放入相应缓冲池
	- sched.download()
			- go func(){
				for{
					if sched.canceled(){
						break
					}
					datum, err := sched.reqBufferPool.Get()
					req, ok := datum.(*module.Request)
					// downloadOne 会根据给定的请求执行下载并把响应放入响应缓冲池。
					sched.downloadOne(req)
							- if req==nil 
							- if sched.canceled()
							- // 获取一个downloader
							- m, err := sched.registrar.Get(module.TYPE_DOWNLOADER)
							- downloader, ok := m.(module.Downloader)
							- resp, err := downloader.Download(req)
				}
			}()
	// 同上,analyze 会从响应缓冲池取出响应并解析，然后把得到的条目或请求放入相应的缓冲池。
	- sched.analyze()
	// 同上,pick 会从条目缓冲池取出条目并处理。
	- sched.pick()
- //放入第一个请求
	- // NewRequest 用于创建一个新的请求实例
	- firstReq:=model.NewRequest(firstHTTPReq,0)
	- // sendReq 会向请求缓冲池发送请求。
	- sched.sendReq(firstReq)
			- //检查一堆一堆一堆
			- go func(){
				sched.reqBufferPool.Put(req)
			}()
			- // 放入处理的url map
			- sched.urlMap.Put(reqURL.String(),struct{}{})


4."调度器的Stop方法":

- // 检查状态
-	sched.cancelFunc()
-	sched.reqBufferPool.Close()
-	sched.respBufferPool.Close()
-	sched.itemBufferPool.Close()
-	sched.errorBufferPool.Close()


5."调度器的Status方法":
- 
	var status Status
	sched.statusLock.RLock()
	status = sched.status
	sched.statusLock.RUnlock()
	return status


6."调度器的ErrorChan()方法"
- // 从缓冲池里导出来
- errBuffer:=sched.errorBufferPool
errCh:=make(chan error,errBuffer.BufferCap())
go func(errBuffer buffer.Pool, errCh chan error){
	datum, err := errBuffer.Get()
	err, ok := datum.(error)
	errCh <- err
}(errBuffer, errCh)
return errCh



7."调度器的Idle()方法":
- //取出所有模块
- modelMap:=sched.registrar.GetAll()
- // 遍历看每个模块的HandlingNumber
- // 看缓冲池的Total


8."调度器的Summary()方法"
- return sched.summary

5."状态机"

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

