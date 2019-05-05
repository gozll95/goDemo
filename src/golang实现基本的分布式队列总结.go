type JobServer struct {
	Workers
	JobQueue []*WorkerClass
	interval time.Duration
	mt       sync.Mutex
	ord      OrdType
}


type WorkerClass struct {
	ClassName string
	Args      []interface{}
}

type OrdType int

const (
	PARALLEL = 1 << iota
	ORDER
)

type workerFunc func(string, ...interface{}) error

type Workers struct {
	workers map[string]workerFunc
}

func mailWorker(queue string, args ...interface{}) error {
	fmt.Println("......mail() begin......")
	for _, arg := range args {
		fmt.Println("   args:", arg)
	}
	fmt.Println("......mail() end......")
	return nil
}


- js := NewJobServer()
- //模拟Worker端注册
// - s.workers: 
// "mail":mailf
// "xxx":xxxf
- js.RegisterWorkerClass("mail", mailWorker)
		- RegisterWorkerClass(className string, f workerFunc)
				- s.mt.Lock()
				- defer s.mt.Unlock() 
				- s.workers[className]=f
- //模拟客户端发送请求
// JobQueue:
// [{className:"mail",args:"xxxx"},{className:"xxx",args:"xxxx"}]
- goroutine
	- js.Enqueue("mail", "xcl_168@aliyun.com", "sub", "body")
			- Enqueue(className string, args ...interface{}) bool
					- s.mt.Lock()
					- w := &WorkerClass{className, args}
					- s.JobQueue = append(s.JobQueue, w)
					- s.mt.Unlock()
//启动服务，开始轮询
- js.StartServer(time.Second*3, ORDER)
		- StartServer(interval time.Duration, ord OrdType)
				- s.interval = interval
				- s.ord = ord
				- //这里是channel func
				- quit := signals() 
						- quit := make(chan bool)
						- goroutine
							- os.Stdin.Read(make([]byte, 1))
							- quit <- true
						- return quit
				- //poller get channel jobs
				- jobs := s.poll(quit)
						- poll(quit <-chan bool) <-chan *WorkerClass
								- jobs := make(chan *WorkerClass) //这里用指针方便通过nil看
								- goroutine
									- for 循环
									- switch case  s.JobQueue == nil:
										- //一些通告信息
										- select case <-quit:
											- return 
									- switch default:
										- s.mt.Lock()
										- j := s.JobQueue[0]
										- s.mt.Unlock()
										- select case job<-j:
										- select case <-quit:
											- retun 
						- return jobs
				- var monitor sync.WaitGroup
				- //顺序执行
				- switch s.ord case ORDER 
					- //range jobs 指定
					- s.work(0, jobs, &monitor)
							- work(id int, jobs <-chan *WorkerClass, monitor *sync.WaitGroup)
									- monitor.Add(1)
									- f:=func(){
										defer monitor.Done()
										for job := range jobs {
											if f, found := s.workers[job.ClassName]; found {
												s.run(f, job)
											} else {
												fmt.Println("[JobServer] [poll] ", job.ClassName, " not found")
											}
										}	
									}
				- //并发执行
				- switch default:
					- ... 
					- s.work(id, jobs, &monitor)
				- monitor.Wait()
					







值得借鉴:
- select case job<-j:
- select case <-quit:


func signals() <-chan bool {
	quit := make(chan bool)
	go func() {
		os.Stdin.Read(make([]byte, 1))
		quit <- true
	}()
	return quit
}

		defer monitor.Done()
		for job := range jobs {
			if f, found := s.workers[job.ClassName]; found {
				s.run(f, job)
			} else {
				fmt.Println("[JobServer] [poll] ", job.ClassName, " not found")
			}
		}
	}