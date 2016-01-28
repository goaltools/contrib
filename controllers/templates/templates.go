// Package templates provides abstractions for work
// with standard Go template engine.
package templates

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	layoutTpl   = flag.String("templates:layout.file", "Layout.html", "name of the layout template file")
	layoutBl    = flag.String("templates:layout.block", "layout", "name of the template block to render")
	elemTplPref = flag.String("templates:element.prefix", "_", "prefix of a template containing view elements")

	delimLeft  = flag.String("templates:delimLeft", "{%", "left action delimiter")
	delimRight = flag.String("templates:delimRight", "%}", "right action delimiter")

	views   = flag.String("templates:path", "./views/", "path to the directory with views")
	defTpl  = flag.String("templates:default.pattern", "%v/%v.html", "default template's name pattern")
	errsDir = flag.String("templates:errors.dir", "Errors", "a directory with error templates")

	devMode  = flag.Bool("mode.dev", false, "development mode with debugging enabled")
	contType = flag.String("templates:content.type", "text/html; charset=utf-8", "Content-Type header's value")

	// Funcs are added to the template's function map.
	// Functions are expected to return just 1 argument or
	// 2 in case the second one is of error type.
	Funcs = template.FuncMap{}

	// Log is a default logger used by the templates controller.
	Log = log.New(os.Stderr, "Templates: ", log.LstdFlags)

	templates map[string]*template.Template
)

// Templates is a controller that provides support of HTML result
// rendering to your application.
// Use SetTemplatePaths to register templates and
// call c.RenderTemplate from your action to render some.
type Templates struct {
	// Context is used for passing variables to templates.
	Context map[string]interface{}

	// StatusCode is a status code that will be returned when rendering.
	// If not specified explicitly, 200 will be used.
	StatusCode int

	defTpl string
}

// Initially sets a name of the template that should be rendered by default
// (i.e. if no templates are defined explicitly). It will looks as follows:
//	CurrentController + / + CurrentAction + .html
// Argument a is guaranteed to contain at least 2 arguments: Controller's name
// as 0th and Action's name as 1st.
// It also allocates and initializes Context.
func (c *Templates) Initially(w http.ResponseWriter, r *http.Request, a []string) bool {
	// Set the default template to render.
	c.defTpl = fmt.Sprintf(*defTpl, a[0], a[1])

	// Allocate a new context.
	c.Context = map[string]interface{}{}
	return false
}

// RenderTemplate is an action that gets a path to template
// and renders it using data from Context.
func (c *Templates) RenderTemplate(templatePath string) http.Handler {
	return &Handler{
		context:  c.Context,
		status:   c.StatusCode,
		template: templatePath,
	}
}

// Render is an equivalent of the following:
//	RenderTemplate(CurrentController + "/" + CurrentAction + ".html")
// The default path pattern may be overriden by adding the following
// line to your configuration file:
//	[templates]
//	default.pattern = %s/%s.tpl
func (c *Templates) Render() http.Handler {
	return c.RenderTemplate(c.defTpl)
}

// RenderError is an action that gets an error and returns
// an Internal Error 500 handler that will render appropriate
// template. In development mode it will also include an error message
// in the template. Template path to be rendered:
//	./errors/500.html
// A way to redefine the path is to update you configuration file:
//	[templates]
//	default.pattern = %s/%s.html
//	errors.dir = errors
func (c *Templates) RenderError(err error) http.Handler {
	if *devMode {
		c.Context["error"] = err
	}
	c.StatusCode = http.StatusInternalServerError
	return c.renderError()
}

// RenderNotFound is similar to RenderError but prints Error 404
// and gets optional messages as input arguments rather than an error.
func (c *Templates) RenderNotFound(msgs ...string) http.Handler {
	c.Context["messages"] = msgs
	c.StatusCode = http.StatusNotFound
	return c.renderError()
}

// renderError renders an "errors/StatusCode.html" template.
func (c *Templates) renderError() http.Handler {
	return c.RenderTemplate(fmt.Sprintf(*defTpl, *errsDir, c.StatusCode))
}

// Redirect gets a URI or URN (e.g. "https://si.te/smt or "/users")
// and returns a handler for user's redirect using 303 status code.
func (c *Templates) Redirect(urn string) http.Handler {
	return http.RedirectHandler(urn, http.StatusSeeOther)
}

// Init triggers loading of templates.
func Init(_ url.Values) {
	templates = load(*views)
}
