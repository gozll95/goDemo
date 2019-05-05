package main

import (
	"fmt"

	"github.com/juju/errors"
)

func main() {
	var someErr = errors.New("some error")
	fmtErr := fmt.Errorf("wrapper")
	err := errors.Wrap(someErr, fmtErr)
	err = errors.Wrap(err, fmt.Errorf("wrapper2"))
	err = errors.Annotate(err, "annotate1")
	err = errors.Annotate(err, "annotate2")
	fmt.Println(errors.ErrorStack(err))

	fmt.Println(err.Error())

	fmt.Println(errors.Cause(err))

}



// https://github.com/juju/errors
/*
/Users/flower/workspace/goDemo/src/juju-errors.go:10: some error
/Users/flower/workspace/goDemo/src/juju-errors.go:12: wrapper
/Users/flower/workspace/goDemo/src/juju-errors.go:13: wrapper2
/Users/flower/workspace/goDemo/src/juju-errors.go:14: annotate1
/Users/flower/workspace/goDemo/src/juju-errors.go:15: annotate2
annotate2: annotate1: wrapper2
wrapper2
*/
