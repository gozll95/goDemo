package client

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zhu/qvm/server/lib/bbbb/oauth"
)

type OauthTransport struct {
	UserName string
	Password string
	TokenCli *oauth.Transport

	trans http.RoundTripper

	token *oauth.Token
	sync.Mutex
}

func NewOauthTransport(usrname, password, tokenUrl string, trans http.RoundTripper) http.RoundTripper {
	if trans == nil {
		trans = http.DefaultTransport
	}

	tokenCli := oauth.Transport{
		Config: &oauth.Config{
			TokenURL: strings.TrimSuffix(tokenUrl, "/") + "/oauth2/token",
		},
	}

	return &OauthTransport{
		UserName: usrname,
		Password: password,
		TokenCli: &tokenCli,
		trans:    trans,
	}
}

func (oa *OauthTransport) expired() bool {
	if oa.token == nil {
		return true
	}
	if oa.token.TokenExpiry == 0 {
		return false
	}
	// 提前120秒换取token
	return oa.token.TokenExpiry <= (time.Now().Unix() + 120)
}

func (oa *OauthTransport) refresh() error {
	oa.Lock()
	defer oa.Unlock()
	if !oa.expired() {
		return nil
	}

	var err error
	oa.token, _, err = oa.TokenCli.ExchangeByPassword(oa.UserName, oa.Password)
	if err != nil {
		return err
	}

	return nil
}

func (oa *OauthTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if oa.expired() {
		err = oa.refresh()
		if err != nil {
			return nil, err
		}
	}
	req.Header.Set("Authorization", "Bearer "+oa.token.AccessToken)
	return oa.trans.RoundTrip(req)
}
