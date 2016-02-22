// Package standard is a wrapper around standard Go template system.
// It is however brings some conventions. For illustration, there is a sample listing:
//	+ Layout.html
//	+ _header.html
//	- Profiles/
//		+ Index.html
//		+ Layout.html
//		+ _header.html
//	- Errors/
//		+ NotFound.html
//		+ _header.html
//		+ _footer.html
// The templates above when parsing will be grouped as follows:
//	* Profiles/Index.html, Profiles/Layout.html, _header.html, Profiles/_header.html
//	* Errors/NotFound.html, Layout.html, _header.html, Errors/_header.html, Errors/_footer.html
// I.e. every template (e.g. "Profiles/Index.html) is always parsed with the closest Layout
// file (e.g. "Profiles/Layout.html"). Moreover, files started with "_" are parsed with
// every of the templates in the same directory or its subdirectories.
//
// Layout templates must have "layout" block that will be executed when rendering a template.
//
// Names of the layout and element ("_smth.html") files, layout block name, and other
// paremeters may be customized by update of your configuration file (if "iniflag" is used)
// and/or by use of flags when running your app.
package standard

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"

	"github.com/colegion/contrib/controllers/global"
)

var (
	layout    = flag.String("templates:layout.file", "Layout.html", "name of the layout template file")
	mainBlock = flag.String("templates:layout.block", "main", "name of the template block to render")
	elemPref  = flag.String("templates:element.file.prefix", "_", "prefix of a template containing view elements")

	delimLeft  = flag.String("templates:delimLeft", "{%", "left template action delimiter")
	delimRight = flag.String("templates:delimRight", "%}", "right template action delimiter")

	views         = flag.String("templates:path", "./views/", "path to the directory with views")
	defTplPattern = flag.String("templates:default.pattern", "%v/%v.html", "default template name's pattern")

	contType = flag.String("templates:content.type", "text/html; charset=utf-8", "Content-Type header's value")
)

// Template implements (github.com/colegion/contrib/controllers/templates).Template interface.
type Template struct {
}

// Funcs gets a map of template functions and sets them for use
// by Init function. This must be called before Init (e.g. from your init function)
// in order to make effect.
func (t Template) Funcs(fm template.FuncMap) {
	funcs = fm
}

// RenderTemplate gets a context, template name and returns a handler that may be
// used for its rendering. If the requested template does not exist,
// an error handler is returned.
func (t Template) RenderTemplate(context global.Object, name string) http.Handler {
	status := context.Status()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set necessary headers of the response.
		w.Header().Set("Content-Type", *contType)

		// If required template does exist, render it.
		err := fmt.Errorf(`template "%s" does not exist`, name)
		if tpl, ok := templates[name]; ok {
			w.WriteHeader(status)
			err = tpl.ExecuteTemplate(w, *mainBlock, context)
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

// Render transforms controller and action to a template name using a pattern
// defined by "default.pattern" flag. For example, it may look as follows:
//	fmt.Sprintf("%s/%s.html", controller, action)
// It then calls RenderTemplate.
func (t Template) Render(context global.Object, controller, action string) http.Handler {
	return t.RenderTemplate(context, fmt.Sprintf(*defTplPattern, controller, action))
}

// Init starts parsing of templates.
func (t Template) Init() {
	templates = load(*views)
}

// templates is a set of all registered and parsed templates.
var templates map[string]*template.Template
