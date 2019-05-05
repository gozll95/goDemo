package watcher

import (
	"time"
)

type WatcherSource interface {
	IsOverThreshold() bool
	OverThresholdProccess() error
}

type Watcher struct {
	WatcherSource
	ticker    *time.Ticker
	threshold int
	quit      chan struct{}
}

func NewWatcher(ttl time.Duration, watcherSource WatcherSource) (w *Watcher) {
	w = &Watcher{
		WatcherSource: watcherSource,
		quit:          make(chan struct{}),
	}
	if ttl > 0 {
		w.ticker = time.NewTicker(ttl)
	}
	return w
}

func (w *Watcher) Start() {
	select {
	case <-w.quit:
		return
	default:
		go w.run()
	}
}

func (w *Watcher) run() {
	for {
		select {
		case <-w.quit:
			return
		case <-w.ticker.C:
			if w.IsOverThreshold() {
				err := w.OverThresholdProccess()
				if err != nil {
					panic(err)
				}
			}

		}
	}
}

func (w *Watcher) Close() {
	close(w.quit)
}
