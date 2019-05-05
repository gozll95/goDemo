package sessions

import "time"

const (
	SESSION_ID_LENGTH        = 32
	COOKIE_VALUE_SPLIT       = "," // value,value,value
	COOKIE_VALUE_PARTS_SPLIT = "|" // value1|value2|value3,name,time
)

type Config struct {
	CookieName   string // session cookie name
	CookieSecure bool   // is cookie use https?
	CookieExpire int    // session cookie expire seconds

	CookieRememberName string // hashed value of user for auto login
	RememberExpire     int    // auto login remember expire seconds

	SessionExpire int  // session expire seconds
	AutoExpire    bool // is provider support auto expire?

	SecretKey string // secure secret key
}

func (c *Config) SessionExpireSeconds() time.Duration {
	return time.Duration(c.SessionExpire) * time.Second
}

func (c *Config) RememberExpireSeconds() time.Duration {
	return time.Duration(c.RememberExpire) * time.Second
}
