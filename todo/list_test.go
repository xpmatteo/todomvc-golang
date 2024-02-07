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
	assert.Equal(t, "foobar", list.Items[MustNewItemId("0")].Title)
	assert.False(t, list.Items[MustNewItemId("0")].IsDone, "new item should not be done")
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

	assert.Equal(itemId("0"), list.Items[MustNewItemId("0")].Id)
	assert.Equal(itemId("1"), list.Items[MustNewItemId("1")].Id)
}

func Test_Toggle_OK(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("first")
	list.Add("second")
	assert.False(list.Items[MustNewItemId("1")].IsDone, "initially")

	_ = list.Toggle(MustNewItemId("1"))
	assert.True(list.Items[MustNewItemId("1")].IsDone, "after one toggle")

	_ = list.Toggle(MustNewItemId("1"))
	assert.False(list.Items[MustNewItemId("1")].IsDone, "after another toggle")
}

func Test_Toggle_error(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("first")

	assert.Error(list.Toggle(MustNewItemId("1")))
	assert.NoError(list.Toggle(MustNewItemId("0")))
}

func Test_Destroy(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.Add("one")
	list.Add("two")

	list.Destroy(MustNewItemId("1"))

	assert.Equal(2, len(list.Items))
}

func Test_Ids_SurviveAfterDestroy(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.Add("one")
	list.Add("two")
	list.Destroy(MustNewItemId("1"))

	err := list.Toggle(MustNewItemId("2"))

	assert.NoError(err)
}

func Test_listAllItems(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.AddCompleted("one")
	list.Add("two")

	actual := list.AllItems()

	expected := []*Item{
		list.Items[MustNewItemId("0")],
		list.Items[MustNewItemId("1")],
		list.Items[MustNewItemId("2")],
	}
	assert.Equal(expected, actual)
}

func Test_listCompletedItems(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.AddCompleted("one")
	list.Add("two")

	actual := list.CompletedItems()

	expected := []*Item{
		list.Items[MustNewItemId("1")],
	}
	assert.Equal(expected, actual)
}

func Test_listMethodsPreserveInsertionOrder(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.AddCompleted("one")
	list.Add("two")
	list.Add("three")
	list.Add("four")

	expected := []*Item{
		list.Items[MustNewItemId("0")],
		list.Items[MustNewItemId("1")],
		list.Items[MustNewItemId("2")],
		list.Items[MustNewItemId("3")],
		list.Items[MustNewItemId("4")],
	}
	assert.Equal(expected, list.AllItems())

	list.Destroy(MustNewItemId("2"))

	remaining := []*Item{
		list.Items[MustNewItemId("0")],
		list.Items[MustNewItemId("1")],
		list.Items[MustNewItemId("3")],
		list.Items[MustNewItemId("4")],
	}
	assert.Equal(remaining, list.AllItems())
}
