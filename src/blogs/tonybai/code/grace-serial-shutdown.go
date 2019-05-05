package main

import (
	"errors"
	"sync"
	"time"
)

type GracefullyShutdowner interface {
	Shutdown(waitTimeout time.Duration) error
}
type ShutdownerFunc func(time.Duration) error

func (f ShutdownerFunc) Shutdown(waitTimeout time.Duration) error {
	return f(waitTimeout)
}

//shutdown.go
func SequentialShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
    start := time.Now()
    var left time.Duration

    for _, g := range shutdowners {
        elapsed := time.Since(start)
        left = waitTimeout - elapsed

        c := make(chan struct{})
        go func(shutdowner GracefullyShutdowner) {
            shutdowner.Shutdown(left)
            c <- struct{}{}
        }(g)

        select {
        case <-c:
            //continue
        case <-time.After(left):
            return errors.New("wait timeout")
        }
    }

    return nil
}