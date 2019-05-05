package sessions

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Manager delegates to provider of session store
type Manager struct {
	provider Provider
}

// NewManager returns a new manager with provider
func NewManager(provider Provider) *Manager {
	// SecretKey MUST be present or panic!
	if provider.Config().SecretKey == "" {
		panic(ErrEmptySecretKey)
	}

	return &Manager{
		provider: provider,
	}
}

// Start restores a session from request cookie or creates a new one
func (mgr *Manager) Start(w http.ResponseWriter, r *http.Request) (sess Store, createdAt time.Time, err error) {
	sid, createdAt, ok := mgr.parseSessionID(r)
	if ok {
		sess, err = mgr.provider.Restore(sid)
		if err == nil || err != ErrNotFound {
			return
		}
	}

	sess, err = mgr.createSession()
	if err == nil {
		createdAt = mgr.writeSessionCookie(r, w, sess.ID())
	}

	return
}

// Refresh updates request session with new ID
func (mgr *Manager) Refresh(w http.ResponseWriter, r *http.Request, params ...map[string]interface{}) (sess Store, err error) {
	oldsid, _, ok := mgr.parseSessionID(r)
	if !ok {
		sess, err = mgr.createSession(params...)
		if err != nil {
			return
		}

		mgr.writeSessionCookie(r, w, sess.ID())
		return
	}

	var (
		sid     string
		retried int
	)

retry:
	sid = NewSessionID()
	sess, err = mgr.provider.Refresh(oldsid, sid)

	switch err {
	case ErrNotFound:
		sess, err = mgr.createSession(params...)
		if err != nil {
			return
		}

	case ErrDuplicateID:
		retried++
		if retried >= 3 {
			return
		}

		goto retry
	}

	mgr.writeSessionCookie(r, w, sess.ID())
	return
}

// Destroy deletes session of current request
func (mgr *Manager) Destroy(w http.ResponseWriter, r *http.Request) error {
	config := mgr.provider.Config()

	sid, _, ok := mgr.parseSessionID(r)
	if !ok {
		return nil
	}

	// ignore provider error is ok!
	mgr.provider.Destroy(sid)

	// force expire client cookie
	cookie := &http.Cookie{
		Name:     config.CookieName,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.CookieSecure,
		Expires:  time.Now(),
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)

	return nil
}

// GC cleans up expired session data by provider implements
func (mgr *Manager) GC(intervals ...time.Duration) {
	gc, ok := mgr.provider.(StoreGCer)
	if !ok {
		return
	}

	var interval time.Duration
	if len(intervals) > 0 {
		interval = intervals[0]

		// defaults to an hour
		if interval < time.Minute {
			interval = time.Hour
		}
	}

	// try gc first
	err := gc.GC()
	if err != nil {
		log.Printf("[SESSION] manager.GC(%v): %v\n", interval, err)
	}

	// start an interval events
	var timer *time.Timer
	timer = time.AfterFunc(interval, func() {
		timer.Stop()

		mgr.GC(interval)
	})
}

// writeSessionCookie sets cookie with session data
func (mgr *Manager) writeSessionCookie(r *http.Request, w http.ResponseWriter, sid string) (createdAt time.Time) {
	config := mgr.provider.Config()

	createdAt = time.Now()

	// secure cookie value of sid
	value, ok := EncodeSecureValue(sid, config.SecretKey, createdAt)
	if !ok {
		return
	}

	cookie := &http.Cookie{
		Name:     config.CookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.CookieSecure,
	}

	if config.CookieExpire >= 0 {
		cookie.MaxAge = config.CookieExpire
	}

	http.SetCookie(w, cookie)
	r.AddCookie(cookie)

	return
}

func (mgr *Manager) createSession(params ...map[string]interface{}) (sess Store, err error) {
	var (
		sid     string
		retried int
	)

retry:
	sid = NewSessionID()

	sess, err = mgr.provider.Create(sid, params...)
	if err == ErrDuplicateID {
		retried++
		if retried >= 3 {
			return
		}

		goto retry
	}

	return
}

func (mgr *Manager) parseSessionID(r *http.Request) (sid string, createdAt time.Time, ok bool) {
	config := mgr.provider.Config()

	cookie, ok := getCookie(r.Cookies(), config.CookieName)
	if !ok {
		return
	}

	value, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return
	}

	value, vTime, ok := DecodeSecureValue(value, config.SecretKey)
	if !ok {
		return
	}

	value = strings.TrimSpace(value)
	if len(value) != SESSION_ID_LENGTH {
		value = ""
		return
	}

	sid = value
	createdAt = vTime
	ok = true
	return
}
