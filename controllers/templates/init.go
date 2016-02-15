package templates

import (
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// load gets a path with views and returns a map of
// parsed templates that can be used by Render actions.
func load(dir string) map[string]*template.Template {
	m := map[string]*template.Template{}

	dir = filepath.Clean(dir)
	log.Printf(`Parsing templates in "%s".`, dir)

	tpls := map[string]string{}
	els := []string{}
	ls := layouts{}
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		// Make sure there are no any errors.
		if err != nil {
			log.Printf("Failed to traverse template files. Error: %v.", err)
			return err
		}

		// Ignore directories.
		if info.IsDir() {
			return nil
		}

		// Get current path without the dir at the beginning.
		// E.g. if we're scanning "views/app/index.html", extract
		// the "app/index.html" part only.
		rel, _ := filepath.Rel(dir, p)
		relNorm := filepath.ToSlash(rel)

		// Check whether current file (e.g. "index.html") is a layout template.
		b := path.Base(relNorm)
		if b == *layoutTpl {
			ls[path.Dir(relNorm)] = true
			return nil
		}

		// Check whether current file is a view element (e.g. "element_button.html").
		if strings.HasPrefix(b, *elemTplPref) {
			els = append(els, p)
			return nil
		}

		// Otherwise, just add it to the list of templates.
		tpls[relNorm] = p
		return nil
	})

	// Print a list of element templates.
	log.Printf("View elements: %v", els)

	// Parse templates and register them.
	var err error
	for relNorm, p := range tpls {
		t := template.New(relNorm).Funcs(Funcs).Delims(*delimLeft, *delimRight)

		// Check whether current template must have
		// a layout file.
		if l, ok := ls.path(path.Dir(relNorm)); ok {
			m[relNorm], err = t.ParseFiles(append(els, l, p)...)
			log.Printf("\t%s (%s)", p, l)
		} else {
			m[relNorm], err = t.ParseFiles(append(els, p)...)
			log.Printf("\t%s", p)
		}

		// Make sure there were no errors during parsing.
		if err != nil {
			log.Panicf(`Failed to parse "%s". Error: %v.`, p, err)
		}
	}
	return m
}

// layouts stores information about directories that have layout
// templates. E.g., if there is "layout.html" in "app/profiles", there will be:
//	app/profiles: true
type layouts map[string]bool

// path gets a template dir's path and returns associated layout template's path. E.g.:
//	- views/
//		+ layout.html
//		- app/
//			+ index.html
//			- profiles/
//				+ index.html
//				+ layout.html
// In case of the listing that's provided above, path must return:
//	layout.html
// if "app/" is provided as argument. And:
//	app/profiles/layout.html
// if "app/profiles/" is provided.
func (l layouts) path(dir string) (string, bool) {
	// If the requested directory has a layout file, return its path.
	if l[dir] {
		return filepath.Join(*views, filepath.FromSlash(dir), *layoutTpl), true
	}

	// If it's not and this is a root directory, return.
	if dir == "." {
		return "", false
	}

	// Otherwise, check higher level directories.
	return l.path(path.Join(dir, ".."))
}
