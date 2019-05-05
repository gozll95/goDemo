package model

import "github.com/zhu/qvm/server/enums"

func (i *EipModel) TaskUid() uint32 {
	return i.Uid
}

func (i *EipModel) TaskRegion() string {
	return i.RegionId
}

func (i *EipModel) TaskResourceType() enums.ResourceType {
	return enums.ResourceTypeIp
}

func (i *EipModel) TaskResourceId() string {
	return i.AllocationId
}
