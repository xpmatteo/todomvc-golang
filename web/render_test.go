package web

import (
	"github.com/stretchr/testify/assert"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_rendersJson(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(headerKeyAccept, "application/json, */*")
	vm := ViewModel{"Foo": 123}

	render(w, r, nil, vm)

	assert := assert.New(t)
	assert.Equal("application/json", w.Header().Get("content-type"))
	assert.Equal("{\"Foo\":123}\n", w.Body.String())
}

func Test_rendersHtml(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	vm := ViewModel{"Foo": 123}
	var templFoo = template.Must(template.New("index").Parse("<p>{{.Foo}}</p>"))

	render(w, r, templFoo, vm)

	assert := assert.New(t)
	assert.Equal("text/html; charset=utf-8", w.Header().Get("content-type"))
	assert.Equal("<p>123</p>", w.Body.String())
}

func Test_failsRenderingTemplate(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	vm := ViewModel{"Foo": 123}
	var templFoo = template.Must(template.New("broken").Parse("{{ len .Foo }}"))

	render(w, r, templFoo, vm)

	assert.Equal(t, 500, w.Code)
}

func Test_startsWith(t *testing.T) {
	assert := assert.New(t)
	assert.True(startsWith("ab", "abc"))
	assert.False(startsWith("xx", "abc"))
	assert.False(startsWith("abc", "ab"))
}
