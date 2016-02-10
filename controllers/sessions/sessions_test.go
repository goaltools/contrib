package sessions

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Init(url.Values{})
	if s == nil {
		t.Errorf("Securecookie expected to be initialized, but it is empty.")
	}
}

func TestCookieExpiration(t *testing.T) {
	*cookieExpire = "2h"
	*cookieMaxAge = 100
	Init(url.Values{})

	r, _ := http.NewRequest("GET", "/", nil)

	c := Sessions{
		Request: r,
	}
	c.Before()

	expMax := time.Now().Add(*expireAfter)
	expMin := expMax.Add(-time.Second)
	cookie := c.cookie("")
	if u := cookie.Expires.Unix(); u < expMin.Unix() || u > expMax.Unix() {
		t.Errorf("Incorrect expiration %v. Expected a value between %v and %v.", u, expMin.Unix(), expMax.Unix())
	}

	if cookie.MaxAge != *cookieMaxAge {
		t.Errorf("Incorrect Max-Age parameter of a cookie. Expected %d, got %d.", *cookieMaxAge, cookie.MaxAge)
	}
}
