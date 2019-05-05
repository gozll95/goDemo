# essence-book-monitor 

1."设计"
监控器

定时监控

定时汇报结果

阈值: 最大空闲计数


1)"摘要":

type summary struct {
	// NumGoroutine 代表Goroutine的数量。
	NumGoroutine int `json:"goroutine_number"`
	// SchedSummary 代表调度器的摘要信息。
	SchedSummary sched.SummaryStruct `json:"sched_summary"`
	// EscapedTime 代表从开始监控至今流逝的时间。
	EscapedTime string `json:"escaped_time"`
}

2."接收和报告错误"
func reportError(scheduler sched.Scheduler,record Record,stopNotifier context.Context){
	go func(){
		// 等待调度器开启
		waitForSchedulerStart(scheduler)
		errorChan:=scheduler.ErrorChan()
		for{
			select{
			case <-stopNotifer.Done():
				return 
			default:
			}
			err,ok:=<-errorChan 
			if ok{
				errMsg=fmt.Spintf("xxx")
				record(2,errMsg)
			}
			time.Sleep(time.Microsecond)
		}
	}()
}

3."记录摘要信息":

4."检查状态,并在满足持续空闲时间的条件时采取必要措施":
func checkStatus(
	scheduler sched.Scheduler,
	checkInterval time.Duration,
	maxIdleCount int,
	autoStop bool,
	checkCountChan chan<-uint64,
	record Record,
	stopFunc context.CancelFunc
){
	go func(){
		

	}()
}

5."流程":
xxx 
xxx 

// 生成监控停止通知器

// 接受和报告错误
// 记录摘要信息
// 检查空闲状态