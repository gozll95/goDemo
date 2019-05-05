package model

import (
	"time"

	"github.com/zhu/qvm/server/enums"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	RdsLocker *_RdsLocker

	RdsLockerCollection = "rds_locker"
	RdsLockerIndexes    = []mgo.Index{
		{
			Key:    []string{"uid", "region_id", "db_instance_id"},
			Unique: true,
		},
	}

	RdsLockerExpire, _ = time.ParseDuration("-10m")
)

type RdsLockerModel struct {
	Id  bson.ObjectId `bson:"_id" json:"id"`
	Uid uint32        `bson:"uid" json:"-"`

	RegionId     string           `bson:"region_id" json:"region_id"`
	DBInstanceId string           `bson:"db_instance_id" json:"db_instance_id"`
	LockStatus   enums.LockStatus `bson:"lock_status" json:"lock_status"`

	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	isNewRecord bool `bson:"-" json:"-"`
	IsUpdated   bool `bson:"-" json:"-"`
}

// new rdslocker default status is unlock
func NewRdsLockerModel(uid uint32, regionId string, dbInstanceId string) *RdsLockerModel {
	now := time.Now()

	return &RdsLockerModel{
		Id:           bson.NewObjectId(),
		Uid:          uid,
		RegionId:     regionId,
		DBInstanceId: dbInstanceId,
		LockStatus:   enums.Unlock,
		CreatedAt:    now,
		UpdatedAt:    now,

		isNewRecord: true,
	}
}

func (rdsLocker *RdsLockerModel) Save() (err error) {
	if !rdsLocker.Id.Valid() {
		return ErrInvalidId
	}

	RdsLocker.Query(func(c *mgo.Collection) {
		t := time.Now()
		if rdsLocker.isNewRecord {
			rdsLocker.CreatedAt = t
			rdsLocker.UpdatedAt = t

			err = c.Insert(rdsLocker)
			if err == nil {
				rdsLocker.isNewRecord = false
			}
		} else {
			migration := bson.M{
				"$set": bson.M{
					"lock_status": rdsLocker.LockStatus,
					"updated_at":  t,
				},
			}

			err = c.UpdateId(rdsLocker.Id, migration)
		}
	})

	return
}

// find and change lock status
// locked -> unlock
// unlock -> locked
func (rdsLocker *RdsLockerModel) FindAndChangeLockStatus(wantedStatus enums.LockStatus) (err error) {
	uid := rdsLocker.Uid
	regionId := rdsLocker.RegionId
	dBInstanceId := rdsLocker.DBInstanceId

	if uid == 0 || regionId == "" || dBInstanceId == "" || !wantedStatus.IsValid() {
		return ErrInvalidParams
	}
	RdsLocker.Query(func(c *mgo.Collection) {
		query := bson.M{
			"uid":            uid,
			"region_id":      regionId,
			"db_instance_id": dBInstanceId,
			"lock_status":    wantedStatus.Opposite(),
		}

		change := mgo.Change{
			Update: bson.M{
				"$set": bson.M{
					"lock_status": wantedStatus,
					"updated_at":  time.Now(),
				},
			},
			ReturnNew: true,
			Upsert:    false,
		}

		_, err = c.Find(query).Apply(change, &rdsLocker)
	})

	return
}

// locked expired
func (rdsLocker *RdsLockerModel) SetUnlockWhenExpiredLocked(expireTime time.Duration) (err error) {
	uid := rdsLocker.Uid
	regionId := rdsLocker.RegionId
	dBInstanceId := rdsLocker.DBInstanceId

	if uid == 0 || regionId == "" || dBInstanceId == "" {
		return ErrInvalidParams
	}
	if expireTime >= time.Duration(int64(0)) {
		expireTime = RdsLockerExpire
	}
	RdsLocker.Query(func(c *mgo.Collection) {
		query := bson.M{
			"uid":            uid,
			"region_id":      regionId,
			"db_instance_id": dBInstanceId,
			"lock_status":    enums.Locked,
			"updated_at": bson.M{
				"$lte": time.Now().Add(expireTime),
			},
		}

		change := mgo.Change{
			Update: bson.M{
				"$set": bson.M{
					"lock_status": enums.Unlock,
					"updated_at":  time.Now(),
				},
			},
			ReturnNew: true,
			Upsert:    false,
		}

		_, err = c.Find(query).Apply(change, &rdsLocker)
	})

	return
}

type _RdsLocker struct{}

func (_ *_RdsLocker) Find(id string) (rdsLocker *RdsLockerModel, err error) {
	if !bson.IsObjectIdHex(id) {
		return nil, ErrInvalidId
	}

	RdsLocker.Query(func(c *mgo.Collection) {
		err = c.FindId(bson.ObjectIdHex(id)).One(&rdsLocker)
	})

	return
}

func (_ *_RdsLocker) FindOrCreate(uid uint32, regionId, dBInstanceId string) (rdsLocker *RdsLockerModel, err error) {
	rdsLocker, err = RdsLocker.FindByUidRegionIdAndDBInstanceId(uid, regionId, dBInstanceId)
	if err == ErrNotFound {
		err = nil
		rdsLocker = NewRdsLockerModel(uid, regionId, dBInstanceId)
		// concurrency will result in duplicate error
		return rdsLocker, rdsLocker.Save()
	}

	return
}

func (_ *_RdsLocker) FindByUidRegionIdAndDBInstanceId(uid uint32, regionId, dBInstanceId string) (rdsLocker *RdsLockerModel, err error) {
	if uid == 0 || regionId == "" || dBInstanceId == "" {
		return nil, ErrInvalidParams
	}
	RdsLocker.Query(func(c *mgo.Collection) {
		query := bson.M{
			"uid":            uid,
			"region_id":      regionId,
			"db_instance_id": dBInstanceId,
		}
		err = c.Find(query).One(&rdsLocker)
	})

	return
}

func (_ *_RdsLocker) RemoveByUidRegionIdAndDBInstanceId(uid uint32, regionId, dBInstanceId string) (err error) {
	if uid == 0 || regionId == "" || dBInstanceId == "" {
		return ErrInvalidParams
	}

	RdsLocker.Query(func(c *mgo.Collection) {
		query := bson.M{
			"uid":            uid,
			"region_id":      regionId,
			"db_instance_id": dBInstanceId,
		}
		err = c.Remove(query)
	})

	return
}

func (_ *_RdsLocker) All() (rdsLockers []*RdsLockerModel, err error) {
	RdsLocker.Query(func(c *mgo.Collection) {
		err = c.Find(nil).All(&rdsLockers)
	})
	return

}

func (_ *_RdsLocker) Query(query func(c *mgo.Collection)) {
	MongoModel().Query(RdsLockerCollection, RdsLockerIndexes, query)
}
