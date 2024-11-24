package service

import (
	dto "PaymentAPI/dto/response"
	"PaymentAPI/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogin(t *testing.T) {
	mockCustomerService := new(CustomerServiceMock)
	mockRefreshTokenService := new(RefreshTokenServiceMock)
	mockBlacklistService := new(BlacklistServiceMock)
	authService := NewAuthService(mockCustomerService, mockRefreshTokenService, mockBlacklistService)

	username := "johndoe"

	customer := dto.CustomerResponse{
		Id:       "id-1",
		Username: username,
	}

	refreshToken := "refresh-token-1"

	mockCustomerService.Mock.On("GetCustomerByUsername", username).
		Return(customer, nil)
	mockRefreshTokenService.Mock.On("GenerateRefreshToken", customer.Id).
		Return(refreshToken, nil)

	login, err := authService.Login(username)
	assert.Nil(t, err)

	claims, err := utils.ParseAndVerifyAccessToken(login.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, customer.Id, claims["sub"])
	assert.NotNil(t, login.RefreshToken)
}

func TestLogout(t *testing.T) {
	mockCustomerService := new(CustomerServiceMock)
	mockRefreshTokenService := new(RefreshTokenServiceMock)
	mockBlacklistService := new(BlacklistServiceMock)
	authService := NewAuthService(mockCustomerService, mockRefreshTokenService, mockBlacklistService)

	accessToken := "access-token-1"

	mockBlacklistService.Mock.On("BlacklistToken", accessToken).
		Return(nil)

	err := authService.Logout(accessToken)
	assert.Nil(t, err)
}
