// Package xml provides functions for rendering
// XML objects.
package xml

import (
	"encoding/xml"
	"flag"
	"net/http"
)

var (
	indent = flag.Bool("xml:indent", false, "use a human readable format of XML")
)

// XML is a controller with helper functions
// for rendering Go objects as JSON.
type XML struct {
}

// RenderXML gets any object and returns an HTTP handler
// that renders the object.
func (c *XML) RenderXML(obj interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b []byte
		var err error
		if *indent {
			b, err = xml.MarshalIndent(obj, "", "\t")
		} else {
			b, err = xml.Marshal(obj)
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
