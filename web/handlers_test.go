package web

import (
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/todomvc-golang/db"
	"github.com/xpmatteo/todomvc-golang/todo"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var templ = template.Must(template.New("index").Parse("items: {{range $item := .Items}}{{$item.Title}},{{end}}"))

var (
	idZero = todo.MustNewItemId("0")
	idOne  = todo.MustNewItemId("1")
	idTwo  = todo.MustNewItemId("2")
)

func Test_indexHandler_ok(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	repository := db.FakeRepository().Add("item0").Add("item1")

	IndexHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "items: item0,item1,", w.Body.String())
}

func Test_indexHandler_unexpectedPath(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/foo", nil)
	repository := db.FakeRepository()

	IndexHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "Not found\n", w.Body.String())
}

func Test_indexHandler_editItem(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/?edit=3", nil)
	templ := template.Must(template.New("index").Parse("<p>{{.EditingItemId}}</p>"))
	repository := db.FakeRepository()

	IndexHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<p>3</p>", w.Body.String())
}

func Test_indexHandler_editItemNotPassed(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	templ := template.Must(template.New("index").Parse("<p>{{.EditingItemId}}</p>"))
	repository := db.FakeRepository()

	IndexHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<p></p>", w.Body.String())
}

func Test_NewItemHandler_ok(t *testing.T) {
	assert := assert.New(t)
	w, r := postRequest("new-todo=foobar")
	repository := db.FakeRepository()
	var templ = template.Must(template.New("index").Parse("items: {{range $item := .Items}}{{$item.Id}}:{{$item.Title}},{{end}}"))

	NewItemHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("items: 0:foobar,", w.Body.String())
	assert.Equal("foobar", repository.Items[0].Title)
}

func Test_editHandler_ok(t *testing.T) {
	assert := assert.New(t)
	w, r := postRequest("todoItemId=0&todoItemTitle=changedTitle")
	repository := db.FakeRepository().Add("foo")

	EditHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("items: changedTitle,", w.Body.String())
	assert.Equal("changedTitle", repository.Items[0].Title)
}

func Test_editHandler_textIsEmpty(t *testing.T) {
	assert := assert.New(t)
	w, r := postRequest("todoItemId=0&todoItemTitle=")
	repository := db.FakeRepository().Add("foo")

	EditHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("items: ", w.Body.String())
	assert.Equal(0, len(repository.Items))
}

func Test_editHandler_elementNotFound(t *testing.T) {
	assert := assert.New(t)
	w, r := postRequest("todoItemId=123&todoItemTitle=changedTitle")
	repository := db.FakeRepository()

	EditHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("bad todoItemId\n", w.Body.String())
}

func Test_toggleHandler_ok(t *testing.T) {
	assert := assert.New(t)
	repository := db.FakeRepository().Add("zero").Add("one").Add("two")
	w, r := postRequest("todoItemId=1")

	ToggleHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("items: zero,one,two,", w.Body.String())
	list, _ := repository.FindList()
	assert.Equal(true, list.Items[1].IsCompleted)
}

func Test_destroyHandler_ok(t *testing.T) {
	assert := assert.New(t)
	repository := db.FakeRepository().Add("zero").Add("one").Add("two")
	w, r := postRequest("todoItemId=1")

	DestroyHandler(templ, repository).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("items: zero,two,", w.Body.String())
	remainingItems, _ := repository.FindList()
	assert.Equal(2, len(remainingItems.AllItems()))
}

func postRequest(body string) (*httptest.ResponseRecorder, *http.Request) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return w, r
}
