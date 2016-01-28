// Package text is a controller for rendering
// plain/text.
package text

import (
	"fmt"
	"net/http"
)

// Text is a controller that provides helpers for
// rendering plain/text.
type Text struct {
}

// RenderText is a handler that works as fmt.Sprintf.
func (c *Text) RenderText(text string, args ...interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := fmt.Sprintf(text, args...)

		// TODO: use custom status code.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(t))
	})
}
