package todo

import (
	"errors"
	"regexp"
)

type ItemId string

var validId = mustCompile(regexp.Compile("^\\d+$"))

func mustCompile(r *regexp.Regexp, err error) *regexp.Regexp {
	if err != nil {
		panic(err.Error())
	}
	return r
}

func NewItemId(s string) (ItemId, error) {
	const maxLength = 10
	if len(s) == 0 || len(s) > maxLength {
		return "invalid", errors.New("invalid id length")
	}
	if !validId.MatchString(s) {
		return "invalid", errors.New("invalid characters in id")
	}
	return ItemId(s), nil
}
