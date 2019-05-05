package sessions

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (mgr *Manager) HasRemember(r *http.Request) (psk, salt string, ok bool) {
	config := mgr.provider.Config()

	cookie, tmpok := getCookie(r.Cookies(), config.CookieRememberName)
	if !tmpok {
		return
	}

	value, createdAt, tmpok := DecodeSecureValue(cookie.Value, config.SecretKey)
	if !tmpok || isExpired(createdAt, config.RememberExpireSeconds()) {
		return
	}

	parts := strings.SplitN(value, COOKIE_VALUE_PARTS_SPLIT, 2)
	if len(parts) != 2 {
		return
	}

	tmppsk, err := url.QueryUnescape(parts[0])
	if err != nil {
		return
	}

	tmpsalt, err := url.QueryUnescape(parts[1])
	if err != nil {
		return
	}

	psk = tmppsk
	salt = tmpsalt
	ok = true
	return
}

func (mgr *Manager) IsValidRemember(psk, salt, saltHash string) (ok bool) {
	if psk == "" || salt == "" || saltHash == "" {
		return
	}

	config := mgr.provider.Config()

	tmppsk, createdAt, tmpok := DecodeSecureValue(saltHash, config.SecretKey+salt)
	if !tmpok || isExpired(createdAt, config.RememberExpireSeconds()) {
		return
	}

	return tmppsk == psk
}

func (mgr *Manager) WriteRememberCookie(w http.ResponseWriter, psk string, salt string) (ok bool) {
	config := mgr.provider.Config()

	createdAt := time.Now()

	value, ok := createRememberCookieValue(config.SecretKey, psk, salt, createdAt)
	if !ok {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     config.CookieRememberName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.CookieSecure,
		MaxAge:   config.RememberExpire,
	})
	return
}

func (mgr *Manager) DestroyRememberCookie(w http.ResponseWriter) {
	config := mgr.provider.Config()

	http.SetCookie(w, &http.Cookie{
		Name:     config.CookieRememberName,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.CookieSecure,
		MaxAge:   -1,
	})
}

func createRememberCookieValue(sk, psk, salt string, createdAt time.Time) (hashValue string, ok bool) {
	if sk == "" || psk == "" || salt == "" {
		return
	}

	saltHash, ok := EncodeSecureValue(psk, sk+salt, createdAt)
	if !ok {
		return
	}

	cookieValue := url.QueryEscape(psk) + COOKIE_VALUE_PARTS_SPLIT + url.QueryEscape(saltHash)
	hashValue, ok = EncodeSecureValue(cookieValue, sk, createdAt)
	return
}
