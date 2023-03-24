package utils

import (
	"net/http"
	"net/url"
	"time"
)

func GetCookie(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	accessToken, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func SetCookie(w http.ResponseWriter, cookieName string, cookieValue string, expiration time.Time) {
	cookie := http.Cookie{Name: cookieName, Value: cookieValue, Expires: expiration}
	http.SetCookie(w, &cookie)
}
