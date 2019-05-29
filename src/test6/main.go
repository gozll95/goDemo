package main

import (
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func main() {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: "your_service_name",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		panic(err)
	}
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// continue main()
	a()

}

var (
	n = time.Now()
)

func a() {
	tracer := opentracing.GlobalTracer()

	// mock rpc 0
	clientSpan := tracer.StartSpan("rpc 0", opentracing.StartTime(n))
	a := opentracing.FinishOptions{
		FinishTime: n.Add(5 * time.Hour),
	}
	defer clientSpan.FinishWithOptions(a)

	// Set some tags on the clientSpan to annotate that it's the client span. The additional HTTP tags are useful for debugging purposes.
	ext.SpanKindRPCClient.Set(clientSpan)
	ext.HTTPUrl.Set(clientSpan, "www.baidu.com")
	ext.HTTPMethod.Set(clientSpan, "GET")

	// mock rpc 0.1
	aApp := opentracing.StartSpan("rpc 0.1", opentracing.ChildOf(clientSpan.Context()), opentracing.StartTime(n.Add(time.Hour))) // Start a new root span.
	a = opentracing.FinishOptions{
		FinishTime: n.Add(2 * time.Hour),
	}
	defer aApp.FinishWithOptions(a)

	// mock rpc 0.2
	bApp := opentracing.StartSpan("rpc 0.2", opentracing.ChildOf(clientSpan.Context()), opentracing.StartTime(n.Add(time.Hour))) // Start a new root span.
	a = opentracing.FinishOptions{
		FinishTime: n.Add(3 * time.Hour),
	}
	defer bApp.FinishWithOptions(a)

}
