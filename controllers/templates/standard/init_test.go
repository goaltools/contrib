package standard

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestElements(t *testing.T) {
	es := elements{
		".":            []string{"_r1.html", "_r2.html"},
		"App":          []string{"_a1.html"},
		"App/Profiles": []string{"_p1.html", "_p2.html"},
	}
	for _, v := range []struct {
		dir string
		exp []string
	}{
		{".", []string{"_r1.html", "_r2.html"}},
		{"App", []string{"_a1.html", "_r1.html", "_r2.html"}},
		{"App/Profiles", []string{"_p1.html", "_p2.html", "_a1.html", "_r1.html", "_r2.html"}},
	} {
		if res := es.path(v.dir); !reflect.DeepEqual(v.exp, res) {
			t.Errorf("Expected %v, got %v.", v.exp, res)
		}
	}
}

func TestLayouts(t *testing.T) {
	ls := layouts{
		".":             false,
		"App":           true,
		"App/Profiles":  true,
		"Errors":        false,
		"Messages":      true,
		"Messages/Test": false,
	}
	for _, v := range []struct {
		dir string
		res string
	}{
		{".", ""},
		{"App", filepath.Join(*views, "App", *layout)},
		{"App/Profiles", filepath.Join(*views, "App", "Profiles", *layout)},
		{"Errors", ""},
		{"Messages", filepath.Join(*views, "Messages", *layout)},
		{"Messages/Test", filepath.Join(*views, "Messages", *layout)},
	} {
		if res := ls.path(v.dir); v.res != res {
			t.Errorf(`Expected: "%s". Got: "%s".`, v.res, res)
		}
	}
}
