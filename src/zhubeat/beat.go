package main

import (
	"fmt"
	"strconv"
	"time"
	"zhubeat/beat/log"
	"zhubeat/beat/output"
	"zhubeat/beat/queue"
	lib_queue "zhubeat/lib/queue"
)

var Job = []string{
	`{"level":"info","method":"GET","msg":"start","path":"/v1/vm/instance","request_id":"f856b899-a98f-4439-8af9-7e4448d08fa4","source_ip":"10.200.20.22","time_start":"2018-06-03T18:19:34.040842122+08:00"}`,
	`{"file":"github.com/zhu/qvm/server/controller/utils.go:187","func":"github.com/zhu/qvm/server/controller.NewBaseRequest","level":"info","msg":"User Uid: 1810637322","request_id":"f856b899-a98f-4439-8af9-7e4448d08fa4","timedate":"2018-06-03T18:19:34.041840711+08:00"}`,
	`{"file":"github.com/zhu/qvm/server/controller/utils.go:303","func":"github.com/zhu/qvm/server/controller.NewEcsRequestWithRegions","level":"info","msg":"Front Region: []","request_id":"f856b899-a98f-4439-8af9-7e4448d08fa4","timedate":"2018-06-03T18:19:34.041937611+08:00"}`,
	`{"level":"info","msg":"finish","request_id":"f856b899-a98f-4439-8af9-7e4448d08fa4","status":200,"time_finish":"2018-06-03T18:19:34.055402735+08:00","time_latency":"14.556048ms"}`,
	`{"file":"github.com/zhu/qvm/server/lib/bbbb/logger/context.go:26","func":"github.com/zhu/qvm/server/lib/bbbb/logger.(*Logger).Xput","level":"debug","msg":"[TRADE:5]","request_id":"69bb7339-0574-46ca-bc98-6508c45a07b6","timedate":"2018-06-03T18:19:34.082205394+08:00"}`,
	`{"level":"info","msg":"finish","request_id":"69bb7339-0574-46ca-bc98-6508c45a07b6","status":200,"time_finish":"2018-06-03T18:19:34.082408836+08:00","time_latency":"42.5374ms"}`,
	`{"file":"github.com/zhu/qvm/server/lib/cron/cron.go:110","func":"github.com/zhu/qvm/server/lib/cron.(*Cron).addEvent.func1","level":"info","msg":"time: 2018-06-04 00:00:00.001475703 +0800 CST m=+196898.480554588 ,start do event(UpdateSystemImage)","timedate":"2018-06-04T00:00:00.001477069+08:00"}`,
	`{"file":"github.com/zhu/qvm/server/application/system_image.go:74","func":"github.com/zhu/qvm/server/application.UpdateSystemImage.func1","level":"info","msg":"there are 396 images updated totally","timedate":"2018-06-04T00:00:24.259603441+08:00"}`,
	`{"file":"github.com/zhu/qvm/server/lib/cron/cron.go:112","func":"github.com/zhu/qvm/server/lib/cron.(*Cron).addEvent.func1","level":"info","msg":"spend: 24.258213852s ,end do event(UpdateSystemImage)","timedate":"2018-06-04T00:00:24.259690659+08:00"}`,
	`{"file":"github.com/zhu/qvm/server/lib/cron/cron.go:106","func":"github.com/zhu/qvm/server/lib/cron.(*Cron).addEvent.func1","level":"warning","msg":"cron event UpdateOms blocked by error E11000 duplicate key error collection: qvm_prod.qvm_cron index: event_1 dup key: { : \"UpdateOms\" }","timedate":"2018-06-04T01:00:00.003623097+08:00"}`,
}

var Dispatcher *output.Dispatcher

func init() {
	batchLen := 10
	threthod := 10
	ttl := 1 * time.Second
	queue.LogQueue = queue.NewJobQueue(batchLen, threthod, ttl, lib_queue.NewQueue())
	queue.LogQueue.RunWatcher()
}

func main() {

	go producer()
	consumer()
	close()
}

func producer() {
	for {
		for i := 0; i < 10000; i++ {
			time.Sleep(500 * time.Millisecond)
			log.Log.Info("xxxxxxxxxx:" + strconv.Itoa(i))
		}
	}
}

func consumer() {
	maxWorkers := 1

	hosts := []string{"127.0.0.1:8000", "127.0.0.1:7000", "127.0.0.1:6000"}
	timeout := 10 * time.Second
	ttl := 500 * time.Millisecond

	clientArgs := output.NewClientArgs("tcp", hosts, timeout, ttl)

	dispatcherTTL := 1 * time.Second
	Dispatcher = output.NewDispatcher(maxWorkers, queue.LogQueue, clientArgs, dispatcherTTL)
	err := Dispatcher.Start()
	if err != nil {
		panic(err)
	}
}

func close() {
	time.Sleep(10000 * time.Second)

	Dispatcher.Close()

	fmt.Println(queue.LogQueue.Dump())

	queue.LogQueue.CloseWatcher()

}

func testlog() {
	log.Log.Info("xxxxxx")
}
