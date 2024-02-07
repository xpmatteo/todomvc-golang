package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GET_rejectsPostRequests(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", nil)

	GETonly(nil).ServeHTTP(w, r)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func Test_GET_acceptsGetRequests(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	GETonly(h).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}
