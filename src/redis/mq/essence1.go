// WorkerPids returns a set of existing workers' ids
func (worker *Worker) WorkerPids() mapset.Set {
	/* Returns a set of all pids (as strings) on
	   this machine.  Used when pruning dead workers. */
	out, err := exec.Command("ps").Output()
	if err != nil {
		log.Fatal(err)
	}

	outLines := strings.Split(strings.TrimSpace(string(out)), "\n")
	inSlice := make([]interface{}, len(outLines)-1) // skip first row
	for i, line := range outLines[1:] {
		inSlice[i] = strings.Split(strings.TrimSpace(line), " ")[0] // pid at index 0
	}

	return mapset.NewSetFromSlice(inSlice)
}




NewDispatcher
// @ dispatcher开启多个worker来进行处理job
// @ dispatcher获取job转发给worker
dispatcher Start(tasks)
- 开启多个worker
	- worker := NewWorker(config, disp.queues, i+1)
	- go 
		//@ worker start
		- err := worker.Start(disp, tasks)
				// 修掉死掉的worker!!!!!!!
				- err := worker.PruneDeadWorkers()
				// 注册worker到set中
				- err = worker.RegisterWorker()
				// 重点
				- worker.work(dispatcher, tasks)
						- for 
							- case job, ok := <-worker.jobChan
							// 这里涉及到
							// job执行
							// 错误重试
							// job stat
							- ExecuteJob(job, tasks)
					
				// 取消注册
				- err = worker.UnregisterWorker()
	- go 
		//@ 源源不断获取job,再分发到worker
		- disp.dispatch(workers)
			- go 
				- for 
					// 获取job 通过BLPOP
					- job, err := ReserveJob(disp.gores, disp.queues)
					// 将job放到dispatch的jobChan
					- disp.jobChan <- job 
			- for 
				- for 
					- range workers
						// 从dispatch的jobChan里拿job放到worker的jobChan中
						- job,ok:=<-disp.jobChan
						- worker.jobChan <- job
	// 等待协程完成
	- wg.Wait()





关于job错误重试,这里实现了一个延时机制
1200是时间
delay_list_1200   task 
delay_zset 1200 1200

data, err := conn.Do("ZRANGEBYSCORE", watchedSchedules, "-inf", key)
key := fmt.Sprintf(delayedQueuePrefix, timeStr)


// 处理延时和失败job
type Scheduler struct {
	gores         *Gores
	timestampChan chan int64 // 存放将要处理的过期时间戳
	shutdownChan  chan bool // 关闭标志chan
}

func (sche *Scheduler) Run()
	//# 开启goroutine不断获取延时时间戳,另一边不断根据延时时间戳获取对应条目进行入列操作
	- sche.HandleDelayedItems()
		- go sche.NextDelayedTimestamps()
			- for
				// 获取下一个需要执行的延时时间戳[利用zset]
				- timestamp := sche.gores.NextDelayedTimestamp()
				// 放入处理时间chan
				- sche.timestampChan <- timestamp 
			// 关闭调度器
			- sche.ScheduleShutdown()
			- for 
				- select case timestamp := <-sche.timestampChan:
					// 根据时间戳获取将要执行的item
					- item := sche.gores.NextItemForTimestamp(timestamp)
					// 将item入列
					- err := sche.gores.Enqueue(item)
				- select case <-sche.shutdownChan:
					- return
				


利用redis存放stat信息