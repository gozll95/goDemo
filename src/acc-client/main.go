package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"zhu.us/oauth"
)

const (
	ClientId     = ""
	ClientSecret = ""
	Host         = "root"
	User         = "root"
	Password     = "root"
)

type AccToken struct {
	user      string
	passwd    string
	xauth     *oauth.Transport
	expiredAt time.Time
	token     string
}

func NewAccToken(user, passwd string, xauth *oauth.Transport) http.RoundTripper {
	return &AccToken{
		user:      user,
		passwd:    passwd,
		xauth:     xauth,
		expiredAt: time.Now(),
	}
}

func (acc *AccToken) getToken() (token string, err error) {
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

func (acc *AccToken) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := acc.getToken()
	if err != nil {
		return nil, err
	}
	authorization := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", authorization)
	return http.DefaultTransport.RoundTrip(req)
}

func QTransport() *oauth.Transport {
	return &oauth.Transport{
		Config: &oauth.Config{
			ClientId:     ClientId,
			ClientSecret: ClientSecret,
			Scope:        "Scope",
			AuthURL:      "<AuthURL>",
			TokenURL:     Host + "/oauth2/token",
			RedirectURL:  "<RedirectURL>",
		},
	}
}

func main() {
	//setup transport
	accToken := NewAccToken(User, Password, QTransport())
	transport := NewCustomTransport(accToken)

	account := AccountInfo{}
	resp, err := GetAccount(transport.Client(), Host, 1381084496, &account)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
	}
	fmt.Println(account)

}

type AccErr struct {
	ErrorCode        uint64 `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

func (a *AccErr) Error() error {
	msg := fmt.Sprintf("acc error_code:%v error_description:%v", a.ErrorCode, a.ErrorDescription)
	return errors.New(msg)
}

type AccountInfo struct {
	Uid         uint32 `json:"uid"`
	Email       string `json:"email"`
	Utype       int    `json:"utype"`
	DisableType int    `json:"disabled_type"`
}

func GetAccount(client *http.Client, host string, uid uint32, result interface{}) (*http.Response, error) {
	u := fmt.Sprintf("%s/admin/user/info?uid=%d", host, uid)
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorResponse := AccErr{}
		if json.NewDecoder(resp.Body).Decode(&errorResponse) != nil {
			return resp, fmt.Errorf("acc response status %d", resp.StatusCode)
		}
		return resp, errorResponse.Error()
	}
	return resp, json.NewDecoder(resp.Body).Decode(&result)
}

type CustomHttpClient interface {
	Client() *http.Client
}

type CustomTransport struct {
	Transport http.RoundTripper
}

func NewCustomTransport(t http.RoundTripper) CustomHttpClient {
	return &CustomTransport{
		Transport: t,
	}
}

func (t *CustomTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.transport().RoundTrip(req)
}

func (t *CustomTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}
