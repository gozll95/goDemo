package locker

import (
	"math/rand"
	"sync"
	"testing"

	"time"

	"github.com/zhu/qvm/server/enums"
	"github.com/zhu/qvm/server/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_RdsLocker_NewRdsLocker(t *testing.T) {
	uid := rand.Uint32()
	regionId := "cn-beijing"
	dBInstanceId := uuid.NewV4().String()
	expireTime, _ := time.ParseDuration("3s")

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			_, err := NewRdsLocker(uid, regionId, dBInstanceId, expireTime)
			assert.Nil(t, err)
		}()
	}
	wg.Wait()

	rdsLockers, err := model.RdsLocker.All()
	assert.Nil(t, err)
	assert.Equal(t, len(rdsLockers), 1)
}

func Test_RdsLocker_Lock(t *testing.T) {
	uid := rand.Uint32()
	regionId := "cn-beijing"
	dBInstanceId := uuid.NewV4().String()
	expireTime, _ := time.ParseDuration("3s")

	locker, err := NewRdsLocker(uid, regionId, dBInstanceId, expireTime)
	assert.Nil(t, err)

	err = locker.Lock()
	assert.Nil(t, err)

	status, err := locker.Status()
	assert.Nil(t, err)
	assert.Equal(t, status, enums.Locked)

	time.Sleep(5 * time.Second)

	status, err = locker.Status()
	assert.Nil(t, err)
	assert.Equal(t, status, enums.Unlock)
}

func Test_RdsLocker_Unlock(t *testing.T) {
	uid := rand.Uint32()
	regionId := "cn-beijing"
	dBInstanceId := uuid.NewV4().String()
	expireTime, _ := time.ParseDuration("-3s")

	locker, err := NewRdsLocker(uid, regionId, dBInstanceId, expireTime)
	assert.Nil(t, err)

	err = locker.Lock()
	assert.Nil(t, err)

	status, err := locker.Status()
	assert.Nil(t, err)
	assert.Equal(t, status, enums.Locked)

	err = locker.Unlock()
	assert.Nil(t, err)

	status, err = locker.Status()
	assert.Nil(t, err)
	assert.Equal(t, status, enums.Unlock)
}
