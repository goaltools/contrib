package denco

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestRouter_SpecialLabels(t *testing.T) {
	// Preserve the old values of NotFound and MethodNotAllowed handlers.
	nf := NotFound
	mna := MethodNotAllowed
	ise := InternalServerError

	rs := []struct {
		Method, Pattern, Label string
		Handler                http.HandlerFunc
	}{
		{"GET", "/", "404", testHandlerFunc},
		{"GET", "/profile/:name", "405", testHandlerFuncHelloWorld},
		{"GET", "/", "500", testHandlerFunc},
	}

	_, err := Build(rs)
	if err != nil {
		t.Errorf("Failed to build a handler using Build() function. Error: %s.", err)
	}

	i := 0
	for _, hs := range [][]http.HandlerFunc{
		{NotFound, testHandlerFunc},
		{MethodNotAllowed, testHandlerFuncHelloWorld},
		{InternalServerError, testHandlerFunc},
	} {
		r, _ := http.NewRequest("GET", "/", nil)
		w1 := httptest.NewRecorder()
		w2 := httptest.NewRecorder()

		hs[0](w1, r)
		hs[1](w2, r)

		if !reflect.DeepEqual(w1.Body, w2.Body) {
			t.Errorf("%v: Handler wasn't set.", rs[i])
		}
		i++
	}

	// Restore the old values of errors handlers.
	NotFound = nf
	MethodNotAllowed = mna
	InternalServerError = ise
}

func TestRouter(t *testing.T) {
	rs := Routes{
		Get("/", testHandlerFunc),
		Get("/profile/:name", testHandlerFunc),
		Get("/profile/:name", testHandlerFuncHelloWorld), // This should override the previous route.
		Post("/profile/:name", testHandlerFunc),
		Head("/profile/:name", testHandlerFunc),
		Put("/profile/:name", testHandlerFunc),
		Delete("/profile/:name", testHandlerFunc),
		Do("GET", "/profile/update", testHandlerFunc),
		Get("/panic", testHandlerFuncPanic),
	}
	rs1 := []struct {
		Method, Pattern, Label string
		Handler                http.HandlerFunc
	}{
		{"GET", "/", "", testHandlerFunc},
		{"GET", "/profile/:name", "", testHandlerFunc},
		{"GET", "/profile/:name", "", testHandlerFuncHelloWorld}, // Must be used as a NotFound handler.
		{"POST", "/profile/:name", "", testHandlerFunc},
		{"HEAD", "/profile/:name", "", testHandlerFunc},
		{"PUT", "/profile/:name", "", testHandlerFunc},
		{"DELETE", "/profile/:name", "", testHandlerFunc},
		{"GET", "/profile/update", "", testHandlerFunc},
		{"GET", "/panic", "", testHandlerFuncPanic},
	}

	// Creating a new router manually.
	r := NewRouter()
	err := r.Handle(rs).Build()
	if err != nil {
		t.Errorf("Failed to build a handler using manual method. Error: %s.", err)
	}

	server := httptest.NewServer(r)
	defer server.Close()

	// Using a Build shortcut.
	h, err := rs.Build()
	if err != nil {
		t.Errorf("Failed to build a handler using Build shortcut. Error: %s.", err)
	}

	server1 := httptest.NewServer(h)
	defer server1.Close()

	// Using a Build function.
	h, err = Build(rs1)
	if err != nil {
		t.Errorf("Failed to build a handler using Build() function. Error: %s.", err)
	}

	server2 := httptest.NewServer(h)
	defer server2.Close()

	for _, v := range []struct {
		status                 int
		method, path, expected string
	}{
		{
			200, "GET", "/",
			fmt.Sprintf("method: GET, path: /, form: %v", url.Values{}),
		},
		{
			200, "GET", "/profile/john",
			fmt.Sprintf("Hello, world!\nmethod: GET, path: /profile/john, form: %v", url.Values{
				"name": {"john"},
			}),
		},
		{
			200, "POST", "/profile/jane",
			fmt.Sprintf("method: POST, path: /profile/jane, form: %v", url.Values{
				"name": {"jane"},
			}),
		},
		{
			200, "HEAD", "/profile/james", "",
		},
		{
			200, "PUT", "/profile/alice",
			fmt.Sprintf("method: PUT, path: /profile/alice, form: %v", url.Values{
				"name": {"alice"},
			}),
		},
		{
			200, "DELETE", "/profile/bob",
			fmt.Sprintf("method: DELETE, path: /profile/bob, form: %v", url.Values{
				"name": {"bob"},
			}),
		},
		{
			200, "GET", "/profile/update",
			fmt.Sprintf("method: GET, path: /profile/update, form: %v", url.Values{}),
		},
		{
			405, "POST", "/", http.StatusText(405) + "\n",
		},
		{
			404, "POST", "/qwerty", "404 page not found\n",
		},
		{
			500, "GET", "/panic", http.StatusText(500) + "\n",
		},
	} {
		for _, s := range []*httptest.Server{server, server1, server2} {
			req, err := http.NewRequest(v.method, s.URL+v.path, nil)
			if err != nil {
				t.Errorf("Failed to create a new request. Error: %s.", err)
				continue
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("Cannot do a request. Error: %s.", err)
				continue
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Did not manage to read a response body. Error: %s.", err)
			}
			actual := string(body)
			if res.StatusCode != v.status || actual != v.expected {
				t.Errorf(
					`%s "%s" => %#v %#v, expected %#v %#v.`,
					v.method, v.path, res.StatusCode, actual, v.status, v.expected,
				)
			}
		}
	}

}

func testHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "method: %s, path: %s, form: %v", r.Method, r.URL.Path, r.Form)
}

func testHandlerFuncHelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!\n")
	testHandlerFunc(w, r)
}
func testHandlerFuncPanic(w http.ResponseWriter, r *http.Request) {
	panic("something")
}
