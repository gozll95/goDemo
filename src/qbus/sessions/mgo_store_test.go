package sessions

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/golib/assert"

	mgo "gopkg.in/mgo.v2"
)

var (
	testMgoDB    = "mgo_store"
	testMgoTable = "sessions"
	testMgoDSN   = "mongodb://localhost:27017/" + testMgoDB
)

type (
	testMgoData struct {
		Name  string
		Phone uint32
	}
)

func Test_Mgo_Store(t *testing.T) {
	assertion := assert.New(t)

	conn, err := mgo.Dial(testMgoDSN)
	assertion.Nil(err)
	defer func() {
		conn.Copy().DB(testMgoDB).DropDatabase()
	}()

	client := conn.Copy().DB(testMgoDB).C(testMgoTable)

	store, err := NewMgoProvider(Config{
		CookieName:    "mgo",
		SessionExpire: 1,
		SecretKey:     "ogm",
	}, func(query func(c *mgo.Collection) error) error {
		return query(client)
	})
	assertion.Nil(err)

	var (
		sid     = NewSessionID()
		name    = NewSessionID()
		phone   = rand.Uint32()
		counter = rand.Int()
	)

	cookie, err := store.Create(sid)
	assertion.Nil(err)

	cookie.Set("user", &testMgoData{name, phone})
	cookie.Set("counter", counter)

	// flush to mongodb
	err = cookie.Flush()
	assertion.Nil(err)

	// restore from mongodb
	rcookie, err := store.Restore(sid)
	assertion.Nil(err)

	msg, ok := rcookie.Get("user")
	assertion.True(ok)

	var user *testMgoData
	err = json.Unmarshal(msg, &user)
	assertion.Nil(err)
	assertion.Equal(name, user.Name)
	assertion.Equal(phone, user.Phone)

	rcounter, err := rcookie.Value("counter").Int()
	assertion.Nil(err)
	assertion.Equal(counter, rcounter)
}
