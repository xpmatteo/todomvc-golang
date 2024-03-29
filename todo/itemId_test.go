package todo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewItemId(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput ItemId
		expectedErr    error
	}{
		{"ok", "22", itemId("22"), nil},
		{"too short", "", nil, errors.New("invalid id length")},
		{"too long", "01234567890123456789", nil, errors.New("invalid id length")},
		{"contains spaces", "1 1", nil, errors.New("invalid characters in id")},
		{"contains spaces outside", " 1 ", nil, errors.New("invalid characters in id")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewItemId(tt.input)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedOutput, id)
		})
	}
}
