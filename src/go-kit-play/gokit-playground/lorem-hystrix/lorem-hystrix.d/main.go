package main

import (
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"fmt"
	"github.com/go-kit/kit/log"
	ratelimitkit "github.com/go-kit/kit/ratelimit"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"time"
	"github.com/juju/ratelimit"
	"flag"
	ilog "log"
	"github.com/ru-rocker/gokit-playground/lorem-hystrix"
)

func main() {

	// parse variable from input command
	var (
		consulAddr = flag.String("consul.addr", "", "consul address")
		consulPort = flag.String("consul.port", "", "consul port")
		advertiseAddr = flag.String("advertise.addr", "", "advertise address")
		advertisePort = flag.String("advertise.port", "", "advertise port")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	//declare metrics
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ru_rocker",
		Subsystem: "lorem_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ru_rocker",
		Subsystem: "lorem_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc lorem_hystrix.Service
	svc = lorem_hystrix.LoremService{}
	svc = lorem_hystrix.LoggingMiddleware(logger)(svc)
	svc = lorem_hystrix.Metrics(requestCount, requestLatency)(svc)

	rlbucket := ratelimit.NewBucket(35*time.Millisecond, 100)
	e := lorem_hystrix.MakeLoremLoggingEndpoint(svc)
	e = ratelimitkit.NewTokenBucketThrottler(rlbucket, time.Sleep)(e)

	// Make health endpoint
	healthEndpoint := lorem_hystrix.MakeHealthEndpoint(svc)
	endpoint := lorem_hystrix.Endpoints{
		LoremEndpoint: e,
		HealthEndpoint: healthEndpoint,
	}

	r := lorem_hystrix.MakeHttpHandler(ctx, endpoint, logger)

	// Register Service to Consul
	registar := lorem_hystrix.Register(*consulAddr,
		*consulPort,
		*advertiseAddr,
		*advertisePort)

	// HTTP transport
	go func() {
		ilog.Println("Starting server at port", *advertisePort)
		// register service
		registar.Register()
		handler := r
		errChan <- http.ListenAndServe(":" + *advertisePort, handler)
	}()


	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	error := <- errChan
	// deregister service
	registar.Deregister()
	ilog.Fatalln(error)
}
