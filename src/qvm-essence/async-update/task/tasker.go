package task

import "github.com/zhu/qvm/server/enums"

type Tasker interface {
	TaskUid() uint32
	TaskRegion() string
	TaskResourceType() enums.ResourceType
	TaskResourceId() string
}
