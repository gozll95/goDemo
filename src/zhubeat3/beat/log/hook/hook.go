package hook

import (
	"zhubeat/beat/job"
	"zhubeat/beat/queue"

	log "github.com/sirupsen/logrus"
)

type ZhuHook struct {
}

func NewZhuHook() *ZhuHook {
	return &ZhuHook{}
}

func (h *ZhuHook) Fire(e *log.Entry) error {

	msg, err := e.String()

	if err != nil {
		return err
	}
	go func(msg string) {
		// check queue is init?
		// init queue
		queue.LogQueue.Push(job.Job(msg))
	}(msg)

	return nil

}

func (h *ZhuHook) Levels() []log.Level {
	return log.AllLevels
}
