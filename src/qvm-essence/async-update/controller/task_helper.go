package controller

import (
	"github.com/zhu/qvm/server/model"
	"github.com/zhu/qvm/server/utils"
)

var (
	TaskHelper *_TaskHelper
)

type _TaskHelper struct{}

func (_ *_TaskHelper) NewInstanceTasks(logger utils.Log, tasks ...*model.InstanceModel) {
	for _, tasker := range tasks {
		err := TaskManager.NewTaskByTasker(tasker)

		if err != nil {
			logger.Errorf("New Task  %#v error %v", tasker, err)
		}
	}
}

func (_ *_TaskHelper) NewDiskTasks(logger utils.Log, tasks ...*model.DiskModel) {
	for _, tasker := range tasks {
		err := TaskManager.NewTaskByTasker(tasker)

		if err != nil {
			logger.Errorf("New Task  %#v error %v", tasker, err)
		}
	}
}

func (_ *_TaskHelper) NewEipTasks(logger utils.Log, tasks ...*model.EipModel) {
	for _, tasker := range tasks {
		err := TaskManager.NewTaskByTasker(tasker)

		if err != nil {
			logger.Errorf("New Task  %#v error %v", tasker, err)
		}
	}
}

func (_ *_TaskHelper) NewSnapshotTasks(logger utils.Log, tasks ...*model.SnapshotModel) {
	for _, tasker := range tasks {
		err := TaskManager.NewTaskByTasker(tasker)

		if err != nil {
			logger.Errorf("New Task  %#v error %v", tasker, err)
		}
	}
}

func (_ *_TaskHelper) NewVSwitchTasks(logger utils.Log, tasks ...*model.VSwitchModel) {
	for _, tasker := range tasks {
		err := TaskManager.NewTaskByTasker(tasker)

		if err != nil {
			logger.Errorf("New Task  %#v error %v", tasker, err)
		}
	}
}

func (_ *_TaskHelper) NewVpcTasks(logger utils.Log, tasks ...*model.VpcModel) {
	for _, tasker := range tasks {
		err := TaskManager.NewTaskByTasker(tasker)

		if err != nil {
			logger.Errorf("New Task  %#v error %v", tasker, err)
		}
	}
}

func (_ *_TaskHelper) NewImageTasks(logger utils.Log, tasks ...*model.ImageModel) {
	for _, tasker := range tasks {
		err := TaskManager.NewTaskByTasker(tasker)

		if err != nil {
			logger.Errorf("New Task  %#v error %v", tasker, err)
		}
	}
}
