package todo

import "testing"
import "github.com/stretchr/testify/assert"

func Test_index(t *testing.T) {
	list := NewList()

	assert.Equal(t, 0, len(list.Items))
}

func Test_AddItem_ok(t *testing.T) {
	list := NewList()

	list.Add("foobar")

	assert.Equal(t, 1, len(list.Items))
	assert.Equal(t, "foobar", list.Items[0].Title)
	assert.False(t, list.Items[0].IsDone, "new item should not be done")
}

func Test_AddItem_Validation(t *testing.T) {
	list := NewList()

	list.Add("")

	assert.Equal(t, 0, len(list.Items), "empty items not allowed")
}

func Test_IDs_sequential(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("first")
	list.Add("second")

	assert.Equal(0, list.Items[0].Id)
	assert.Equal(1, list.Items[1].Id)
}
