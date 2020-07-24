package gin

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError(t *testing.T) {
	baseError := errors.New("test error")
	err := &Error{
		Err:  baseError,
		Type: ErrorTypePrivate,
	}
	assert.Equal(t, err.Error(), baseError.Error())
	assert.Equal(t, H{"error": baseError.Error()}, err.JSON())

}
