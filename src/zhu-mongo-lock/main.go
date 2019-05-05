package model

import (
	"fmt"
	"time"

	"github.com/zhu/qvm/server/enums"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	MnsSend *_Offset

	offsetCollection = "mns_send"
	taskIndexes    = []mgo.Index{
		{
			Key:    []string{"uid", "region_id", "resource_type", "resource_id"},
			Unique: true,
		},
		{
			Key: []string{"updated_at"},
		},
	}
)

type TaskModel struct {
	Id           bson.ObjectId      `bson:"_id" json:"_id"`
	Uid          uint32             `bson:"uid" json:"uid"`
	RegionId     string             `bson:"region_id" json:"region_id"`
	ResourceType enums.ResourceType `bson:"resource_type" json:"resource_type"`
	ResourceId   string             `bson:"resource_id" json:"resource_id"`
	TargetStatus string             `bson:"target_status" json:"target_status"`

	Timeout   int       `bson:"timeout" json:"timeout"` // default unit second
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	isNewRecord bool `bson:"-" json:"-"`
}

func NewTaskModel(uid uint32, regionId string, resourceType enums.ResourceType, resourceId string) *TaskModel {
	now := time.Now()

	return &TaskModel{
		Id:           bson.NewObjectId(),
		Uid:          uid,
		RegionId:     regionId,
		ResourceType: resourceType,
		ResourceId:   resourceId,
		Timeout:      60 * 5, // default is 5 minute
		CreatedAt:    now,
		UpdatedAt:    now,

		isNewRecord: true,
	}
}

func (task *TaskModel) SetTargetStatus(status string) {
	task.TargetStatus = status
}

func (task *TaskModel) SetTimeoutAt(timeout int) {
	if timeout > 0 {
		task.Timeout = timeout
	}
}

func (task *TaskModel) Save() (err error) {
	if !task.Id.Valid() {
		return ErrInvalidId
	}

	Task.Query(func(c *mgo.Collection) {
		t := time.Now()
		if task.isNewRecord {
			task.CreatedAt = t
			task.UpdatedAt = t

			err = c.Insert(task)
			if err == nil {
				task.isNewRecord = false
			}
		} else {
			migration := bson.M{
				"$set": bson.M{
					"target_status": task.TargetStatus,
					"timeout":       task.Timeout,
					"updated_at":    t,
				},
			}

			err = c.UpdateId(task.Id, migration)
		}
	})

	return
}

func (task *TaskModel) Remove() (err error) {
	if !task.Id.Valid() {
		return ErrInvalidId
	}

	Task.Query(func(c *mgo.Collection) {

		err = c.RemoveId(task.Id)
	})

	return
}

type _Task struct{}

func (_ *_Task) Find(id string) (task *TaskModel, err error) {
	if !bson.IsObjectIdHex(id) {
		return nil, ErrInvalidId
	}

	Task.Query(func(c *mgo.Collection) {
		err = c.FindId(bson.ObjectIdHex(id)).One(&task)
	})

	return
}

func (_ *_Task) FindUpdatedIter() (iter *mgo.Iter) {
	Task.Query(func(c *mgo.Collection) {
		query := bson.M{
			"updated_at": bson.M{
				"$lt": time.Now().Add(time.Second * -18), // 5 * task
			},
		}

		iter = c.Find(query).Iter()
	})

	return
}

func (_ *_Task) Next(iter *mgo.Iter) (task *TaskModel, isError bool) {
	if iter == nil {
		fmt.Println("............")
		return nil, false
	}

	isError = iter.Next(&task)

	return
}

func (_ *_Task) FindOrCreate(uid uint32, regionId string, resourceType enums.ResourceType, resourceId string) (task *TaskModel, err error) {
	task, err = Task.FindByResourceId(uid, regionId, resourceType, resourceId)
	if err == ErrNotFound {
		err = nil
		task = NewTaskModel(uid, regionId, resourceType, resourceId)
		return
	}

	return
}

func (_ *_Task) FindByResourceId(uid uint32, regionId string, resourceType enums.ResourceType, resourceId string) (task *TaskModel, err error) {
	if resourceId == "" {
		return nil, ErrInvalidId
	}

	Task.Query(func(c *mgo.Collection) {
		query := bson.M{
			"uid":           uid,
			"region_id":     regionId,
			"resource_type": resourceType,
			"resource_id":   resourceId,
		}

		err = c.Find(query).One(&task)
	})

	return
}

func (_ *_Task) AllByUidRegionAndResourceType(uid uint32, regionId string, resourceType enums.ResourceType) (tasks []*TaskModel, err error) {
	if uid == 0 {
		return nil, ErrInvalidParams
	}

	Task.Query(func(c *mgo.Collection) {
		query := bson.M{
			"uid":           uid,
			"region_id":     regionId,
			"resource_type": resourceType,
		}

		err = c.Find(query).Sort("-_id").Limit(100).All(&tasks)
	})

	return
}

func (_ *_Task) Query(query func(c *mgo.Collection)) {
	MongoModel().Query(taskCollection, taskIndexes, query)
}
