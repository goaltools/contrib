package templates

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var templates = map[string]*template.Template{}

// load gets three input arguments:
// 1. Path to the root of the user application (e.g. "./").
// 2. Path to the views directory relative to the project root (e.g. "views").
// 3. A list of template paths relative to the views directory (in a form of map[string]string).
// It checks whether all the templates exist, parses and registers them.
// It panics if some of the requested templates do not exist or cannot be parsed.
func load(root string, views string, templatePaths map[string]string) {
	log.Println("Loading templates...")
	root = filepath.Join(root, views)

	// Iterating over all available template paths.
	for _, path := range templatePaths {
		// Find base for the current template
		// (either in the current dir or in one of the previous levels).
		var base, cd string
		limit, i := 100, 0
		for {
			if i++; i == limit {
				log.Println("Limit reached when loading templates.")
				break
			}

			b := filepath.Base(path)
			dir := filepath.Join(path[:len(path)-len(b)], cd)
			cd += "../"

			// Check whether this template is a base. If so, do not load
			// any other bases.
			if b == *baseTemplate {
				break
			}

			// Check whether base template exists in the directory.
			base = filepath.Join(dir, *baseTemplate)
			if _, ok := templates[base]; ok || contains(templatePaths, base) {
				break
			}
			base = ""

			// Check whether we have unsuccessfully achieved the top level
			// of the path.
			if strings.HasPrefix(dir, "../") {
				break
			}
		}

		log.Printf("\t%s (%s)", path, base)

		// If the base was found, use it. Otherwise, go without it.
		var err error
		t := template.New(path).Funcs(Funcs).Delims(*delimLeft, *delimRight)
		if base != "" {
			templates[path], err = t.ParseFiles(
				filepath.Join(root, base),
				filepath.Join(root, path),
			)
			showError(root, base, path, err)
			continue
		}
		templates[path], err = t.ParseFiles(filepath.Join(root, path))
		showError(root, base, path, err)
	}
}

// contains returns true if a requested value found
// in the requested slice of strings.
func contains(lst map[string]string, value string) bool {
	for k := range lst {
		if lst[k] == value {
			return true
		}
	}
	return false
}

// showErrors writes an error to log.
func showError(root, base, path string, err error) {
	if err == nil {
		return
	}
	pwd, _ := os.Getwd()
	log.Panicf(
		`Cannot parse "%s" with "%s" as a base template (pwd "%s"). Error: %v.`,
		filepath.Join(root, path), filepath.Join(root, base), pwd, err,
	)
}
