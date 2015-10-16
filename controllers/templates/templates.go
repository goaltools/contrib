// Package templates provides abstractions for work
// with standard Go template engine.
package templates

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/sasbury/mini"
)

var (
	baseTemplate   = flag.String("templates:base", "Base.html", "name of the main template")
	renderTemplate = flag.String("templates:render", "base", "name of template block to render")
	delimLeft      = flag.String("templates:delimLeft", "{%", "left action delimiter")
	delimRight     = flag.String("templates:delimRight", "%}", "right action delimiter")
	listF          = flag.String("templates:views", "assets/views/views.ini", "file with a list of views")

	root  = flag.String("root.directory", "./", "path where files of the project are stored")
	views = flag.String("views.directory", "views/", "path with the views, relative to the root")

	// Funcs are added to the template's function map.
	// Functions are expected to return just 1 argument or
	// 2 in case the second one is of error type.
	Funcs template.FuncMap
)

// Templates is a controller that provides support of HTML result
// rendering to your application.
// Use SetTemplatePaths to register templates and
// call c.RenderTemplate from your action to render some.
type Templates struct {
	// Context is used for passing variables to templates.
	Context map[string]interface{}

	// Status is a status code that will be returned when rendering.
	Status int

	defaultTemplate string
}

// Initially sets default template name as CurrentController + CurrentAction + .html.
// Third argument is garanteed to contain Controller as a 0th argument
// and Action as a 1st.
func (c *Templates) Initially(w http.ResponseWriter, r *http.Request, a []string) bool {
	c.defaultTemplate = fmt.Sprintf("%s/%s.html", a[0], a[1])
	return false
}

// Before initializes Context that will be passed to template.
func (c *Templates) Before() http.Handler {
	c.Context = map[string]interface{}{}
	return nil
}

// Render is an equivalent of
// RenderTemplate(ControllerName+"/"+ActionName+".html").
func (c *Templates) Render() http.Handler {
	return &Handler{
		context:  c.Context,
		template: c.defaultTemplate,
	}
}

// RenderTemplate is an action that gets a path to template
// and renders it using data from Context.
func (c *Templates) RenderTemplate(templatePath string) http.Handler {
	return &Handler{
		context:  c.Context,
		template: templatePath,
	}
}

// RenderError is an action that renders Error 500 page.
func (c *Templates) RenderError() http.Handler {
	c.Status = http.StatusInternalServerError
	return c.RenderTemplate("Errors/InternalError.html")
}

// RenderNotFound is an action that renders Error 404 page.
func (c *Templates) RenderNotFound() http.Handler {
	c.Status = http.StatusNotFound
	return c.RenderTemplate("Errors/NotFound.html")
}

// Init initializes triggers loading of templates.
func Init() {
	load(*root, *views, loadViewsList())
}

func loadViewsList() map[string]string {
	// Open the configuration file with a list of views.
	c, err := mini.LoadConfiguration(*listF)
	if err != nil {
		log.Fatal(err)
	}

	// Get all the keys that are found there.
	ks := c.KeysForSection("views")

	// Generate a list of file names.
	m := map[string]string{}
	for i := range ks {
		m[ks[i]] = c.StringFromSection("views", ks[i], "")
	}
	return m
}
