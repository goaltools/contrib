package requests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestRequestsInitially(t *testing.T) {
	vs := url.Values{
		"key1": {"value1"},
		"key2": {"value2_a", "value2_b"},
		"key3": {"value3_a", "value3_b"},
	}

	// Test different kinds of content types.
	for _, h := range []string{
		"application/x-www-form-urlencoded",
	} {
		c := &Requests{}

		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "test?x=z", strings.NewReader(vs.Encode()))
		assertNil(t, err)
		r.Header.Set("Content-Type", h)

		// Imitating values that was passed by contrib/routers/denco
		// using Form of the request.
		// Only the first value of every key will be joined with the result r.Form
		// as the router doesn't pass more than that anyway.
		r.Form = url.Values{
			"router_key1": {"value1"},
			"router_key2": {"value2_a", "value2_b"},
			"key3":        {"router_value3_a", "router_value3_b"},
		}

		exp := url.Values{
			"x":           {"z"},
			"key1":        {"value1"},
			"key2":        {"value2_a", "value2_b"},
			"router_key1": {"value1"},
			"router_key2": {"value2_a"},
			"key3":        {"value3_a", "value3_b", "router_value3_a"},
		}

		// After Initially is called both requests' values and
		// the router's values must be combined.
		finish := c.Initially(w, r, []string{"TestController", "TestAction"})
		if finish {
			t.Errorf("Magic method Initially unexpectedly returned `finish == true` (Content-Type: %v).", h)
			t.FailNow()
		}

		if !reflect.DeepEqual(r, c.Request) {
			t.Errorf("Request was not saved for use in the child controllers (Content-Type: %v).", h)
		}

		if !reflect.DeepEqual(c.Request.Form, exp) {
			t.Errorf("Expected %v, %v. Got %v, %v (Content-Type: %v).", exp, false, c.Request.Form, finish, h)
		}
	}
}

func assertNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Got unexpected error: %v.", err)
	}
}
