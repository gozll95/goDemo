package sessions

import (
	"encoding/json"
	"time"

	"gopkg.in/redis.v3"
)

type redisStore struct {
	Sid       string                     `json:"sid"`
	Data      map[string]json.RawMessage `json:"data"`
	CreatedAt time.Time                  `json:"ctime"`
	UpdatedAt time.Time                  `json:"atime"`
}

func (r *redisStore) CacheKey() string {
	return CacheKey(r.Sid)
}

func newRedisStore(sid string, params ...map[string]interface{}) *redisStore {
	store := &redisStore{
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

type redisData struct {
	store    *redisStore
	provider *RedisProvider

	changed bool // flag when values has changed
	closed  bool // flag when data destroy
}

var _ Store = new(redisData)

func newRedisData(p *RedisProvider, store *redisStore, changed bool) *redisData {
	data := &redisData{
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

func (m *redisData) ID() string {
	return m.store.Sid
}

func (m *redisData) Has(key string) bool {
	_, ok := m.store.Data[key]

	return ok
}

func (m *redisData) Set(key string, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		return
	}

	m.store.Data[key] = json.RawMessage(b)
	m.changed = true
}

func (m *redisData) Get(key string) (msg json.RawMessage, ok bool) {
	msg, ok = m.store.Data[key]
	if !ok {
		return
	}

	var buf []byte
	err := json.Unmarshal(msg, &buf)
	if err != nil {
		ok = false
	} else {
		msg = json.RawMessage(buf)
	}

	return
}

func (m *redisData) Value(key string) *Value {
	return NewSessionValue(m.Get(key))
}

func (m *redisData) Delete(key string) error {
	delete(m.store.Data, key)
	m.changed = true

	return nil
}

// Dump returns all values copied
func (m *redisData) Dump() map[string]json.RawMessage {
	values := make(map[string]json.RawMessage, len(m.store.Data))
	for key, value := range m.store.Data {
		values[key] = value
	}

	return values
}

// Flush persists session values to store
func (m *redisData) Flush() error {
	// has destroy
	if m.closed {
		return nil
	}

	// no changes
	if !m.changed {
		return nil
	}

	err := m.save()
	if err == nil {
		m.changed = false
	}

	return err
}

func (m *redisData) Touch() error {
	return m.provider.client.Expire(m.store.CacheKey(), m.provider.config.SessionExpireSeconds()).Err()
}

// Clean removes all values in session
func (m *redisData) Clean() {
	m.store.Data = make(map[string]json.RawMessage)
	m.changed = true
}

// Destroy removes session in store
func (m *redisData) Destroy() error {
	err := m.provider.Destroy(m.store.Sid)
	if err == nil {
		m.store.Data = make(map[string]json.RawMessage)
		m.closed = true
	}

	return err
}

func (m *redisData) save() error {
	b, err := json.Marshal(m.store)
	if err != nil {
		return err
	}

	return m.provider.client.Set(m.store.CacheKey(), b, m.provider.config.SessionExpireSeconds()).Err()
}

type RedisProvider struct {
	config *Config
	client *redis.Client
}

var _ Provider = new(RedisProvider)

func NewRedisProvider(config Config, client *redis.Client) (provider Provider, err error) {
	err = client.Ping().Err()
	if err != nil {
		return
	}

	provider = &RedisProvider{
		config: &config,
		client: client,
	}

	return
}

func (p *RedisProvider) Config() *Config {
	config := *(p.config)

	return &config
}

func (p *RedisProvider) Create(sid string, params ...map[string]interface{}) (sess Store, err error) {
	ok, err := p.client.Exists(CacheKey(sid)).Result()
	if err != nil {
		return
	}

	if ok {
		err = ErrDuplicateID
		return
	}

	sess = newRedisData(p, newRedisStore(sid, params...), len(params) > 0)
	err = sess.Flush()
	return
}

func (p *RedisProvider) Restore(sid string) (sess Store, err error) {
	b, err := p.client.Get(CacheKey(sid)).Bytes()
	if err == redis.Nil {
		err = ErrNotFound
	}
	if err != nil {
		return
	}

	var rs *redisStore
	err = json.Unmarshal(b, &rs)
	if err != nil {
		return
	}

	sess = newRedisData(p, rs, false)
	return
}

func (p *RedisProvider) Refresh(old string, sid string) (store Store, err error) {
	ok, err := p.client.Exists(CacheKey(sid)).Result()
	if err != nil {
		return
	}

	if ok {
		err = ErrDuplicateID
		return
	}

	store, err = p.Restore(old)
	if err != nil {
		return
	}

	store.(*redisData).store.Sid = sid

	err = store.Flush()
	return
}

func (p *RedisProvider) Destroy(sid string) (er error) {
	return p.client.Del(CacheKey(sid)).Err()
}
