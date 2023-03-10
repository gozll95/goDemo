package endpoint

import (
	"context"
	"fmt"

	// "go-kit-gin/entity"
	"go-kit-gin/protocol"
	"go-kit-gin/service"
	// "github.com/go-kit/kit/endpoint"
)

//MakeAccountEndpoint 生成Account断点
func MakeAccountEndpoint(svc service.AppService) Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(protocol.AccountReq)
		var resp protocol.Resp
		resp, _ = svc.Account(req)
		return resp, nil
	}
}

//MakeCreateAccountEndpoint 生成CreateAccount端点
func MakeCreateAccountEndpoint(svc service.AppService) Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(protocol.CreateAccountReq)
		var resp protocol.Resp
		resp, _ = svc.CreateAccount(req)

		return resp, nil
	}
}

//MakeUpdateAccountEndpoint 更新账户端点
func MakeUpdateAccountEndpoint(svc service.AppService) Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(protocol.UpdateAccountReq)
		var resp protocol.Resp
		fmt.Println(req)
		return resp, nil
	}
}

//MakeLoginEndpoint 登录端点
func MakeLoginEndpoint(svc service.AppService) Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(protocol.LoginReq)
		var resp protocol.Resp
		resp, _ = svc.Login(req)
		return resp, nil
	}
}
