package manager

import "context"

type WokerFactory func(chan chan Job, context.Context) Worker

type Worker interface {
	Start()
	//Close()
}
