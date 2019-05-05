package locker

import (
	"errors"
	"fmt"
	"time"

	"github.com/zhu/qvm/server/enums"
	"github.com/zhu/qvm/server/model"
	"github.com/zhu/qvm/server/utils"
	mgo "gopkg.in/mgo.v2"
)

type RdsLocker struct {
	locker *model.RdsLockerModel
	expire time.Duration
}

func NewRdsLocker(uid uint32, regionId, dbInstanceId string, expire time.Duration) (rdsLocker Locker, err error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			rdsLocker, err := model.RdsLocker.FindOrCreate(uid, regionId, dbInstanceId)
			if mgo.IsDup(err) {
				continue
			}
			if err == nil {
				return &RdsLocker{
					locker: rdsLocker,
					expire: expire,
				}, nil
			}
			if err != nil {
				return nil, err
			}
		case <-time.After(5 * time.Second):
			errMsg := fmt.Sprintf("NewRdsLocker(%v,%v,%v,%v):timeout", uid, regionId, dbInstanceId, expire)
			return nil, errors.New(errMsg)
		}
	}
}

func (r *RdsLocker) Lock() (err error) {
	err = r.locker.FindAndChangeLockStatus(enums.Locked)
	if err != nil {
		return
	}
	time.AfterFunc(r.expire, func() {
		err = r.locker.SetUnlockWhenExpiredLocked(r.expire * -1)
		if err != nil {
			utils.StdLog.Errorf("r.locker(%#v).SetUnlockWhenExpiredLocked(%v):%v", r.locker, err)
		}
	})
	return nil
}

func (r *RdsLocker) Unlock() (err error) {
	return r.locker.FindAndChangeLockStatus(enums.Unlock)
}

func (r *RdsLocker) Status() (status enums.LockStatus, err error) {
	findRdsLocker, err := model.RdsLocker.Find(r.locker.Id.Hex())
	if err != nil {
		return
	}
	return findRdsLocker.LockStatus, nil
}
