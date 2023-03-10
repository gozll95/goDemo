package protocol

import (
	"go-kit-gin/entity"
)

//CreateAccountReq 创建账户请求
type CreateAccountReq struct {
	entity.UserAccount
}

//CreateAccountResp 创建账户响应
type CreateAccountResp struct {
}

//AccountReq 查询账户请求
type AccountReq struct {
	Accesstoken string
}

//AccountResp 查询账户响应
type AccountResp struct {
	entity.UserAccount
}

//UpdateAccountReq 更新账户请求
type UpdateAccountReq struct {
	entity.UserAccount
}

//UpdateAccountResp 更新账户请求响应
type UpdateAccountResp struct {
}

type LoginReq struct {
	entity.UserAccount
}

type LoginResp struct {
	AccessToken string
}
