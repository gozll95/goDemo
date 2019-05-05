package model

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"time"

	"github.com/zhu/qvm/server/enums"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
)

func Test_RdsLocker_Find(t *testing.T) {
	uid := rand.Uint32()
	regionId := "cn-beijing"
	dBInstanceId := uuid.NewV4().String()
	rdsLocker := NewRdsLockerModel(uid, regionId, dBInstanceId)

	err := rdsLocker.Save()
	assert.Nil(t, err)

	// find
	findRdsLocker, err := RdsLocker.Find(rdsLocker.Id.Hex())
	assert.Nil(t, err)
	assert.Equal(t, rdsLocker.Uid, findRdsLocker.Uid)
	assert.Equal(t, rdsLocker.RegionId, findRdsLocker.RegionId)
	assert.Equal(t, rdsLocker.DBInstanceId, findRdsLocker.DBInstanceId)
	assert.Equal(t, rdsLocker.LockStatus, enums.Unlock)

	// find and lock
	err = findRdsLocker.FindAndChangeLockStatus(enums.Locked)
	assert.Nil(t, err)

	// duplicate find and lock
	err = findRdsLocker.FindAndChangeLockStatus(enums.Locked)
	assert.NotNil(t, err)

	// find and unlock
	err = findRdsLocker.FindAndChangeLockStatus(enums.Unlock)
	assert.Nil(t, err)

	// duplicate find and unlock
	err = findRdsLocker.FindAndChangeLockStatus(enums.Unlock)
	assert.NotNil(t, err)

	// find and lock again
	err = findRdsLocker.FindAndChangeLockStatus(enums.Locked)
	assert.Nil(t, err)

	// all
	findRdsLockers, err := RdsLocker.All()
	assert.Nil(t, err)
	assert.Equal(t, len(findRdsLockers), 1)

	//SetUnlockWhenExpiredLocked
	expireTime, _ := time.ParseDuration("-5s")

	time.Sleep(3 * time.Second)
	err = findRdsLocker.SetUnlockWhenExpiredLocked(expireTime)
	assert.NotNil(t, err)

	time.Sleep(3 * time.Second)
	err = findRdsLocker.SetUnlockWhenExpiredLocked(expireTime)
	assert.Nil(t, err)
	assert.Equal(t, findRdsLocker.LockStatus, enums.Unlock)

	fmt.Println(findRdsLocker)

	// remove
	err = RdsLocker.RemoveByUidRegionIdAndDBInstanceId(uid, regionId, dBInstanceId)
	assert.Nil(t, err)
}

func Test_RdsLocker_FindOrCreate(t *testing.T) {
	uid := rand.Uint32()
	regionId := "cn-beijing"
	dBInstanceId := uuid.NewV4().String()

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			_, err := RdsLocker.FindOrCreate(uid, regionId, dBInstanceId)
			if mgo.IsDup(err) {
				return
			}
			if err != nil {
				assert.NotNil(t, err)
			}
		}()
	}
	wg.Wait()

	findRdsLockers, err := RdsLocker.All()
	assert.Nil(t, err)
	assert.Equal(t, len(findRdsLockers), 1)

	// remove
	err = RdsLocker.RemoveByUidRegionIdAndDBInstanceId(uid, regionId, dBInstanceId)
	assert.Nil(t, err)

}
