package graph

import (
	"errors"
	"fmt"
)

var (
	GetPositionErr           = errors.New("GetPositionErr")
	GraphCycleErr            = errors.New("Graph has a cycle")
	MartixVexMustNotNil      = errors.New("martix vex must not be nil")
	MartixOverSize           = fmt.Errorf("martix over size %v", MAXSIZE)
	VexAlreadyExistInMartix  = errors.New("vex alerady exist in martix")
	EdgeAlreadyExistInMartix = errors.New("edge alerady exist in martix")
	VexNotFoundInMartix      = errors.New("vex not found in martix")
)
