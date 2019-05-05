package status

import (
	"sync"
)

type StatusScheduler struct {
	status     Status
	statusLock sync.RWMutex
}

func NewStatusScheduler() *StatusScheduler {
	return &StatusScheduler{}
}

// 用于状态的检查，并在条件满足时设置状态。
func (sched *StatusScheduler) CheckAndSetStatus(wantedStatus Status) (oldStatus Status, err error) {
	sched.statusLock.Lock()
	defer sched.statusLock.Unlock()
	oldStatus = sched.status
	err = sched.status.CheckStatusFactory(wantedStatus, nil)
	if err == nil {
		sched.status = wantedStatus
	}
	return
}

func (sched *StatusScheduler) CheckAndSetStatusWithErr(oldStatus, wantedStatus Status, err error) {
	if err != nil {
		sched.status = oldStatus
		return
	}
	sched.statusLock.Lock()
	defer sched.statusLock.Unlock()
	sched.status = wantedStatus
}

func (sched *StatusScheduler) Status() Status {
	sched.statusLock.Lock()
	defer sched.statusLock.Unlock()
	return sched.status
}
