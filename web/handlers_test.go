package web

import (
	"github.com/stretchr/testify/assert"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo"
)

var templ = template.Must(template.New("index").Parse("<p>{{.Model}}</p>"))

func Test_indexHandler_ok(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)

	MakeIndexHandler(templ, "foo").ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p>foo</p>", w.Body.String())
}

func Test_indexHandler_escapesEntities(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)

	MakeIndexHandler(templ, " a < b ").ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p> a &lt; b </p>", w.Body.String())
}

func Test_indexHandler_unexpectedPath(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/foo", nil)

	MakeIndexHandler(templ, "foo").ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "Not found\n", w.Body.String())
}

func Test_indexHandler_unexpectedMethod(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", nil)

	MakeIndexHandler(templ, "foo").ServeHTTP(w, r)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "Method not allowed\n", w.Body.String())
}

func Test_indexHandler_editItem(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/?edit=3", nil)
	templ := template.Must(template.New("index").Parse("<p>{{.EditingItemId}}</p>"))

	MakeIndexHandler(templ, "foo").ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<p>3</p>", w.Body.String())
}

func Test_indexHandler_editItemNotPassed(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	templ := template.Must(template.New("index").Parse("<p>{{.EditingItemId}}</p>"))

	MakeIndexHandler(templ, "foo").ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<p></p>", w.Body.String())
}

func Test_destroyHamdler_ok(t *testing.T) {
	assert := assert.New(t)
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", strings.NewReader("todoItemId=0"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	model := todo.NewList()
	model.Add("foo")

	MakeDestroyHandler(model).ServeHTTP(w, r)

	assert.Equal(http.StatusSeeOther, w.Code)
	assert.Equal("/", w.Header().Get("location"))
	assert.Equal(0, len(model.Items))
}
