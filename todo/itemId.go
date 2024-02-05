package todo

import (
	"errors"
	"regexp"
)

// ItemId represents the ID of an item
// This type serves two purposes:
//   - We want to avoid "prinitive obsession", that is, representing everything with primitive types
//   - We want to ensure that not any string can be used as an ItemId. This is important, in a web app, for security.
//     Since we will be building ItemId's using strings coming from outside the application, we want to frustrate any
//     attempts to break the application by supplying inappropriate strings.
//
// The exported type is opaque; all you know is that it's something that can be represented as a string.
// The `implementsItemId` method is not exported, making it impossible for other packages to create an instance of this
// interface
type ItemId interface {
	String() string
	implementsItemId()
}

// The actual implementation of an ItemId is itemId, which is a wrapper type around string
type itemId string

// Satisfy the requirement to implement the exported interface ItemId
func (id itemId) implementsItemId() {
}

func (id itemId) String() string {
	return string(id)
}

// A valid id is a nonempty sequence of digits
var validId = regexp.MustCompile("^\\d+$")

// NewItemId is the only way to create an ItemId
func NewItemId(s string) (ItemId, error) {
	const maxLength = 10
	if len(s) == 0 || len(s) > maxLength {
		return nil, errors.New("invalid id length")
	}
	if !validId.MatchString(s) {
		return nil, errors.New("invalid characters in id")
	}
	return itemId(s), nil
}

func MustNewTypeId(s string) ItemId {
	id, err := NewItemId(s)
	if err != nil {
		panic(err.Error())
	}
	return id
}
