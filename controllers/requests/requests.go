package requests

import (
	"flag"
	"net/http"
)

var (
	// maxMem is a value per every upload (i.e. there are may be thousands of users who upload
	// files 32MB of which are stored in memory). The remainder (out of 32MB) is stored
	// on disk in temporary files.
	maxMem = flag.Int64("requests:max.memory", 32, "number of MB to store in memory when parsing a file")
)

// Requests is a controller that does two things:
// 1. Calls Request.ParseForm to parse GET / POST requests;
// 2. Makes Request available in your controller (use c.Request).
type Requests struct {
	Request *http.Request `bind:"request"`
}

// Before calls ParseForm of the c.Request.
// At the same time, if used with a standard goal routing package,
// parameters extracted from URN are saved to the Form field of the Request.
func (c *Requests) Before() http.Handler {
	// Save the old value of Form, "github.com/colegion/contrib/routers/denco"
	// uses it to pass parameters extracted from URN.
	t := c.Request.Form

	// Set r.Form to nil, otherwise ParseForm / ParseMultipartForm will not work.
	c.Request.Form = nil

	// Parse the body depending on the Content-Type.
	var err error
	switch c.Request.Header.Get("Content-Type") {
	case "multipart/form-data":
		err = c.Request.ParseMultipartForm(*maxMem << 20)
	default:
		err = c.Request.ParseForm()
	}

	// Make sure the parsing was successful.
	// Otherwise, return a "bad request" error.
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		})
	}

	// Add the old values from router to the new r.Form.
	// Copying only one value per key as the router does not pass more than that.
	for k := range t {
		c.Request.Form.Add(k, t.Get(k))
	}
	return nil
}
