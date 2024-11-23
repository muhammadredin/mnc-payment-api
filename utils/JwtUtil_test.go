package utils

import (
	dto "PaymentAPI/dto/response"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateAccessToken(t *testing.T) {
	customer := dto.CustomerResponse{
		Id:       "id-1",
		Username: "johndoe",
	}

	token, err := GenerateAccessToken(customer)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
}
