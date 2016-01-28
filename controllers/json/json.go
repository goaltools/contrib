// Package json provides functions for rendering
// JSON objects.
package json

import (
	"encoding/json"
	"flag"
	"net/http"
)

var (
	indent = flag.Bool("json:indent", false, "use a human readable format of JSON")
)

// JSON is a controller with helper functions
// for rendering Go objects as JSON.
type JSON struct {
}

// RenderJSON gets any object and returns an HTTP handler
// that renders the object.
func (c *JSON) RenderJSON(obj interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b []byte
		var err error
		if *indent {
			b, err = json.MarshalIndent(obj, "", "\t")
		} else {
			b, err = json.Marshal(obj)
		}

		// Make sure there are no errors.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: use custom status code.
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
}
