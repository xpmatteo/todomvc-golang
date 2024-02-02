package todo

import "testing"
import "github.com/stretchr/testify/assert"

func Test_index(t *testing.T) {
	list := NewList()

	assert.Equal(t, 0, len(list.Items))
}

type Item struct {
	Title  string
	IsDone bool
}

type List struct {
	Items []Item
}

func NewList() List {
	return List{}
}
