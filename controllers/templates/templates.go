// Package templates provides abstractions for work
// with standard Go template engine.
package templates

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/colegion/contrib/controllers/global"
)

var (
	layoutTpl   = flag.String("templates:layout.file", "Layout.html", "name of the layout template file")
	layoutBl    = flag.String("templates:layout.block", "layout", "name of the template block to render")
	elemTplPref = flag.String("templates:element.prefix", "_", "prefix of a template containing view elements")

	delimLeft  = flag.String("templates:delimLeft", "{%", "left action delimiter")
	delimRight = flag.String("templates:delimRight", "%}", "right action delimiter")

	views  = flag.String("templates:path", "./views/", "path to the directory with views")
	defTpl = flag.String("templates:default.pattern", "%v/%v.html", "default template's name pattern")

	contType = flag.String("templates:content.type", "text/html; charset=utf-8", "Content-Type header's value")

	// Funcs are added to the template's function map.
	// Functions are expected to return just 1 argument or
	// 2 in case the second one is of error type.
	Funcs = template.FuncMap{}

	// templates is a set of all registered and parsed templates.
	templates map[string]*template.Template
)

// Templates is a controller that provides support of HTML result
// rendering to your application.
// Use SetTemplatePaths to register templates and
// call c.RenderTemplate from your action to render some.
type Templates struct {
	// Global brings a Context that is used for passing variables to templates.
	*global.Global
}

// After renders the result template and writes it to the response.
func (c *Templates) After() http.Handler {
	// Check whether template rendering was requested.
	path := c.Context.Template()
	if path == "" {
		return nil
	}
	status := c.Context.Status()

	// If so, render the requested template.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set necessary headers of the response.
		w.Header().Set("Content-Type", *contType)

		// If required template does exist, execute it.
		err := fmt.Errorf(`template "%s" does not exist`, path)
		if tpl, ok := templates[path]; ok {
			w.WriteHeader(status)
			err = tpl.ExecuteTemplate(w, *layoutBl, c.Context)
			if err == nil { // Exit if execution has finished successfully.
				return
			}
		}

		// If the error above occured while rendering an error 500 page,
		// log the error and show an empty page to the user.
		if status == http.StatusInternalServerError {
			http.Error(w, err.Error(), status)
			return
		}

		// Otherwise, try to render 500 error page
		// by throwing a panic that must be caught by a router.
		panic(err.Error())
	})
}

// RenderTemplate is an action that gets a path to template
// and saves it to Context for later rendering.
func (c *Templates) RenderTemplate(templatePath string) http.Handler {
	c.Context.SetTemplate(templatePath)
	return nil
}

// Render is an equivalent of the following:
//	RenderTemplate(CurrentController + "/" + CurrentAction + ".html")
// The default path pattern may be overriden by adding the following
// line to your configuration file:
//	[templates]
//	default.pattern = %s/%s.tpl
func (c *Templates) Render() http.Handler {
	return c.RenderTemplate(
		fmt.Sprintf(*defTpl, c.CurrentController, c.CurrentAction),
	)
}

// Init triggers loading of templates.
func Init(url.Values) {
	templates = load(*views)
}
