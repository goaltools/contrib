package standard

import (
	"path"
	"path/filepath"
)

// elements stores information about directories and element templates
// associated with them. The example value is:
//	App/Profiles: _file1.html, file2.html, fileN.html
type elements map[string][]string

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
func (l layouts) path(dir string) (string, bool) {
	// If the requested directory has a layout file, return its path.
	if l[dir] {
		// Join the directory with the full path to the views.
		// Use the OS specific format of the path (and hence "filepath" package).
		return filepath.Join(*views, filepath.FromSlash(dir), *layout), true
	}

	// If this is a root directory but there is no layout file, return nothing.
	if dir == "." {
		return "", false
	}

	// Otherwise, check higher level directories.
	return l.path(path.Join(dir, ".."))
}
