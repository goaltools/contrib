// Package grace is an abstraction around standard Go server
// and facebook's httpgrace.
// When on Windows a regular ListenAndServe is used.
// On other platforms graceful restarts and shutdowns are provided
// using gracehttp's Serve method.
package grace
