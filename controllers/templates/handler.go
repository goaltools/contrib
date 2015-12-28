package templates

import (
	"net/http"
)

// Handler is a templates handler that implements http.Handler interface.
type Handler struct {
	context  map[string]interface{} // Variables to be passed to the template.
	template string                 // Path to the template to be rendered.
	status   int                    // Expected status code of the response.
}

// Apply writes to response the result received from action.
func (t *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set status of the response.
	if t.status == 0 {
		t.status = http.StatusOK
	}
	w.Header().Set("Content-Type", *contType)

	// If required template exists, execute it.
	if tpl, ok := templates[t.template]; ok {
		w.WriteHeader(t.status)
		err := tpl.ExecuteTemplate(w, *layoutBl, t.context)
		if err != nil {
			go Log.Println(err)
		}
		return
	}

	// Otherwise, show internal server error.
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
	go Log.Printf(`Template "%s" does not exist.`, t.template)
}
