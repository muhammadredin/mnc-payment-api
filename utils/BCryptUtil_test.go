package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBcryptEncoder(t *testing.T) {
	password := "password"
	encodedPassword := BCryptEncoder(password)

	t.Run("ShouldReturnTrueOnEqual", func(t *testing.T) {
		assert.NotEmpty(t, encodedPassword, "Password encoder return is empty")
		isEqual := BCryptCompare(password, []byte(encodedPassword))
		assert.True(t, isEqual, "Password not encrypted properly")
	})

	t.Run("ShouldReturnFalseOnNotEqual", func(t *testing.T) {
		isEqual := BCryptCompare("passwords", []byte(encodedPassword))
		assert.False(t, isEqual, "Return is not false")
	})
}
