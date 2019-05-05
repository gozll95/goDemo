package main

/*
import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Middleware func(endpoint.Endpoint) endpoint.Endpoint

func loggingMiddleware(logger log.Logger) Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(_ context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", "calling endpoint")
			defer logger.Log("msg", "called endpoint")
			return next(nil, request)
		}
	}
}

// Transports expose the service to the network. In this first example we utilize JSON over HTTP.
func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	svc := stringService{}

	var uppercase endpoint.Endpoint
	uppercase = makeUppercaseEndpoint(svc)
	uppercase = loggingMiddleware(log.With(logger, "method", "uppercase"))(uppercase)

	var count endpoint.Endpoint
	count = makeCountEndpoint(svc)
	count = loggingMiddleware(log.With(logger, "method", "count"))(count)

	uppercaseHandler := httptransport.NewServer(
		uppercase,
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		count,
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)

	http.ListenAndServe(":8080", nil)
}
*/

/* request
$ curl -XPOST -d'{"s":"hello, world"}' localhost:8080/uppercase
{"v":"HELLO, WORLD","err":null}
$ curl -XPOST -d'{"s":"hello, world"}' localhost:8080/count
{"v":12}
*/

/* response:
method=uppercase msg="calling endpoint"
method=uppercase msg="called endpoint"
method=count msg="calling endpoint"
method=count msg="called endpoint"
*/

// 使用端点中间件来解决传输域问题，例如断路和速率限制。使用服务中间件来解决业务领域问题，例如日志记录和检测。说到仪器......
