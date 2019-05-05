package engine

type ConcurrentEngine struct {
	Scheduler        Scheduler
	WorkerCount      int
	ItemChan         chan Item
	RequestProcessor Processor
}

type Processor func(r Request) (ParseResult, error)

type Scheduler interface {
	ReadyNotifier
	Submit(Request)
	WorkerChan() chan Request
	Run()
}

type ReadyNotifier interface {
	WokerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult)
	e.Scheduler.Run()
	for i := 0; i < e.WorkerCount; i++ {
		//createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
		e.createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	for _, r := range seeds {
		if isDuplicate(r.Url) {
			//log.Printf("Duplicate request: %s", r.Url)
			continue
		}
		e.Scheduler.Submit(r)
	}
	//profileCount := 0
	for {
		result := <-out
		for _, item := range result.Items {
			//if profile, ok := item.(model.Profile); ok {
			//x
			//	log.Printf("Got item #%d: %v", profileCount, profile)
			//}
			go func() {
				e.ItemChan <- item
			}()
			//	profileCount++

		}
		for _, r := range result.Requests {
			if isDuplicate(r.Url) {
				//log.Printf("Duplicate request: %s", r.Url)
				continue
			}
			e.Scheduler.Submit(r)
		}
	}
}

//func createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier)
func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			ready.WokerReady(in)
			request := <-in
			//result, err := worker(request)             // Call RPC
			result, err := e.RequestProcessor(request) // Call RPC
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

var visitedUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}
