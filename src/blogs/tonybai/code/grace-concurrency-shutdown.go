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

// shutdown.go
func ConcurrencyShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	c := make(chan struct{})

	go func() {
		var wg sync.WaitGroup
		for _, g := range shutdowners {
			wg.Add(1)
			go func(shutdowner GracefullyShutdowner) {
				shutdowner.Shutdown(waitTimeout)
				wg.Done()
			}(g)
		}
		wg.Wait()
		c <- struct{}{}
	}()

	select {
	case <-c:
		return nil
	case <-time.After(waitTimeout):
		return errors.New("wait timeout")
	}
}
