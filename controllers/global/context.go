package global

import (
	"flag"
	"net/http"
)

var (
	tplPathKey    = flag.String("context:template.path.key", "@sys.template.path", "a context key for template path")
	statusCodeKey = flag.String("context:status.code.key", "@sys.status.code", "a context key for status code")
)

// SetTemplate sets a path to the template that must be rendered.
func (o Object) SetTemplate(path string) {
	o[*tplPathKey] = path
}

// Template returns a path of the template that must be rendered.
// If not templates were set, empty string is returned.
func (o *Object) Template() string {
	return o.String(*tplPathKey)
}

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
