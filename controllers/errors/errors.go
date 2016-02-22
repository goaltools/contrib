package errors

import (
	"net/http"

	"github.com/colegion/contrib/controllers/templates"
)

// Errors is a controller that brings support of errors
// processing to your application.
type Errors struct {
	*templates.Templates
}

// NotFound is an action that renders a 404 page not found error.
//@route /404 404
func (c *Errors) NotFound() http.Handler {
	c.Context.SetStatus(http.StatusNotFound)
	return c.Render()
}

// MethodNotAllowed is an action that renders a 405 method not allowed error.
//@route /405 405
func (c *Errors) MethodNotAllowed() http.Handler {
	c.Context.SetStatus(http.StatusMethodNotAllowed)
	return c.Render()
}

// InternalServerError is an action that renders a 500 internal server error.
//@route /500 500
func (c *Errors) InternalServerError() http.Handler {
	c.Context.SetStatus(http.StatusInternalServerError)
	return c.Render()
}
