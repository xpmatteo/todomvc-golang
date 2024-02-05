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
	assert.Equal(t, "foobar", list.Items[MustNewTypeId("0")].Title)
	assert.False(t, list.Items[MustNewTypeId("0")].IsDone, "new item should not be done")
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

	assert.Equal(itemId("0"), list.Items[MustNewTypeId("0")].Id)
	assert.Equal(itemId("1"), list.Items[MustNewTypeId("1")].Id)
}

func Test_Toggle_OK(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("first")
	list.Add("second")
	assert.False(list.Items[MustNewTypeId("1")].IsDone, "initially")

	_ = list.Toggle(MustNewTypeId("1"))
	assert.True(list.Items[MustNewTypeId("1")].IsDone, "after one toggle")

	_ = list.Toggle(MustNewTypeId("1"))
	assert.False(list.Items[MustNewTypeId("1")].IsDone, "after another toggle")
}

func Test_Toggle_error(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("first")

	assert.Error(list.Toggle(MustNewTypeId("1")))
	assert.NoError(list.Toggle(MustNewTypeId("0")))
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

	_ = list.Toggle(MustNewTypeId("1"))

	assert.Equal(2, list.ItemsLeft())
}

func Test_Destroy(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.Add("one")
	list.Add("two")

	list.Destroy(MustNewTypeId("1"))

	assert.Equal(2, len(list.Items))
}

func Test_Ids_SurviveAfterDestroy(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.Add("one")
	list.Add("two")
	list.Destroy(MustNewTypeId("1"))

	err := list.Toggle(MustNewTypeId("2"))

	assert.NoError(err)
}

func Test_listAllItems(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.Add("one")
	list.Add("two")
	_ = list.Toggle(MustNewTypeId("1"))

	actual := list.AllItems()

	expected := []*Item{
		list.Items[MustNewTypeId("0")],
		list.Items[MustNewTypeId("1")],
		list.Items[MustNewTypeId("2")],
	}
	assert.Equal(expected, actual)
}

func Test_listCompletedItems(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero")
	list.Add("one")
	list.Add("two")
	_ = list.Toggle(MustNewTypeId("1"))

	actual := list.CompletedItems()

	expected := []*Item{
		list.Items[MustNewTypeId("1")],
	}
	assert.Equal(expected, actual)
}
