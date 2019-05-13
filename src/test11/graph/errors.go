package graph

import (
	"errors"
	"fmt"
)

var (
	GetPositionErr      = errors.New("GetPositionErr")
	GraphCycleErr       = errors.New("Graph has a cycle")
	MartixVexMustNotNil = errors.New("martix vex must not be nil")
	MartixOverSize      = fmt.Errorf("martix over size %v", MAXSIZE)
)
