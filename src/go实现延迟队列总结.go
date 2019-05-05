type DelayMessage struct {
	//当前下标
	curIndex int
	//环形槽
	slots [3600]map[string]*Task
	//关闭
	closed chan bool
	//任务关闭
	taskClose chan bool
	//时间关闭
	timeClose chan bool
	//启动时间
	startTime time.Time
}data interface{}
}

//任务
type Task struct {
	//循环次数
	cycleNum int
	//执行的函数
	exec   TaskFunc
	params []interface{}
}


//创建延迟消息
dm := NewDelayMessage()
//启动延迟消息
dm.Start()
		- //goroutine 处理每1秒的任务
		- go dm.taskLoop()
			- for循环
				- select 
					- case <-dm.taskClose
					- default:
							//取出任务
							- tasks := dm.slots[dm.curIndex]
							- //遍历任务
								- ...
									- //该执行了
									- go v.exec(v.params...)
									- //删除运行过的任务
									- delete(tasks, k)
								- else v.cycleNum--

		- //goroutine 处理每1秒移动下标
		- go dm.timeLoop()
			- tick := time.NewTicker(time.Second)
			- for循环
				- select 
					- case <-dm.timeClose
					- case <-tick.C
						- dm.curIndex++
		- select 
			- case <-dm.closed:
				- dm.taskClose <- true
				- dm.timeClose <- true