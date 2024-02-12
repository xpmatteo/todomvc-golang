package web

import (
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/todomvc-golang/todo"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_dataForTheTemplate_items(t *testing.T) {
	model := todo.NewList()
	model.Add("zero", nil)
	model.AddCompleted("one")

	cases := []struct {
		name          string
		url           string
		expectedItems []*todo.Item
	}{
		{"all", "/", []*todo.Item{model.Items[0], model.Items[1]}},
		{"active", "/active", []*todo.Item{model.Items[0]}},
		{"complete", "/completed", []*todo.Item{model.Items[1]}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			data := viewModel(model, httptest.NewRequest(http.MethodGet, test.url, nil))
			assert.Equal(t, test.expectedItems, data["Items"])
		})
	}
}

func Test_dataForTheTemplate_ItemsCount(t *testing.T) {
	model := todo.NewList()
	model.Add("zero", nil)
	model.AddCompleted("one")

	data := viewModel(model, httptest.NewRequest(http.MethodGet, "/completed", nil))

	assert.Equal(t, 2, data["ItemsCount"])
}

func Test_dataForTheTemplate_itemsLeftLabel(t *testing.T) {
	assert := assert.New(t)
	list := todo.NewList()
	list.Add("zero", idZero)
	list.Add("one", idOne)

	data := viewModel(list, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal("2 items left", data["ItemsLeft"])

	_ = list.Toggle(idZero)
	data = viewModel(list, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal("1 item left", data["ItemsLeft"])

	_ = list.Toggle(idOne)
	data = viewModel(list, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal("0 items left", data["ItemsLeft"])
}
