package utils

import (
	"PaymentAPI/constants"
	"errors"
	"github.com/gin-gonic/gin"
)

func GetAuthenticatedUser(c *gin.Context) (string, error) {
	// Get access token from context
	accessToken, exists := c.Get("accessToken")
	if !exists {
		return "", errors.New(constants.AccessTokenNotFoundError)
	}

	// Check authenticated user
	id, err := GetCustomerIdFromClaims(accessToken.(string))
	if err != nil {
		return "", err
	}

	return id, nil
}
