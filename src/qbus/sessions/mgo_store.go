package sessions

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	idField        = "_id"
	sidField       = "sid"
	dataField      = "data"
	createdAtField = "ctime"
	updatedAtField = "atime"
)

type mgoStore struct {
	Id        bson.ObjectId              `bson:"_id"`
	Sid       string                     `bson:"sid"`
	Data      map[string]json.RawMessage `bson:"data"`
	CreatedAt time.Time                  `bson:"ctime"`
	UpdatedAt time.Time                  `bson:"atime"`
}

func newMgoStore(sid string, params ...map[string]interface{}) *mgoStore {
	store := &mgoStore{
		Id:        bson.NewObjectId(),
		Sid:       sid,
		Data:      make(map[string]json.RawMessage),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if len(params) > 0 {
		for key, value := range params[0] {
			b, err := json.Marshal(value)
			if err == nil {
				store.Data[key] = json.RawMessage(b)
			}
		}
	}

	return store
}

type mgoData struct {
	store    *mgoStore
	provider *MgoProvider

	changed bool // flag when values has changed
	closed  bool // flag when data destroy
}

var _ Store = new(mgoData)

func newMgoData(p *MgoProvider, store *mgoStore, changed bool) *mgoData {
	data := &mgoData{
		store:    store,
		provider: p,
		changed:  changed,
	}

	if store.Data == nil {
		store.Data = make(map[string]json.RawMessage)
		data.changed = true
	}

	return data
}

func (m *mgoData) ID() string {
	return m.store.Sid
}

func (m *mgoData) Has(key string) bool {
	_, ok := m.store.Data[key]

	return ok
}

func (m *mgoData) Set(key string, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		return
	}

	m.store.Data[key] = json.RawMessage(b)
	m.changed = true
}

func (m *mgoData) Get(key string) (msg json.RawMessage, ok bool) {
	msg, ok = m.store.Data[key]
	return
}

func (m *mgoData) Value(key string) *Value {
	return NewSessionValue(m.Get(key))
}

func (m *mgoData) Delete(key string) error {
	delete(m.store.Data, key)
	m.changed = true

	return nil
}

// Dump returns all values copied
func (m *mgoData) Dump() map[string]json.RawMessage {
	values := make(map[string]json.RawMessage, len(m.store.Data))
	for key, value := range m.store.Data {
		values[key] = value
	}

	return values
}

// Flush persists session values to store
func (m *mgoData) Flush() error {
	// has destroy
	if m.closed {
		return nil
	}

	// no changes
	if !m.changed {
		return nil
	}

	err := m.provider.save(m.store)
	if err == nil {
		m.changed = false
	}

	return err
}

func (m *mgoData) Touch() error {
	err := m.provider.connect(func(c *mgo.Collection) error {
		merr := c.Update(bson.M{sidField: m.store.Sid}, bson.M{
			"$set": bson.M{
				updatedAtField: time.Now(),
			},
		})

		return merr
	})

	return err
}

// Clean removes all values in session
func (m *mgoData) Clean() {
	m.store.Data = make(map[string]json.RawMessage)
	m.changed = true
}

// Destroy removes session in store
func (m *mgoData) Destroy() error {
	err := m.provider.Destroy(m.store.Sid)
	if err == nil {
		m.store.Data = make(map[string]json.RawMessage)
		m.closed = true
	}

	return err
}

type MgoProvider struct {
	config  *Config
	connect func(func(c *mgo.Collection) error) error
}

var _ Provider = new(MgoProvider)

func NewMgoProvider(config Config, connect func(func(c *mgo.Collection) error) error) (provider Provider, err error) {
	err = connect(func(c *mgo.Collection) (merr error) {
		merr = c.Database.Session.Ping()
		if merr != nil {
			return
		}

		// ensure sid is unque index
		merr = c.EnsureIndex(mgo.Index{Key: []string{sidField}, Name: sidField, Unique: true})
		if merr != nil {
			return
		}

		// auto expire?
		if config.AutoExpire {
			index := mgo.Index{
				Name:        updatedAtField,
				Key:         []string{updatedAtField},
				ExpireAfter: time.Duration(config.SessionExpire) * time.Second,
			}

			merr = c.EnsureIndex(index)
			if merr != nil {
				return
			}
		}

		return
	})
	if err != nil {
		return
	}

	provider = &MgoProvider{
		config:  &config,
		connect: connect,
	}

	return
}

func (p *MgoProvider) Config() *Config {
	config := *(p.config)

	return &config
}

func (p *MgoProvider) Create(sid string, params ...map[string]interface{}) (sess Store, err error) {
	err = p.connect(func(c *mgo.Collection) (merr error) {
		store := newMgoStore(sid, params...)

		merr = c.Insert(store)
		if mgo.IsDup(merr) {
			merr = ErrDuplicateID
		}

		if merr == nil {
			sess = newMgoData(p, store, len(params) > 0)
		}

		return
	})

	return
}

func (p *MgoProvider) Restore(sid string) (sess Store, err error) {
	err = p.connect(func(c *mgo.Collection) (merr error) {
		var store *mgoStore

		merr = c.Find(bson.M{sidField: sid}).One(&store)
		if merr == mgo.ErrNotFound {
			merr = ErrNotFound
			return
		}

		if merr == nil {
			if !p.config.AutoExpire {
				// is session has expired?
				if time.Since(store.UpdatedAt) > p.config.SessionExpireSeconds() {
					merr = ErrNotFound

					// error can secure skip
					_ = p.Destroy(sid)
					return
				}
			}

			sess = newMgoData(p, store, false)
		}

		return
	})

	return
}

// TODO: need to verify expiration!
func (p *MgoProvider) Refresh(old string, sid string) (sess Store, err error) {
	err = p.connect(func(c *mgo.Collection) (merr error) {
		merr = c.Update(bson.M{sidField: old}, bson.M{
			"$set": bson.M{
				sidField:       sid,
				updatedAtField: time.Now(),
			},
		})

		if mgo.IsDup(merr) {
			merr = ErrDuplicateID
			return
		}

		if merr == mgo.ErrNotFound {
			merr = ErrNotFound
			return
		}

		return
	})

	if err != nil {
		return nil, err
	}

	return p.Restore(sid)
}

func (p *MgoProvider) Destroy(sid string) (er error) {
	er = p.connect(func(c *mgo.Collection) (err error) {
		err = c.Remove(bson.M{
			sidField: sid,
		})
		return err
	})
	return
}

func (p *MgoProvider) GC() (err error) {
	if p.config.AutoExpire {
		return
	}

	err = p.connect(func(c *mgo.Collection) error {
		_, merr := c.RemoveAll(bson.M{
			updatedAtField: bson.M{
				"$lte": time.Now().Add(-p.config.SessionExpireSeconds()),
			},
		})

		return merr
	})

	return
}

func (p *MgoProvider) save(store *mgoStore) (err error) {
	err = p.connect(func(c *mgo.Collection) error {
		_, merr := c.Upsert(bson.M{
			sidField: store.Sid,
		}, bson.M{
			idField:        store.Id,
			sidField:       store.Sid,
			dataField:      store.Data,
			createdAtField: store.CreatedAt,
			updatedAtField: time.Now(),
		})

		return merr
	})

	return
}
