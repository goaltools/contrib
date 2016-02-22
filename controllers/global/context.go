package global

import (
	"flag"
	"net/http"
)

var (
	statusCodeKey = flag.String("context:status.code.key", "@sys.status.code", "a context key for status code")
)

// SetStatus sets a status code that must be returned with the response.
func (o Object) SetStatus(code int) {
	o[*statusCodeKey] = code
}

// Status returns a status code that must be returned with the response.
// If no codes were set, HTTP 200 OK status is returned.
func (o *Object) Status() int {
	i := o.Int(*statusCodeKey)
	if i == 0 {
		i = http.StatusOK
	}
	return i
}
