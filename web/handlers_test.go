package web

import (
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/todomvc-golang/todo"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var templ = template.Must(template.New("index").Parse("<p>{{.Items}}</p>"))

var (
	idZero = todo.MustNewItemId("0")
	idOne  = todo.MustNewItemId("1")
	idTwo  = todo.MustNewItemId("2")
)

var aModel = todo.NewList()

func Test_indexHandler_ok(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)

	MakeIndexHandler(templ, aModel).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p>[]</p>", w.Body.String())
}

func Test_indexHandler_unexpectedPath(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/foo", nil)

	MakeIndexHandler(templ, aModel).ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "Not found\n", w.Body.String())
}

func Test_indexHandler_unexpectedMethod(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", nil)

	MakeIndexHandler(templ, aModel).ServeHTTP(w, r)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "Method not allowed\n", w.Body.String())
}

func Test_indexHandler_editItem(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/?edit=3", nil)
	templ := template.Must(template.New("index").Parse("<p>{{.EditingItemId}}</p>"))

	MakeIndexHandler(templ, aModel).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<p>3</p>", w.Body.String())
}

func Test_indexHandler_editItemNotPassed(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	templ := template.Must(template.New("index").Parse("<p>{{.EditingItemId}}</p>"))

	MakeIndexHandler(templ, aModel).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<p></p>", w.Body.String())
}

func Test_editHandler_ok(t *testing.T) {
	assert := assert.New(t)
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", strings.NewReader("todoItemId=0&todoItemTitle=bar"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	model := todo.NewList()
	model.Add("foo")
	templ := template.Must(template.New("index").Parse("<p>{{len .Model.Items}}</p>"))

	MakeEditHandler(templ, model).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("<p>1</p>", w.Body.String())
	assert.Equal("bar", model.Items[idZero].Title)
}

func Test_editHandler_textIsEmpty(t *testing.T) {
	assert := assert.New(t)
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", strings.NewReader("todoItemId=0&todoItemTitle="))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	model := todo.NewList()
	model.Add("foo")

	MakeEditHandler(templ, model).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("<p>[]</p>", w.Body.String())
	assert.Equal(0, len(model.Items))
}

func Test_destroyHandler_ok(t *testing.T) {
	assert := assert.New(t)
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", strings.NewReader("todoItemId=0"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	model := todo.NewList()
	model.Add("foo")

	MakeDestroyHandler(templ, model).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("<p>[]</p>", w.Body.String())
	assert.Equal(0, len(model.Items))
}

func Test_dataForTheTemplate_items(t *testing.T) {
	model := todo.NewList()
	model.Add("zero")
	model.AddCompleted("one")

	cases := []struct {
		name          string
		url           string
		expectedItems []*todo.Item
	}{
		{"all", "/", []*todo.Item{model.Items[idZero], model.Items[idOne]}},
		{"active", "/active", []*todo.Item{model.Items[idZero]}},
		{"complete", "/completed", []*todo.Item{model.Items[idOne]}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			data := makeDataForTemplate(model, httptest.NewRequest(http.MethodGet, test.url, nil))
			assert.Equal(t, test.expectedItems, data["Items"])
		})
	}
}

func Test_dataForTheTemplate_ItemsCount(t *testing.T) {
	model := todo.NewList()
	model.Add("zero")
	model.AddCompleted("one")

	data := makeDataForTemplate(model, httptest.NewRequest(http.MethodGet, "/completed", nil))

	assert.Equal(t, 2, data["ItemsCount"])
}
