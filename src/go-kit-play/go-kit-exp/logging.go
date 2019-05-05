package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

<<<<<<< HEAD
func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next StringService) StringService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	StringService
}

func (mw logmw) Uppercase(s string) (output string, err error) {
=======
// 装饰endpoint
type loggingMiddleware struct {
	logger log.Logger
	next   StringService
}

func (mw loggingMiddleware) Uppercase(s string) (output string, err error) {
>>>>>>> 9ccf17f8216e7b960ff871bfc2dc378b13895252
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

<<<<<<< HEAD
	output, err = mw.StringService.Uppercase(s)
	return
}

func (mw logmw) Count(s string) (n int) {
=======
	output, err = mw.next.Uppercase(s)
	return
}

func (mw loggingMiddleware) Count(s string) (n int) {
>>>>>>> 9ccf17f8216e7b960ff871bfc2dc378b13895252
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

<<<<<<< HEAD
	n = mw.StringService.Count(s)
=======
	n = mw.next.Count(s)
>>>>>>> 9ccf17f8216e7b960ff871bfc2dc378b13895252
	return
}
