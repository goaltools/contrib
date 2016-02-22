package global

import (
	"net/http"
	"testing"
)

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
