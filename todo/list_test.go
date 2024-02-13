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

	list.Add("foobar", nil)

	assert.Equal(t, 1, len(list.Items))
	assert.Equal(t, "foobar", list.Items[0].Title)
	assert.False(t, list.Items[0].IsDone, "new item should not be done")
}

func Test_AddItem_Validation(t *testing.T) {
	list := NewList()

	list.Add("", nil)

	assert.Equal(t, 0, len(list.Items), "empty items not allowed")
}

func Test_edit_ok(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	id := MustNewItemId("111")
	list.Add("foo", id)

	err := list.Edit(id, "newTitle")

	assert.NoError(err)
	assert.Equal("newTitle", list.Items[0].Title)
	assert.Equal(true, list.Items[0].IsModified)
	assert.Equal(false, list.Items[0].IsDeleted)
}

func Test_edit_deletes(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	id := MustNewItemId("111")
	list.Add("foo", id)

	err := list.Edit(id, "")

	assert.NoError(err)
	assert.Equal(true, list.Items[0].IsDeleted)
}

func Test_edit_notExistent(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	id := MustNewItemId("111")
	list.Add("foo", id)
	list.Destroy(id)

	err := list.Edit(id, "new title")

	assert.Error(err)
}

func Test_edit_destroyedItem(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	id := MustNewItemId("111")

	err := list.Edit(id, "new title")

	assert.Error(err)
}

func Test_Toggle_OK(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add1(&Item{Title: "foo", Id: MustNewItemId("100")})
	list.Add1(&Item{Title: "bar", Id: MustNewItemId("200")})

	_ = list.Toggle(MustNewItemId("200"))
	assert.True(list.Items[1].IsDone, "after one toggle")

	_ = list.Toggle(MustNewItemId("200"))
	assert.False(list.Items[1].IsDone, "after another toggle")
}

func Test_Toggle_error(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add1(&Item{Title: "foo", Id: MustNewItemId("1")})

	assert.Error(list.Toggle(MustNewItemId("100")))
	assert.NoError(list.Toggle(MustNewItemId("1")))
}

func Test_Destroy(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add1(&Item{Title: "100", Id: MustNewItemId("100")})
	list.Add1(&Item{Title: "200", Id: MustNewItemId("200")})
	list.Add1(&Item{Title: "300", Id: MustNewItemId("300")})

	list.Destroy(MustNewItemId("200"))

	assert.Equal([]string{"100", "300"}, titles(list.AllItems()))
}

func Test_listAllItems(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero", nil)
	list.AddCompleted("one")
	list.Add("two", nil)

	actual := list.AllItems()

	expected := []string{
		"zero", "one", "two",
	}
	assert.Equal(expected, titles(actual))
}

func Test_listCompletedItems(t *testing.T) {
	assert := assert.New(t)
	list := NewList()
	list.Add("zero", nil)
	list.AddCompleted("one")
	list.Add("two", nil)

	actual := list.CompletedItems()

	expected := []string{
		"one",
	}
	assert.Equal(expected, titles(actual))
}

func titles(actual []*Item) []string {
	result := make([]string, len(actual))
	for i, item := range actual {
		result[i] = item.Title
	}
	return result
}
