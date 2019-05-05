package model

import (
	"math/rand"
	"testing"

	"github.com/zhu/qvm/server/enums"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Task_Find(t *testing.T) {
	task := NewTaskModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeInstance,
		uuid.NewV4().String(),
	)

	err := task.Save()
	assert.Nil(t, err)

	// duplicate save
	err = task.Save()
	assert.Nil(t, err)

	// find
	newTask, err := Task.Find(task.Id.Hex())
	assert.Nil(t, err)
	assert.Equal(t, task.Uid, newTask.Uid)
	assert.Equal(t, task.ResourceType, newTask.ResourceType)
	assert.Equal(t, task.ResourceId, newTask.ResourceId)
	assert.Equal(t, 5*60, task.Timeout)
}

func Test_Task_FindNext(t *testing.T) {
	task := NewTaskModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeInstance,
		uuid.NewV4().String(),
	)

	err := task.Save()
	assert.Nil(t, err)

	//
	//time.Sleep(time.Second * 16)
	iter := Task.FindUpdatedIter()
	newTask, isError := Task.Next(iter)
	assert.Nil(t, newTask)
	assert.False(t, isError)
	// find
	//assert.True(t, isError)
	//assert.Equal(t, task.Uid, newTask.Uid)
	//assert.Equal(t, task.ResourceType, newTask.ResourceType)
	//assert.Equal(t, task.ResourceId, newTask.ResourceId)
	//assert.Equal(t, 5*60, task.Timeout)
	//
	//newTask, isError = Task.Next(iter)
	//assert.False(t, isError)
	//assert.Nil(t, newTask)
}

func Test_Task_FindOrCreate(t *testing.T) {
	task := NewTaskModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeInstance,
		uuid.NewV4().String(),
	)

	err := task.Save()
	assert.Nil(t, err)

	// find
	newTask, err := Task.FindOrCreate(task.Uid, task.RegionId, task.ResourceType, task.ResourceId)
	assert.Nil(t, err)
	assert.Equal(t, task.Id, newTask.Id)
	assert.False(t, task.isNewRecord)

	// create
	newTask1, err := Task.FindOrCreate(task.Uid, task.RegionId, task.ResourceType, uuid.NewV4().String())
	assert.Nil(t, err)
	assert.True(t, newTask1.isNewRecord)
}

func Test_Task_FindByResourceId(t *testing.T) {
	task := NewTaskModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeInstance,
		uuid.NewV4().String(),
	)

	err := task.Save()
	assert.Nil(t, err)

	// duplicate save
	err = task.Save()
	assert.Nil(t, err)

	// find
	newTask, err := Task.FindByResourceId(task.Uid, task.RegionId, task.ResourceType, task.ResourceId)
	assert.Nil(t, err)
	assert.Equal(t, task.Uid, newTask.Uid)
	assert.Equal(t, task.ResourceType, newTask.ResourceType)
	assert.Equal(t, task.ResourceId, newTask.ResourceId)
	assert.Equal(t, 5*60, task.Timeout)
}

func Test_Task_Remove(t *testing.T) {
	task := NewTaskModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeInstance,
		uuid.NewV4().String(),
	)

	task.Save()
	err := task.Remove()
	assert.Nil(t, err)

	// find
	newTask, err := Task.Find(task.Id.Hex())
	assert.Nil(t, newTask)
	assert.Equal(t, ErrNotFound, err)
}

func Test_Task_AllByUidAndResourceType(t *testing.T) {
	var (
		uid      = rand.Uint32()
		regionId = uuid.NewV4().String()
	)

	task1 := NewTaskModel(
		uid,
		regionId,
		enums.ResourceTypeInstance,
		uuid.NewV4().String(),
	)
	task2 := NewTaskModel(
		uid,
		regionId,
		enums.ResourceTypeIp,
		uuid.NewV4().String(),
	)
	task3 := NewTaskModel(
		uid,
		regionId,
		enums.ResourceTypeInstance,
		uuid.NewV4().String(),
	)
	task4 := NewTaskModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeInstance,
		uuid.NewV4().String(),
	)
	task5 := NewTaskModel(
		uid,
		uuid.NewV4().String(),
		enums.ResourceTypeInstance,
		uuid.NewV4().String(),
	)

	task1.Save()
	task2.Save()
	task3.Save()
	task4.Save()
	task5.Save()

	// find
	tasks, err := Task.AllByUidRegionAndResourceType(uid, regionId, enums.ResourceTypeInstance)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, task1.Id, tasks[1].Id)
	assert.Equal(t, task3.Id, tasks[0].Id)
}
