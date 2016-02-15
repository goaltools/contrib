package global

import (
	"net/http"
	"testing"
)

func TestObject_Template(t *testing.T) {
	o := Object{}
	if o.Template() != "" {
		t.Fail()
	}
	o.SetTemplate("/some/path")
	if o.Template() != "/some/path" {
		t.Fail()
	}
}

func TestObject_Status(t *testing.T) {
	o := Object{}
	if o.Status() != http.StatusOK {
		t.Fail()
	}
	o.SetStatus(http.StatusForbidden)
	if o.Status() != http.StatusForbidden {
		t.Fail()
	}
}
