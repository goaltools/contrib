// When not on Windows use graceful restarts and shutdowns.
//
// +build !windows

package grace

import (
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
)

// Serve is a wrapper around gracehttp.Serve. As opposed
// to the standard net/http package, gracehttp server may be terminated
// and/or restarted without dropping any connections.
func Serve(s *http.Server) error {
	return gracehttp.Serve(s)
}
