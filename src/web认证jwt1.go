package main

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func main() {
	mySigningKey := []byte("hzwy23")
	// Create the Claims
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix() - 1000),
		ExpiresAt: int64(time.Now().Unix() + 1000),
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Println("签名后的token信息:", ss)
	t, err := jwt.Parse(ss, func(*jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		fmt.Println("parase with claims failed.", err)
		return
	}
	fmt.Println("还原后的token信息claims部分:", t.Claims)
}

//http://blog.csdn.net/hzwy23/article/details/53224724

/*
json web token 简介

json web token 简称 jwt.他是一种轻量级的规范.这种规范允许客户端和服务端之间传递一些非敏感信息.
常用于用户认证和授权系统.
jwt组成部分

Header
Claims
Signature
Header 组成部分

typ: "JWT",
alg: "HS256",
1. typ是默认的一种标识.标识这条信息采用JWT规范.
2. alg表示签名使用的加密算法.通常有ES256,ES512,RS256等等

Claims 组成部分

Audience  string `json:"aud,omitempty"`
ExpiresAt int64  `json:"exp,omitempty"`
Id        string `json:"jti,omitempty"`
IssuedAt  int64  `json:"iat,omitempty"`
Issuer    string `json:"iss,omitempty"`
NotBefore int64  `json:"nbf,omitempty"`
Subject   string `json:"sub,omitempty"`
这段结构体是我们在golang中使用到的字段. 可以在这个的基础上进行组合,定义新的Claims部分.

1. aud 标识token的接收者.
2. exp 过期时间.通常与Unix UTC时间做对比过期后token无效
3. jti 是自定义的id号
4. iat 签名发行时间.
5. iss 是签名的发行者.
6. nbf 这条token信息生效时间.这个值可以不设置,但是设定后,一定要大于当前Unix UTC,否则token将会延迟生效.
7. sub 签名面向的用户

通过设置exp与nbf来管理token的生命周期.

Signature 组成部分

将Header与Claims信息拼接起来[base64(header)+”.”+base64(claims)],采用Header中指定的加密算法进行加密,得到Signature部分.
token组成部分

base64(header) + “.” + base64(claims) + “.” + 加密签名
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI5NTg5MTA4MDYsImlzcyI6InRlc3QiLCJuYmYiOjE0Nzk0NTczMTZ9.57gqtlk1nNezXSa0VgWBOwu2b2FCDJ6wXizuJF6IY10
1
1
如上边的token,由2个点好分割.第一部分是iheader的base64编码,第二部分是claims的base64编码,第三部分是加密签名信息.


*/
