package sessions

import (
	"testing"
)

func TestInit(t *testing.T) {
	Init()
	if s == nil {
		t.Errorf("Securecookie expected to be initialized, but it is empty.")
	}
}
