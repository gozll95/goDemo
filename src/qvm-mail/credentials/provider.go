package credential

// A Provider is the interface for any component which will provide credential
// Value. A provider is required to manage its own Expired state, and what to
// be expired means.
//
// The Provider should not need to implement its own mutexes, because
// that will be managed by Credential.
type Provider interface {
	Parse() (Value, error)
	IsExpired() bool
}

// StaticProvider is a set of credential which are set pragmatically, and
// will never expire.
type StaticProvider struct {
	Value
}

// NewStaticCredentials returns a pointer to a new Credential object
// wrapping a static credential value provider.
func NewStaticProvider(id, secret, token string) *Credential {
	return NewCredential(&StaticProvider{
		Value: Value{
			AccessKeyId:     id,
			AccessKeySecret: secret,
			SessionToken:    token,
		},
	})
}

// Parse returns the credential or error if the credential are invalid.
func (s *StaticProvider) Parse() (Value, error) {
	if s.AccessKeyId == "" || s.AccessKeySecret == "" {
		return Value{}, ErrCredentialEmpty
	}

	return s.Value, nil
}

// IsExpired returns if the credentials are expired.
//
// For StaticProvider, the credentials never expired.
func (s *StaticProvider) IsExpired() bool {
	return false
}
