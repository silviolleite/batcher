package batcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/silviolleite/batcher"
)

func TestError_Error(t *testing.T) {
	testCases := []struct {
		name string
		e    error
		want string
	}{
		{
			"Should return channel closed error message",
			batcher.ErrInputChannelClosed,
			"batch channel is closed",
		},
		{
			"Should return input nil error message",
			batcher.ErrInputIsNil,
			"input is nil",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.e.Error())
		})
	}
}
