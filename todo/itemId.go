package todo

import (
	"errors"
	"regexp"
)

type itemId string

var validId = regexp.MustCompile("^\\d+$")

func NewItemId(s string) (itemId, error) {
	const maxLength = 10
	if len(s) == 0 || len(s) > maxLength {
		return "invalid", errors.New("invalid id length")
	}
	if !validId.MatchString(s) {
		return "invalid", errors.New("invalid characters in id")
	}
	return itemId(s), nil
}

func (id itemId) String() string {
	return string(id)
}
