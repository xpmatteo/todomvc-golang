package web

import (
	"bytes"
	"encoding/xml"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xpmatteo/todomvc-golang/todo"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_wellFormedHtml(t *testing.T) {
	templ := template.Must(template.ParseFiles("../templates/index.html"))
	list := todo.NewList()
	model := viewModel(list, get("/"))

	buf := renderTemplate(templ, model)

	assertWellFormedHTML(t, buf)
}

func Test_todoItemsAreShown(t *testing.T) {
	templ := template.Must(template.ParseFiles("../templates/index.html"))
	list := todo.NewList()
	list.Add("Foo", nil)
	list.Add("Bar", nil)
	model := viewModel(list, get("/"))

	buf := renderTemplate(templ, model)

	assertWellFormedHTML(t, buf)
	document := parseHtml(t, buf)
	require.Equal(t, 2, document.Find("ul.todo-list li").Length())
	assert.Equal(t, "Foo", strings.TrimSpace(document.Find("ul.todo-list li:nth-of-type(1)").Text()))
	assert.Equal(t, "Bar", strings.TrimSpace(document.Find("ul.todo-list li:nth-of-type(2)").Text()))
}

func Test_completedItemsGetCompletedClass(t *testing.T) {
	templ := template.Must(template.ParseFiles("../templates/index.html"))
	list := todo.NewList()
	list.Add("Foo", nil)
	list.AddCompleted("Bar")
	model := viewModel(list, get("/"))

	buf := renderTemplate(templ, model)

	assertWellFormedHTML(t, buf)
	document := parseHtml(t, buf)
	selection := document.Find("ul.todo-list li.completed")
	require.Equal(t, 1, selection.Length())
	assert.Equal(t, "Bar", strings.TrimSpace(selection.Text()))
}

func Test_editingItems(t *testing.T) {
	templ := template.Must(template.ParseFiles("../templates/index.html"))
	list := todo.NewList()
	list.Add("Foo", todo.MustNewItemId("22"))
	model := viewModel(list, get("/?edit=22"))

	buf := renderTemplate(templ, model)

	assertWellFormedHTML(t, buf)
	document := parseHtml(t, buf)
	selection := document.Find("ul.todo-list li.editing")
	require.Equal(t, 1, selection.Length())
	assert.Equal(t, "Foo", strings.TrimSpace(selection.Text()))
}

func get(target string) *http.Request {
	return httptest.NewRequest(http.MethodGet, target, nil)
}

func parseHtml(t *testing.T, buf bytes.Buffer) *goquery.Document {
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("Error rendering template %s", err)
	}
	return document
}

func renderTemplate(templ *template.Template, model any) bytes.Buffer {
	var buf bytes.Buffer
	err := templ.Execute(&buf, model)
	if err != nil {
		panic(err)
	}
	return buf
}

func assertWellFormedHTML(t *testing.T, buf bytes.Buffer) {
	decoder := xml.NewDecoder(bytes.NewReader(buf.Bytes()))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity

	for {
		_, err := decoder.Token()
		switch err {
		case io.EOF:
			return // We're done, it's valid!
		case nil:
			// do nothing
		default:
			t.Fatalf("Error parsing html: %s", err)
		}
	}
}
