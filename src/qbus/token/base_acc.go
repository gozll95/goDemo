package token

import (
	"fmt"
	"time"

	"zhu.us/oauth"
)

type AccToken struct {
	user      string
	passwd    string
	xauth     *oauth.Transport
	expiredAt time.Time
	token     string
}

func NewAccToken(user, passwd string, xauth *oauth.Transport) Token {
	fmt.Println("newacctoken")
	return &AccToken{
		user:      user,
		passwd:    passwd,
		xauth:     xauth,
		expiredAt: time.Now(),
	}
}

func (acc *AccToken) GetToken() (token string, err error) {
	if time.Now().After(acc.expiredAt) {

		adminToken, _, err := acc.xauth.ExchangeByPassword(acc.user, acc.passwd)

		if err != nil {
			return "", err
		}
		expiredAt := time.Unix(adminToken.TokenExpiry, 0)
		acc.expiredAt = expiredAt
		acc.token = adminToken.AccessToken
		return adminToken.AccessToken, nil
	}
	return acc.token, nil
}
