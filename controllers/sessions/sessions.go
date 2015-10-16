// Package sessions implements COOKIE based sessions.
package sessions

import (
	"flag"
	"net/http"

	"github.com/gorilla/securecookie"
)

var (
	cookieName   = flag.String("sessions:cookie.name", "_Session", "name of the cookie with session data")
	cookieDomain = flag.String("sessions:cookie.domain", "", "domain of cookie")
	cookieSecure = flag.Bool("sessions:cookie.secure", false, "")

	httpOnly  = flag.Bool("sessions:cookie.http.only", false, "")
	appSecret = flag.String("sessions:app.secret", string(securecookie.GenerateRandomKey(64)), "")

	hashKey []byte

	s *securecookie.SecureCookie
)

// Sessions is a controller that makes Session field
// available for your actions when you're using this
// controller as a parent.
type Sessions struct {
	Session map[string]string
}

// Before is a magic action of Sessions controller.
func (c *Sessions) Before() http.Handler {
	return nil
}

// Initially is a magic method that gets session info from a request
// and initializes Session field.
func (c *Sessions) Initially(w http.ResponseWriter, r *http.Request, a []string) bool {
	c.Session = map[string]string{}
	if cookie, err := r.Cookie(*cookieName); err == nil {
		s.Decode(*cookieName, cookie.Value, &c.Session)
	}
	return false
}

// Finally is a magic method that will be executed at the very end of request
// life cycle and is responsible for creating a signed cookie with session info.
func (c *Sessions) Finally(w http.ResponseWriter, r *http.Request, a []string) bool {
	if encoded, err := s.Encode(*cookieName, c.Session); err == nil {
		cookie := &http.Cookie{
			Name:     *cookieName,
			Value:    encoded,
			Domain:   *cookieDomain,
			HttpOnly: *httpOnly,
			Path:     "/",
			Secure:   *cookieSecure,
		}
		http.SetCookie(w, cookie)
	}
	return false
}

// Init is a function that is used for initialization of
// Sessions controller.
func Init() {
	hashKey = []byte(*appSecret)
	s = securecookie.New(hashKey, nil)
}
