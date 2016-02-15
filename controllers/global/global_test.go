package global

import (
	"testing"
)

func TestGlobalBefore(t *testing.T) {
	g := Global{}
	g.Before()
	if g.Context == nil {
		t.Error("Before method must allocate Context.")
	}
}

func TestObject(t *testing.T) {
	o := Object{
		"string1": "value1",
		"int1":    100,
	}
	for _, v := range []struct {
		key string
		i   int
		s   string
	}{
		{"keyThatDoesntExist", 0, ""},
		{"string1", 0, "value1"},
		{"int1", 100, ""},
	} {
		i := o.Int(v.key)
		if i != v.i {
			t.Errorf(`Incorrect result of Int. Expected "%d", got "%d".`, v.i, i)
		}
		s := o.String(v.key)
		if s != v.s {
			t.Errorf(`Incorrect result of String. Expected "%s", got "%s".`, v.s, s)
		}
	}
}
