package locker

import "github.com/zhu/qvm/server/enums"

type Locker interface {
	Status() (status enums.LockStatus, err error)
	Lock() (err error)
	Unlock() (err error)
}
