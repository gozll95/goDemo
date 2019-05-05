package main

import "github.com/go-kit/kit/endpoint"

type Subscriber interface {
	Endpoints() ([]endpoint.Endpoint, error)
}

//type Factory func(instance string) (endpoint.Endpoint, error)
//instance=host:port




var e endpoint.Endpoint

e = makeUppercaseProxy(instance)

// add  circuitbreaker.Gobreaker
e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
// add  kitratelimit.NewTokenBucketLimiter
e = kitratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(float64(maxQPS), int64(maxQPS)))(e)
}


// 调用other服务,并且支持服务发现+负载均衡+重试机制的
func proxyingMiddleware(instances string, logger log.Logger) ServiceMiddleware {
	// If instances is empty, don't proxy.
	if instances == "" {
		logger.Log("proxy_to", "none")
		return func(next StringService) StringService { return next }
	}

	// Set some parameters for our client.
	var (
		qps         = 100                    // beyond which we will return an error
		maxAttempts = 3                      // per request, before giving up
		maxTime     = 250 * time.Millisecond // wallclock time, before giving up
	)

	// Otherwise, construct an endpoint for each instance in the list, and add
	// it to a fixed set of endpoints. In a real service, rather than doing this
	// by hand, you'd probably use package sd's support for your service
	// discovery system.
	var (
		instanceList = split(instances)
		subscriber   sd.FixedSubscriber
	)
	logger.Log("proxy_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeUppercaseProxy(instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = kitratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(float64(qps), int64(qps)))(e)
		subscriber = append(subscriber, e)
	}

	// Now, build a single, retrying, load-balancing endpoint out of all of
	// those individual endpoints.
	balancer := lb.NewRoundRobin(subscriber)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	// And finally, return the ServiceMiddleware, implemented by proxymw.
	return func(next StringService) StringService {
		return proxymw{next, retry}
	}
}

/*
1.
Subscriber 适配多种服务发现机制
type Subscriber interface {
	Endpoints() ([]endpoint.Endpoint, error)
}

内部实现
type Factory func(instance string) (endpoint.Endpoint, error)

2.
负载均衡[wrap Subscriber] 从一组节点中选择一个
type Balancer interface {
	Endpoint() (endpoint.Endpoint, error)
}


3.重试
Retry[wrap(Balancer(wrap(Subscriber)))] 
func Retry(max int, timeout time.Duration, lb Balancer) endpoint.Endpoint



*/