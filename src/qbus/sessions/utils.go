package sessions

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func NewSessionID() (sid string) {
	h := md5.New()
	str, _ := uuid.NewV4()
	h.Write(str.Bytes())

	sid = hex.EncodeToString(h.Sum(nil))
	return
}

func CacheKey(sid string) string {
	return fmt.Sprintf("sessions:%s", sid)
}

func EncodeSecureValue(raw string, secretKey string, createdAt time.Time) (value string, ok bool) {
	if raw == "" || secretKey == "" {
		return
	}

	timeValue := strconv.FormatInt(createdAt.UnixNano(), 10)

	h := hmac.New(sha1.New, []byte(secretKey))
	_, err := h.Write([]byte(raw + timeValue))
	if err != nil {
		return
	}

	src := url.QueryEscape(raw)
	src += COOKIE_VALUE_SPLIT + timeValue
	src += COOKIE_VALUE_SPLIT + hex.EncodeToString(h.Sum(nil))

	value = base64.URLEncoding.EncodeToString([]byte(src))
	ok = true
	return
}

func DecodeSecureValue(value string, secretKey string) (raw string, createdAt time.Time, ok bool) {
	rawBytes, _ := base64.URLEncoding.DecodeString(value)
	value = string(rawBytes)

	parts := strings.SplitN(value, COOKIE_VALUE_SPLIT, 3)
	if len(parts) < 3 {
		return
	}

	vRaw := strings.TrimSpace(parts[0])
	vCreated := strings.TrimSpace(parts[1])
	vHash := strings.TrimSpace(parts[2])

	if vRaw == "" || vCreated == "" || vHash == "" {
		return
	}

	vTime, _ := strconv.ParseInt(vCreated, 10, 64)
	if vTime <= 0 {
		return
	}

	vRaw, err := url.QueryUnescape(vRaw)
	if err != nil {
		return
	}

	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(vRaw + vCreated))

	if hex.EncodeToString(h.Sum(nil)) != vHash {
		return
	}

	raw = vRaw
	createdAt = time.Unix(0, vTime)
	ok = true
	return
}

func isExpired(createdAt time.Time, expired time.Duration) bool {
	return createdAt.Add(expired).Before(time.Now())
}

// get last same name cookie from cookies
// http://play.golang.org/p/LDfjMnJnhI
func getCookie(cookies []*http.Cookie, name string) (cookie *http.Cookie, ok bool) {
	for i := len(cookies) - 1; i >= 0; i-- {
		if cookies[i].Name == name {
			cookie = cookies[i]
			ok = true
			break
		}
	}

	return
}
