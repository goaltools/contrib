package templates

import (
	"html/template"
	"net/http"
	"net/url"

	"github.com/colegion/contrib/controllers/global"
	"github.com/colegion/contrib/controllers/requests"
	"github.com/colegion/contrib/controllers/templates/standard"
)

// TemplateSystem is a template engine that is used by this controller.
// By default it is a wrapper around standard templates + some conventions.
// This may be replaced by a template system of your choice:
//	func init() {
//		templates.TemplateSystem = myTSThatImplementTemplateInterface
//	}
// If you want to set template functions, do the following:
//	templates.TemplateSystem.Funcs(template.FuncMap{
//		...
//	})
var TemplateSystem = Template(standard.Template{})

// Template is an interface that must be implemented by template systems
// in order to be compatible with Templates controller.
type Template interface {
	// Funcs gets functions that will be available for use in the templates. The functions
	// must return just 1 argument or 2 in case the second one is of error type.
	Funcs(template.FuncMap)

	// Init is a function that must initialize the template system,
	// load and parse templates, etc.
	Init()

	// RenderTemplate gets a context, template name and must return a
	// handler that would execute the template or error 500.
	RenderTemplate(context global.Object, name string) http.Handler

	// Render is an equivalent of RenderTemplate except instead of
	// template name it gets current controller and action names
	// and renders a template associated with them.
	Render(context global.Object, controller, action string) http.Handler
}

// Templates is a controller that provides support of HTML result
// rendering to your application.
// Use SetTemplatePaths to register templates and
// call c.RenderTemplate from your action to render some.
type Templates struct {
	// Global brings a Context that is used for passing variables to templates.
	*global.Global

	*requests.Requests
}

// RenderTemplate is an action that gets a path to template
// and returns an HTTP handler for its rendering.
func (c *Templates) RenderTemplate(templatePath string) http.Handler {
	return TemplateSystem.RenderTemplate(c.Context, templatePath)
}

// Render executes a te
func (c *Templates) Render() http.Handler {
	return TemplateSystem.Render(c.Context, c.CurrentController, c.CurrentAction)
}

// Init triggers loading of templates.
func Init(url.Values) {
	TemplateSystem.Init()
}
