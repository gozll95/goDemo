package task

import (
	"sync"
	"time"

	"github.com/zhu/qvm/server/enums"
	"github.com/zhu/qvm/server/model"
	"github.com/zhu/qvm/server/utils"
)

const closed = "closed"

type TaskManager struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Worker chan bool

	Closed chan bool

	NoLeftLock sync.RWMutex
	IsNoLeft   bool
}

func NewTaskManager(worker int) *TaskManager {
	return &TaskManager{
		Worker: make(chan bool, worker),
		Closed: make(chan bool),
	}
}

func (t *TaskManager) NewTask(uid uint32, regionId string, resourceType enums.ResourceType, resourceId string) (err error) {
	return t.NewStatusTask(uid, regionId, resourceType, resourceId, "")

}

func (t *TaskManager) NewStatusTask(uid uint32, regionId string, resourceType enums.ResourceType, resourceId, targetStatus string) (err error) {
	task, err := model.Task.FindOrCreate(uid, regionId, resourceType, resourceId)
	if err != nil {
		return
	}

	task.SetTargetStatus(targetStatus)
	err = task.Save()

	if err != nil {
		return
	}

	_, err = model.Ticket.FindbyRegionIdAndResourcetype(uid, regionId, resourceType)
	if err == model.ErrNotFound {
		// Note: only no ticket save the ticket
		// or update ticket updatedAt time will effect ticket wait time
		ticket := model.NewTicketModel(uid, regionId, resourceType)
		return ticket.Save()
	}

	return err

}

func (t *TaskManager) NewTaskByTasker(tasker Tasker) (err error) {
	return t.NewTask(
		tasker.TaskUid(),
		tasker.TaskRegion(),
		tasker.TaskResourceType(),
		tasker.TaskResourceId(),
	)
}

func (t *TaskManager) NewTaskByTaskers(taskers ...Tasker) (err error) {
	for _, tasker := range taskers {
		err = t.NewTaskByTasker(tasker)
		if err != nil {
			return
		}
	}

	return
}

// create task which only will be excuted once
func (t *TaskManager) NewOnceTask(uid uint32, regionId string, resourceType enums.ResourceType, resourceId string) (err error) {
	task, err := model.Task.FindOrCreate(uid, regionId, resourceType, resourceId)
	if err != nil {
		return
	}

	task.OnlyOnce = true

	err = task.Save()
	if err != nil {
		return
	}

	_, err = model.Ticket.FindbyRegionIdAndResourcetype(uid, regionId, resourceType)
	if err == model.ErrNotFound {
		// Note: only no ticket save the ticket
		// or update ticket updatedAt time will effect ticket wait time
		ticket := model.NewTicketModel(uid, regionId, resourceType)
		return ticket.Save()
	}

	return err
}

func (t *TaskManager) Start() {

	go t.run()
	go t.refreshTicket()

}

func (t *TaskManager) Close() {
	defer func() {
		tickets, _ := model.Ticket.All()
		for _, ticket := range tickets {
			err := ticket.Unlock()
			if err != nil {
				utils.StdLog.Errorf("task manager unlock ticket %s error %v", ticket.Id, err)
			}
		}
		utils.StdLog.Info("all unlock")
	}()
	utils.StdLog.Debugf("accept Closed")
	close(t.Closed)
}

// refresh ticket
func (t *TaskManager) refreshTicket() {
	defer func() {
		if err := recover(); err != nil {
			if err == closed {
				utils.StdLog.Infof("task refreshTicket received closed signal")
			} else {
				t.refreshTicket()
			}
			return
		}
		t.refreshTicket()
	}()

	for {
		select {
		case <-t.Closed:
			panic(closed)
		default:
			// sleep half of ticket refresh time to wait task being older
			time.Sleep(time.Second * 9)

			iter := model.Task.FindUpdatedIter()
			task, isError := model.Task.Next(iter)

			for isError {
				_, err := model.Ticket.FindbyRegionIdAndResourcetype(task.Uid, task.RegionId, task.ResourceType)
				if err == model.ErrNotFound {
					// Note: only no ticket save the ticket
					// or update ticket updatedAt time will effect ticket wait time
					ticket := model.NewTicketModel(task.Uid, task.RegionId, task.ResourceType)
					err = ticket.Save()
				}

				if err != nil {
					utils.StdLog.Errorf("refresh ticket find or save error %v", err)
				}

				task, isError = model.Task.Next(iter)
			}
		}

	}
}

func (t *TaskManager) run() {
	defer func() {
		if err := recover(); err != nil {
			if err == closed {
				utils.StdLog.Infof("task run received closed signal")
			} else {
				t.run()
			}
			return
		}
		t.run()
	}()

	for {
		select {
		case <-t.Closed:
			panic(closed)
		case t.Worker <- true:
			go func() {
				// check No Left true
				t.NoLeftLock.RLock()

				// if no left task sleep 3 second and release no left
				if t.IsNoLeft {
					t.NoLeftLock.RUnlock()

					time.Sleep(time.Second * 3)

					t.NoLeftLock.Lock()
					t.IsNoLeft = false
					t.NoLeftLock.Unlock()
				} else {
					t.NoLeftLock.RUnlock()
				}

				isNoLeft := t.execute()

				if isNoLeft {
					// set no left marker
					t.NoLeftLock.Lock()
					t.IsNoLeft = true
					t.NoLeftLock.Unlock()
				}

				<-t.Worker
			}()

		}

	}
}
