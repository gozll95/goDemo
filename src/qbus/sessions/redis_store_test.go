package sessions

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/golib/assert"

	"gopkg.in/redis.v3"
)

var (
	testRedisOpts = &redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	}
)

type (
	testRedisData struct {
		Name  string
		Phone uint32
	}
)

func Test_Redis_Store(t *testing.T) {
	assertion := assert.New(t)

	client := redis.NewClient(testRedisOpts)
	defer func() {
		client.FlushDb()
	}()
	assertion.Nil(client.Ping().Err())

	store, err := NewRedisProvider(Config{
		CookieName: "redis",
		SecretKey:  "sider",
	}, client)
	assertion.Nil(err)

	var (
		sid     = NewSessionID()
		name    = NewSessionID()
		phone   = rand.Uint32()
		counter = rand.Int()
	)
	client.Del(sid)

	cookie, err := store.Create(sid)
	assertion.Nil(err)

	cookie.Set("user", &testRedisData{name, phone})
	cookie.Set("counter", counter)

	// flush to redis
	err = cookie.Flush()
	assertion.Nil(err)

	// restore from redis
	rcookie, err := store.Restore(sid)
	assertion.Nil(err)

	msg, ok := rcookie.Get("user")
	assertion.True(ok)

	var user *testRedisData
	err = json.Unmarshal(msg, &user)
	assertion.Nil(err)
	assertion.Equal(name, user.Name)
	assertion.Equal(phone, user.Phone)

	rcounter, err := rcookie.Value("counter").Int()
	assertion.Nil(err)
	assertion.Equal(counter, rcounter)
}
