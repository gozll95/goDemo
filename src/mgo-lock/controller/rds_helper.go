package controller

import (
	"time"

	"github.com/zhu/qvm/server/errors"
	common_rds "github.com/zhu/qvm/server/lib/aliyun/common/rds"
	"github.com/zhu/qvm/server/lib/locker"
	params_rds "github.com/zhu/qvm/server/lib/params/rds"
)

type doFunc func() error
type asyncFunc func() error

func GetRdsStatus(dbInstanceID string, request *RdsRequestWithRegion) (status common_rds.InstanceStatus, err error) {
	regionId := request.Region()
	rdsInstanceClient := request.Rds().NewInstance()

	describeArgs := params_rds.DescribeDBInstanceAttribute{
		RegionId:     regionId,
		DBInstanceID: dbInstanceID,
	}

	res, err := rdsInstanceClient.DescribeDBInstanceAttribute(&describeArgs)
	if err != nil {
		request.Logger().Errorf("RDS(?).DescribeDBInstanceAttribute().DescribeDBInstanceAttribute(%#v): %v", describeArgs, err)
		return common_rds.InstanceStatus("N/A"), err
	}
	return res.Items.DBInstanceAttribute[0].DBInstanceStatus, nil
}

func GetRdsStatusIntermittently(dbInstanceID string, expire time.Duration, request *RdsRequestWithRegion) (err error) {
	sendTicker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-time.After(expire):
			request.Logger().Errorf("GetRdsStatusIntermittently(%v,%v,%#v):timeout", dbInstanceID, expire)
			return
		case <-sendTicker.C:
			status, err := GetRdsStatus(dbInstanceID, request)
			if err != nil {
				return err
			}
			if status == common_rds.RunningInstanceStatus {
				return nil
			}
		}
	}
}

func test(dbInstanceID string, expire time.Duration, request *RdsRequestWithRegion) (err error) {
	var (
		uid       uint32
		regionId  string
		rdsLocker locker.Locker
	)

	// init vars
	uid = request.User().Uid
	regionId = request.Region()

	// wait status until be running
	GetRdsStatusIntermittently(dbInstanceID, expire, request)

	// call aliyun to describe
	describeArgs := common_rds.DescribeDBInstanceNetInfo{
		RegionId:     regionId,
		DBInstanceId: dbInstanceID,
	}

}
func RdsLockDo(name, dbInstanceId string, expire time.Duration, doFunc doFunc, asyncFunc asyncFunc, request *RdsRequestWithRegion) (err error) {
	var (
		uid       uint32
		regionId  string
		rdsLocker locker.Locker
	)

	// init vars
	uid = request.User().Uid
	regionId = request.Region()

	// get rds status
	status, err := GetRdsStatus(dbInstanceId, request)
	if err != nil {
		request.Logger().Errorf("GetRdsStatus(%v,%#v): %v", dbInstanceId, request, err)
		return errors.InternalError
	}
	switch status {
	case common_rds.RunningInstanceStatus:
		// get locker
		rdsLocker, err = locker.NewRdsLocker(uid, regionId, dbInstanceId, expire)
		if err != nil {
			request.Logger().Errorf("locker.NewRdsLocker(%v,%v,%v,%v): %v", uid, regionId, dbInstanceId, expire, err)
			return errors.InternalError
		}
		err = rdsLocker.Lock()
		if err != nil {
			request.Logger().Errorf("rdsLocker.Lock(): %v", err)
			return errors.NotSupportRdsStatus
		}

		err = doFunc()
		if err != nil {
			rdsLocker.Unlock()
			return errors.InternalError
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					request.Logger().Errorf("asyncFunc defer func: %v", uid, regionId, dbInstanceId, expire, err)
				}

				rdsLocker.Unlock()
			}()
			err = asyncFunc()
			if err != nil {
				request.Logger().Errorf("asyncFunc: %v", uid, regionId, dbInstanceId, expire, err)
			}
		}()

	default:
		return errors.NotSupportRdsStatus
	}
	return nil
}
