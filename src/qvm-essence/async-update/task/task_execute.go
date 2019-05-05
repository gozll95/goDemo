package task

import (
	"time"

	"github.com/zhu/qvm/server/model"
	"github.com/zhu/qvm/server/utils"
)

// true: no ticket left
func (t *TaskManager) execute() bool {

	ticket, err := model.Ticket.FindAndLock()
	if err == model.ErrNotFound {
		return true
	}

	if err != nil {
		utils.StdLog.Errorf("task manager execute error %v", err)
		return false
	}

	tasks, err := model.Task.AllByUidRegionAndResourceType(ticket.Uid, ticket.RegionId, ticket.ResourceType)
	if err != nil {
		utils.StdLog.Errorf("task manager execute error %v", err)
		return false
	}

	defer func() {
		if len(tasks) != 0 {
			// free ticket
			err := ticket.Unlock()
			if err != nil {
				utils.StdLog.Errorf("task manager unlock ticket %s error %v", ticket.Id, err)
			}
		}
	}()

	// ticket has not tasks remove this ticket
	if len(tasks) == 0 {
		err = ticket.Remove()
		if err != nil {
			utils.StdLog.Errorf("task manager remove ticket %s error %v", ticket.Id, err)
		}

		return false
	}

	// merge tasks to one request
	resourceIds := make([]string, len(tasks))
	for i, task := range tasks {
		resourceIds[i] = task.ResourceId
	}

	checkers, err := t.do(ticket.Uid, ticket.RegionId, ticket.ResourceType, resourceIds)
	if err != nil {
		utils.StdLog.Errorf("task manager do task error %v", err)
		return false
	}

	checkerMap := make(map[string]Checker)
	for _, checker := range checkers {
		checkerMap[checker.ID()] = checker
	}

	now := time.Now()
	for _, task := range tasks {
		checker, ok := checkerMap[task.ResourceId]
		if !ok {
			// find task not in checkers
			// then remove this task
			// Note: this most happen when user has a resource task running and user delete this resource
			err = task.Remove()
			if err != nil && err != model.ErrNotFound {
				utils.StdLog.Errorf("task manager remove task %s error %v", task.Id, err)
			}
			continue
		}

		// updated message
		if checker.Updated() {
			msg, err := model.ResourceMessage.FindOrCreate(task.Uid, task.RegionId, task.ResourceType, task.ResourceId)
			if err != nil {
				utils.StdLog.Errorf("task manager find or create resource message error %v", err)
				continue
			}

			err = msg.Save()
			if err != nil {
				utils.StdLog.Errorf("task manager save resource message error %v", err)
				continue

			}
		}

		// remove if it is an once task
		if task.OnlyOnce {
			err = task.Remove()
		}

		// check is target status
		if task.TargetStatus != "" && task.TargetStatus == checker.ResourceStatus() {
			err = task.Remove()
		}

		// resource is in stable status
		if task.TargetStatus == "" && IsStableStatus(task.ResourceType, checker.ResourceStatus()) {
			err = task.Remove()

		}

		// if task is time out then remove
		if now.After(task.UpdatedAt.Add(time.Second * time.Duration(task.Timeout))) {
			err = task.Remove()
		}

		if err != nil && err != model.ErrNotFound {
			utils.StdLog.Errorf("task manager remove task %s error %v", task.Id, err)
		}
	}

	return false
}
