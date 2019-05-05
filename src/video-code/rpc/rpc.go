package rpcdemo

import "errors"

// Service.Method

type DemoService struct {
}

// Public
type Args struct {
	A, B int
}

// rpc的参数一定要是两个,第二个必须是指针类型
func (DemoService) Div(args Args, result *float64) error {
	if args.B == 0 {
		return errors.New("division by zero")
	}
	*result = float64(args.A) / float64(args.B)
	return nil
}
