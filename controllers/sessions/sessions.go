// Package sessions implements COOKIE based sessions.
package sessions

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/securecookie"
)

var (
	cookieName   = flag.String("sessions:cookie.name", "_Session", "name of the cookie with session data")
	cookieDomain = flag.String("sessions:cookie.domain", "", "domain of cookie")
	cookieSecure = flag.Bool("sessions:cookie.secure", false, "prevent transmission of a cookie in clear text")

	cookieMaxAge = flag.Int("sessions:cookie.maxage", 0, "time in seconds for when a cookie will be deleted")
	cookieExpire = flag.String("sessions:cookie.expires.duration", "", "a time duration when a cookie expires")

	httpOnly  = flag.Bool("sessions:cookie.http.only", false, "")
	appSecret = flag.String("sessions:app.secret", string(securecookie.GenerateRandomKey(64)), "")

	hashKey []byte

	s *securecookie.SecureCookie

	expireAfter *time.Duration
)

// Sessions is a controller that makes Session field
// available for your actions when you're using this
// controller as a parent.
type Sessions struct {
	Session map[string]string

	Request  *http.Request       `bind:"request"`
	Response http.ResponseWriter `bind:"response"`
}

// Before is a magic action that gets session info from a request
// and initializes Session field.
func (c *Sessions) Before() http.Handler {
	c.Session = map[string]string{}
	if cookie, err := c.Request.Cookie(*cookieName); err == nil {
		s.Decode(*cookieName, cookie.Value, &c.Session)
	}
	return nil
}

// After is a magic action that will be executed at the very end of request
// life cycle and is responsible for creating a signed cookie with session info.
func (c *Sessions) After() http.Handler {
	if encoded, err := s.Encode(*cookieName, c.Session); err == nil {
		cookie := c.cookie(encoded)
		http.SetCookie(c.Response, cookie)
	}
	return nil
}

func (c *Sessions) cookie(data string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     *cookieName,
		Value:    data,
		Domain:   *cookieDomain,
		HttpOnly: *httpOnly,
		Path:     "/",
		Secure:   *cookieSecure,
		MaxAge:   *cookieMaxAge,
	}
	if expireAfter != nil {
		cookie.Expires = time.Now().Local().Add(*expireAfter)
	}
	return cookie
}

// Init is a function that is used for initialization of
// Sessions controller.
func Init(url.Values) {
	hashKey = []byte(*appSecret)
	s = securecookie.New(hashKey, nil)

	// Convert "expiration" string to a time duration
	// if it is not empty.
	if *cookieExpire != "" {
		e, err := time.ParseDuration(*cookieExpire)
		if err != nil {
			log.Panicf(`Cannot parse expire time duration. Error: %v.`, err)
		}
		expireAfter = &e
		log.Printf(`Parameter "cookie.expires.duration" is equal to %v.`, expireAfter)
	}

	// Print a message if MaxAge is set.
	if *cookieMaxAge != 0 {
		log.Printf(`Parameter "cookie.maxage" is equal to %vs.`, *cookieMaxAge)
	}
}
