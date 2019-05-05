package main

import (
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// 客户端端点
// proxymw implements StringService, forwarding Uppercase requests to the
// provided endpoint, and serving all other (i.e. Count) requests via the
// next StringService.
type proxymw struct {
	next      StringService     // Serve most requests via this service...
	uppercase endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
}

func (mw proxymw) Uppercase(s string) (string, error) {
	response, err := mw.uppercase(uppercaseRequest{S: s})
	if err != nil {
		return "", err
	}
	resp := response.(uppercaseResponse)
	if resp.Err != "" {
		return resp.V, errors.New(resp.Err)
	}
	return resp.V, nil
}

type ServiceMiddleware func(next StringService) StringService

func proxyingMiddleware(proxyURL string) ServiceMiddleware {
	return func(next StringService) StringService {
		return proxymw{next, makeUppercaseEndpoint(proxyURL)}
	}
}

func makeUppercaseEndpoint(proxyURL string) endpoint.Endpoint {
	return httptransport.NewClient(
		"GET",
		mustParseURL(proxyURL),
		encodeUppercaseRequest,
		decodeUppercaseResponse,
	).Endpoint()
}

func main() {

	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(proxyingMiddleware("proxyURL")(svc)),
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
	http.ListenAndServe(":8080", nil)
}
