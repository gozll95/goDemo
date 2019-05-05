package status

import (
	"errors"
	"fmt"
	"sync"
)

type Status uint32

const (
	// SCHED_STATUS_UNINITIALIZED 代表未初始化的状态。
	SCHED_STATUS_UNINITIALIZED Status = 0
	// SCHED_STATUS_INITIALIZING 代表正在初始化的状态。
	SCHED_STATUS_INITIALIZING Status = 1
	// SCHED_STATUS_INITIALIZED 代表已初始化的状态。
	SCHED_STATUS_INITIALIZED Status = 2
	// SCHED_STATUS_STARTING 代表正在启动的状态。
	SCHED_STATUS_STARTING Status = 3
	// SCHED_STATUS_STARTED 代表已启动的状态。
	SCHED_STATUS_STARTED Status = 4
	// SCHED_STATUS_STOPPING 代表正在停止的状态。
	SCHED_STATUS_STOPPING Status = 5
	// SCHED_STATUS_STOPPED 代表已停止的状态。
	SCHED_STATUS_STOPPED Status = 6
)

// GetStatusDescription 用于获取状态的文字描述。
func (s Status) string() string {
	switch s {
	case SCHED_STATUS_UNINITIALIZED:
		return "uninitialized"
	case SCHED_STATUS_INITIALIZING:
		return "initializing"
	case SCHED_STATUS_INITIALIZED:
		return "initialized"
	case SCHED_STATUS_STARTING:
		return "starting"
	case SCHED_STATUS_STARTED:
		return "started"
	case SCHED_STATUS_STOPPING:
		return "stopping"
	case SCHED_STATUS_STOPPED:
		return "stopped"
	default:
		return "unkown"
	}
}

func (s Status) CheckStatusFactory(wantedStatus Status, lock sync.Locker) (err error) {
	var currentStatus Status
	currentStatus = s
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	switch currentStatus {
	case SCHED_STATUS_INITIALIZING:
		return errors.New("the scheduler is being initialized!")
	case SCHED_STATUS_STARTING:
		return errors.New("the scheduler is being started!")
	case SCHED_STATUS_STOPPING:
		return errors.New("the scheduler is being stopped!")
	}
	if err != nil {
		return
	}
	if currentStatus == SCHED_STATUS_UNINITIALIZED &&
		(wantedStatus == SCHED_STATUS_STARTING ||
			wantedStatus == SCHED_STATUS_STOPPING) {
		err = errors.New("the scheduler has not yet been initialized!")
		return
	}
	switch wantedStatus {
	case SCHED_STATUS_INITIALIZING:
		switch currentStatus {
		case SCHED_STATUS_STARTED:
			err = errors.New("the scheduler has been started!")
		}
	case SCHED_STATUS_STARTING:
		switch currentStatus {
		case SCHED_STATUS_UNINITIALIZED:
			err = errors.New("the scheduler has not been initialized!")
		case SCHED_STATUS_STARTED:
			err = errors.New("the scheduler has been started!")
		}
	case SCHED_STATUS_STOPPING:
		if currentStatus != SCHED_STATUS_STARTED {
			err = errors.New("the scheduler has not been started!")
		}
	default:
		errMsg :=
			fmt.Sprintf("unsupported wanted status for check! (wantedStatus: %d)",
				wantedStatus)
		err = errors.New(errMsg)
	}
	return
}
