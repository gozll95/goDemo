package credential

import (
	"sync"
)

// NOTE: It's modified from https://github.com/aws/aws-sdk-go/tree/master/aws/credentials

type Credential struct {
	mux          sync.Mutex
	value        Value
	provider     Provider
	forceRefresh bool
}

func NewCredential(provider Provider) *Credential {
	return &Credential{
		provider:     provider,
		forceRefresh: true,
	}
}

// Get returns the Value, or error if the Value failed
// to be retrieved.
//
// It will return the cached Value if it has not expired. If the
// Value has expired the Provider's Parse() will be called
// to refresh the credential.
//
// If Credential.Expire() was called the Value will be force
// expired, and the next call to Find() will cause them to be refreshed.
func (c *Credential) Get() (Value, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.isExpired() {
		value, err := c.provider.Parse()
		if err != nil {
			return Value{}, err
		}

		c.value = value
		c.forceRefresh = false
	}

	return c.value, nil
}

// Expire expires the credential and forces to be refreshed on the
// next call to Find().
//
// This will override the Provider's expired state, and force Credential
// to call the Provider's Parse().
func (c *Credential) Expire() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.forceRefresh = true
}

// IsExpired returns if the credential are no longer valid, and need
// to be parsed.
//
// If the Credential were forced to be expired with Expire() this will
// reflect that override.
func (c *Credential) IsExpired() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.isExpired()
}

func (c *Credential) isExpired() bool {
	return c.forceRefresh || c.provider.IsExpired()
}
