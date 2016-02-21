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
