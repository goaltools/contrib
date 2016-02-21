package standard

import (
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	funcs template.FuncMap
)

// load gets a path with views and returns a map of
// parsed templates that can be used by render actions.
func load(dir string) map[string]*template.Template {
	m := map[string]*template.Template{}

	dir = filepath.Clean(dir) // Transform the input path to OS specific format.
	log.Printf(`Parsing templates in "%s".`, dir)

	tpls := map[string]string{}
	es := elements{}
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
		// E.g. if we're scanning "views/App/Index.html", extract
		// the "App/index.html" part only.
		rel, _ := filepath.Rel(dir, p)
		relNorm := filepath.ToSlash(rel) // Transform the path to normalized (i.e. *nix) form.

		// Check whether current file (e.g. "Index.html") is a layout template.
		b := path.Base(relNorm)
		if b == *layout {
			ls[path.Dir(relNorm)] = true
			return nil
		}

		// Check whether current file is a view element (e.g. "_button.html").
		if strings.HasPrefix(b, *elemPref) {
			es.set(path.Dir(relNorm), p)
			return nil
		}

		// Otherwise, just add it to the list of templates.
		tpls[relNorm] = p
		return nil
	})

	// Parse templates and register them.
	var err error
	for relNorm, p := range tpls {
		t := template.New(relNorm).Funcs(funcs).Delims(*delimLeft, *delimRight)
		dir := path.Dir(relNorm)

		// Check whether current template must have
		// a layout file and/or element views.
		l := ls.path(dir)
		e := es.path(dir)
		log.Print("\t")
		switch true {
		case l != "" && len(e) > 0: // Both layout and element files.
			log.Printf("%s (%s); %v", p, l, e)
			m[relNorm], err = t.ParseFiles(append(e, p, l)...)
		case l != "": // Layout file only.
			log.Printf("%s (%s)", p, l)
			m[relNorm], err = t.ParseFiles(p, l)
		case len(e) > 0: // Element files only.
			log.Printf("%s (%s); %v", p, p, e)
			m[relNorm], err = t.ParseFiles(append(e, p)...)
		default: // Neither layout nor element files.
			log.Printf("%s (%s)", p, p)
			m[relNorm], err = t.ParseFiles(p)
		}

		// Make sure there were no errors during parsing.
		if err != nil {
			log.Panicf(`Failed to parse "%s". Error: %v.`, p, err)
		}
	}
	return m
}

// elements stores information about directories and element templates
// associated with them. The example value is:
//	App/Profiles: _file1.html, file2.html, fileN.html
type elements map[string][]string

// set gets a directory path and an element path and adds
// the latter to the elements.
func (e elements) set(dir, path string) {
	if _, ok := e[dir]; !ok {
		e[dir] = []string{path}
		return
	}
	e[dir] = append(e[dir], path)
}

// path gets a template directory path and returns a number of template
// element paths that were found in the directory or in directories
// of the higher level.
func (e elements) path(dir string) (res []string) {
	// Check whether current directory has element templates.
	if ls, ok := e[dir]; ok {
		res = append(res, ls...)
	}

	// If this is a root directory, return the result as is.
	if dir == "." {
		return
	}

	// Otherwise, check higher level directories as well.
	res = append(res, e.path(path.Join(dir, ".."))...)
	return
}

// layouts stores information about directories that have layout templates.
// If there is "App/Profiles/Layout.html", the value will be:
//	App/Profiles: true
type layouts map[string]bool

// path gets a template directory path and returns associated Layout template path.
// E.g. in case of the following listing:
//	- views/
//		+ Layout.html
//		- App/
//			+ Index.html
//			- Profiles/
//				+ Index.html
//				+ Layout.html
// If "App/" is an input dir argument, the result must be:
//	views/Layout.html
// And if "App/Profiles/" is an input dir argument, the result must be:
//	views/App/Profiles/Layout.html
// The input path must be in a normalized (i.e. *nix) form.
func (l layouts) path(dir string) string {
	// If the requested directory has a layout file, return its path.
	if l[dir] {
		// Join the directory with the full path to the views.
		// Use the OS specific format of the path (and hence "filepath" package).
		return filepath.Join(*views, filepath.FromSlash(dir), *layout)
	}

	// If this is a root directory but there is no layout file, return nothing.
	if dir == "." {
		return ""
	}

	// Otherwise, check higher level directories.
	return l.path(path.Join(dir, ".."))
}
