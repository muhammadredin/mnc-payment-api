package utils

import (
	"PaymentAPI/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateAccessToken(t *testing.T) {
	customer := entity.Customer{
		Id:       "id-1",
		Username: "johndoe",
		Password: "password",
	}

	token, err := GenerateAccessToken(customer)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
}
