package global

import (
	"net/http"
)

// Object is a type that is used for global variables.
type Object map[string]interface{}

// Global is a controller that provides a registry for
// global request variables.
type Global struct {
	// Context is used for passing variables to templates.
	Context Object

	// CurrentAction and CurrentController will be automatically initialized
	// with the values of current controller / action at the time of
	// this controller's allocation.
	CurrentAction     string `bind:"action"`
	CurrentController string `bind:"controller"`
}

// Before action allocates Context.
func (c *Global) Before() http.Handler {
	c.Context = Object{}
	return nil
}

// Int returns a value associated with the key as an integer.
func (o *Object) Int(key string) int {
	v := o.get(key)
	if v == nil {
		return 0
	}
	i, ok := v.(int)
	if !ok {
		return 0
	}
	return i
}

// String returns a value associated with the key as a string.
func (o *Object) String(key string) string {
	v := o.get(key)
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

// get returns a value associated with the key.
func (o Object) get(key string) interface{} {
	if v, ok := o[key]; ok {
		return v
	}
	return nil
}
