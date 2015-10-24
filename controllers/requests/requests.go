package requests

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	// maxMem is a value per every upload (i.e. there are may be thousands of users who upload
	// files 32MB of which are stored in memory). The remainder (out of 32MB) is stored
	// on disk in temporary files.
	maxMem = flag.Int64("requests:max.memory", 32, "number of MB to store in memory when parsing a file")

	// Log is a default logger of the controller.
	Log = log.New(os.Stderr, "Requests Controller: ", log.LstdFlags)
)

// Requests is a controller that does two things:
// 1. Calls Request.ParseForm to parse GET / POST requests;
// 2. Makes Request available in your controller (use c.Request).
type Requests struct {
	Request *http.Request
}

// Initially calls ParseForm on the request and saves it to c.Request.
// At the same time, if used with a standard goal routing package,
// parameters extracted from URN are saved to the Form field of the Request.
func (c *Requests) Initially(w http.ResponseWriter, r *http.Request, as []string) bool {
	// Save the old value of Form, "github.com/colegion/contrib/routers/denco"
	// uses it to pass parameters extracted from URN.
	t := r.Form

	// Set r.Form to nil, otherwise ParseForm / ParseMultipartForm will not work.
	r.Form = nil

	// Parse the body depending on the Content-Type.
	var err error
	switch r.Header.Get("Content-Type") {
	case "multipart/form-data":
		err = r.ParseMultipartForm(*maxMem << 20)
	case "application/x-www-form-urlencoded":
		err = r.ParseForm()
	}

	// Make sure the parsing was successfull.
	if err != nil {
		go Log.Printf("Failed to parse request body: %v.", err)
		return true
	}

	// Add the old values from router to the new r.Form.
	// Copying only one value per key as the router does not pass more than that.
	for k := range t {
		r.Form.Add(k, t.Get(k))
	}

	// Save the request, so it can be accessed from child controllers.
	c.Request = r
	return false
}
