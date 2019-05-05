package sessions

import "encoding/json"

type Store interface {
	ID() string                             // return current session ID
	Has(key string) bool                    // is key exist?
	Set(key string, value interface{})      // set a session value with the key, returns an error if existed
	Get(key string) (json.RawMessage, bool) // get a session value by key
	Value(key string) *Value                // get a session value by key and convert it to Value object for easying usage
	Delete(key string) error                // delete a session value by key
	Dump() map[string]json.RawMessage       // dump all values
	Flush() error                           // save all memory data to the provider
	Touch() error                           // sync session store expire time
	Clean()                                 // clean all memory data
	Destroy() error                         // delete session in store
}

type Provider interface {
	Config() *Config
	Create(sid string, params ...map[string]interface{}) (Store, error)
	Restore(sid string) (Store, error)
	Refresh(oldsid, sid string) (Store, error)
	Destroy(sid string) error
}

// internal usage
type StoreGCer interface {
	GC() error
}
