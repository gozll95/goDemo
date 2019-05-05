package model

import (
	"time"

	"github.com/zhu/qvm/server/enums"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	Ticket *_Ticket

	ticketCollection = "ticket"
	ticketIndexes    = []mgo.Index{
		{
			Key:    []string{"uid", "region_id", "resource_type"},
			Unique: true,
		},
		{
			Key: []string{"status", "updated_at"},
		},
	}
)

const MaxBucketSize = 200

type TicketModel struct {
	Id           bson.ObjectId          `bson:"_id" json:"_id"`
	Uid          uint32                 `bson:"uid" json:"uid"`
	RegionId     string                 `bson:"region_id" json:"region_id"`
	ResourceType enums.ResourceType     `bson:"resource_type" json:"resource_type"`
	Status       enums.TicketStatusType `bson:"status" json:"status"`

	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	isNewRecord bool `bson:"-" json:"-"`
}

func NewTicketModel(uid uint32, regionId string, resourceType enums.ResourceType) *TicketModel {
	now := time.Now()

	return &TicketModel{
		Id:           bson.NewObjectId(),
		Uid:          uid,
		RegionId:     regionId,
		ResourceType: resourceType,
		Status:       enums.TicketStatusTypeIdle,
		CreatedAt:    now,
		UpdatedAt:    now,

		isNewRecord: true,
	}
}

func (ticket *TicketModel) Save() (err error) {
	if !ticket.Id.Valid() {
		return ErrInvalidId
	}

	Ticket.Query(func(c *mgo.Collection) {
		t := time.Now()
		if ticket.isNewRecord {
			ticket.CreatedAt = t
			ticket.UpdatedAt = t

			err = c.Insert(ticket)
			if err == nil {
				ticket.isNewRecord = false
			}
		} else {
			migration := bson.M{
				"$set": bson.M{
					"status":     ticket.Status,
					"updated_at": t,
				},
			}

			err = c.UpdateId(ticket.Id, migration)
		}
	})

	return
}

func (ticketModel *TicketModel) Unlock() (err error) {
	ticketModel.Status = enums.TicketStatusTypeIdle
	return ticketModel.Save()
}

func (ticket *TicketModel) Remove() (err error) {
	if !ticket.Id.Valid() {
		return ErrInvalidId
	}

	Ticket.Query(func(c *mgo.Collection) {
		err = c.RemoveId(ticket.Id)
	})

	return
}

type _Ticket struct{}

func (_ *_Ticket) Find(id string) (ticket *TicketModel, err error) {
	if !bson.IsObjectIdHex(id) {
		return nil, ErrInvalidId
	}

	Ticket.Query(func(c *mgo.Collection) {
		err = c.FindId(bson.ObjectIdHex(id)).One(&ticket)
	})

	return
}

func (_ *_Ticket) All() (tickets []*TicketModel, err error) {

	Ticket.Query(func(c *mgo.Collection) {
		err = c.Find(nil).All(&tickets)
	})
	return

}

func (_ *_Ticket) FindbyRegionIdAndResourcetype(uid uint32, regionId string, resourceType enums.ResourceType) (ticket *TicketModel, err error) {
	if uid == 0 || regionId == "" {
		return nil, ErrInvalidId
	}

	Ticket.Query(func(c *mgo.Collection) {
		query := bson.M{
			"uid":           uid,
			"region_id":     regionId,
			"resource_type": resourceType,
		}
		err = c.Find(query).One(&ticket)
	})

	return
}

// this will return a idle item, and modify it's status to running to lock this item
// if nothing found will return mgo.NotFound
func (_ *_Ticket) FindAndLock() (ticket *TicketModel, err error) {
	Ticket.Query(func(c *mgo.Collection) {
		query := bson.M{
			"status": enums.TicketStatusTypeIdle,
			"updated_at": bson.M{
				"$lte": time.Now().Add(time.Second * 3 * -1), // updated_at lock item in near three seconds
			},
		}

		change := mgo.Change{
			Update: bson.M{
				"$set": bson.M{
					"status":     enums.TicketStatusTypeBusy,
					"updated_at": time.Now(),
				},
			},
			ReturnNew: true,
			Upsert:    false,
		}

		_, err = c.Find(query).Apply(change, &ticket)
	})

	return
}

func (_ *_Ticket) Query(query func(c *mgo.Collection)) {
	MongoModel().Query(ticketCollection, ticketIndexes, query)
}
