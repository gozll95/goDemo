package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	var svc StringService
	svc = stringService{}
	svc = proxyingMiddleware(context.Background(), *proxy, logger)(svc)
	svc = loggingMiddleware(logger)(svc)
	svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)
	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}

/*
The complete service so far is stringsvc3.

$ go get github.com/go-kit/kit/examples/stringsvc3
$ stringsvc3 -listen=:8001 &
listen=:8001 caller=proxying.go:25 proxy_to=none
listen=:8001 caller=main.go:72 msg=HTTP addr=:8001
$ stringsvc3 -listen=:8002 &
listen=:8002 caller=proxying.go:25 proxy_to=none
listen=:8002 caller=main.go:72 msg=HTTP addr=:8002
$ stringsvc3 -listen=:8003 &
listen=:8003 caller=proxying.go:25 proxy_to=none
listen=:8003 caller=main.go:72 msg=HTTP addr=:8003
$ stringsvc3 -listen=:8080 -proxy=localhost:8001,localhost:8002,localhost:8003
listen=:8080 caller=proxying.go:29 proxy_to="[localhost:8001 localhost:8002 localhost:8003]"
listen=:8080 caller=main.go:72 msg=HTTP addr=:8080
$ for s in foo bar baz ; do curl -d"{\"s\":\"$s\"}" localhost:8080/uppercase ; done
{"v":"FOO","err":null}
{"v":"BAR","err":null}
{"v":"BAZ","err":null}
listen=:8001 caller=logging.go:28 method=uppercase input=foo output=FOO err=null took=5.168µs
listen=:8080 caller=logging.go:28 method=uppercase input=foo output=FOO err=null took=4.39012ms
listen=:8002 caller=logging.go:28 method=uppercase input=bar output=BAR err=null took=5.445µs
listen=:8080 caller=logging.go:28 method=uppercase input=bar output=BAR err=null took=2.04831ms
listen=:8003 caller=logging.go:28 method=uppercase input=baz output=BAZ err=null took=3.285µs
listen=:8080 caller=logging.go:28 method=uppercase input=baz output=BAZ err=null took=1.388155ms
*/
