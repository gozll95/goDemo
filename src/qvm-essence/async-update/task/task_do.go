package task

import (
	"github.com/zhu/qvm/server/enums"
)

// request aliyun backend update db and return checker
func (t *TaskManager) do(uid uint32, regionId string, resourceType enums.ResourceType, resourceIds []string) (checkers []Checker, err error) {
	switch resourceType {
	case enums.ResourceTypeInstance:
		return instanceUpdate(uid, regionId, resourceType, resourceIds)

	case enums.ResourceTypeDisk:
		return diskUpdate(uid, regionId, resourceType, resourceIds)

	case enums.ResourceTypeIp:
		return eipUpdate(uid, regionId, resourceType, resourceIds)

	case enums.ResourceTypeSnapshot:
		return snapUpdate(uid, regionId, resourceType, resourceIds)

	case enums.ResourceTypeSecurityGroup:

	case enums.ResourceTypeVswitch:
		return vswitchUpdate(uid, regionId, resourceType, resourceIds)

	case enums.ResourceTypeVpc:
		return vpcUpdate(uid, regionId, resourceType, resourceIds)

	case enums.ResourceTypeImage:
		return imageUpdate(uid, regionId, resourceType, resourceIds)

	case enums.ResourceTypeListener:
		//这里listener的resourceIds为protocol:lbId:ListenerPort
		return listenerUpdate(uid, regionId, resourceType, resourceIds)

	case enums.ResourceTypeLoadBalancer:
		return loadbalancerUpdate(uid, regionId, resourceType, resourceIds)

	case enums.ResourceTypeBsn:
		return bsnUpdate(uid, regionId, resourceType, resourceIds)
	}

	return
}
