package main

import (
	"time"

	"github.com/opentracing/opentracing-go"
)

var (
	n = time.Now()
)

func createSpan() {
	// MOCK 0
	root := opentracing.StartSpan("GET /0", opentracing.StartTime(n)) // Start a new root span.
	a := opentracing.FinishOptions{
		FinishTime: n.Add(time.Hour),
	}
	defer root.FinishWithOptions(a)

	// MOCK 0.1
	aApp := opentracing.StartSpan("GET /0.1", opentracing.ChildOf(root.Context()), opentracing.StartTime(n.Add(2*time.Hour))) // Start a new root span.
	a = opentracing.FinishOptions{
		FinishTime: n.Add(3 * time.Hour),
	}
	defer aApp.FinishWithOptions(a)

}
