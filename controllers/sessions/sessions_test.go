package sessions

import (
	"net/url"
	"testing"
)

func TestInit(t *testing.T) {
	Init(url.Values{})
	if s == nil {
		t.Errorf("Securecookie expected to be initialized, but it is empty.")
	}
}
