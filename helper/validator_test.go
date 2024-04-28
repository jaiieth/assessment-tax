package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	assert.NotNil(t, v)
}

func TestValidate(t *testing.T) {
	t.Run("StructWithNoErrors", func(t *testing.T) {
		cv := NewValidator()
		input := struct {
			Name  string `validate:"required"`
			Email string `validate:"required,email"`
		}{
			Name:  "John Doe",
			Email: "johndoe@example.com",
		}
		err := cv.Validate(input)
		assert.NoError(t, err)
	})

	t.Run("StructWithNilValue", func(t *testing.T) {
		cv := NewValidator()
		var input *struct {
			Name  string `validate:"required"`
			Email string `validate:"required,email"`
		}
		err := cv.Validate(input)
		assert.Error(t, err)
	})
}
