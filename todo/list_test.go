package todo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func Test_Toggle_OK(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("first")
	list.Add("second")
	assert.False(list.Items[1].IsDone, "initially")

	_ = list.Toggle(1)
	assert.True(list.Items[1].IsDone, "after one toggle")

	_ = list.Toggle(1)
	assert.False(list.Items[1].IsDone, "after another toggle")
}

func Test_Toggle_error(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("first")

	assert.Error(list.Toggle(-1))
	assert.Error(list.Toggle(1))
	assert.NoError(list.Toggle(0))
}

func Test_CountItemsLeft(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.Add("one")
	list.Add("two")

	assert.Equal(3, list.ItemsLeft())
}

func Test_CountItemsLeft_countsOnlyNotDone(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.Add("one")
	list.Add("two")

	_ = list.Toggle(1)

	assert.Equal(2, list.ItemsLeft())
}
