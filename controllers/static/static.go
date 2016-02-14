package static

import (
	"flag"
	"net/http"
)

var (
	path   = flag.String("static:root.directory", "./static", "path to the directory with static assets")
	prefix = flag.String("static:path.prefix", "/", "a prefix that's added to static assets' paths")
)

// Static is a controller that brings static
// assets' serving functionality to your app.
type Static struct {
}

// Serve is a wrapper around Go's standard FileServer
// and StripPrefix HTTP handlers.
//@get /*filepath
func (c *Static) Serve(filepath string) http.Handler {
	return http.StripPrefix(*prefix, http.FileServer(http.Dir(*path)))
}
