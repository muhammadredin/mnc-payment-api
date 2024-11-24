package service

import (
	dto "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	refreshToken := entity.RefreshToken{
		RefreshToken: "refresh-token-1",
		CustomerId:   customer.Id,
		ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour).Format(time.RFC3339),
	}

	mockCustomerService.Mock.On("GetCustomerByUsername", username).
		Return(customer, nil)
	mockRefreshTokenService.Mock.On("GenerateRefreshToken", customer.Id).
		Return(refreshToken, nil)

	login, err := authService.Login(username)
	assert.Nil(t, err)

	id, err := utils.GetCustomerIdFromClaims(login.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, customer.Id, id)
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

func TestGetNewAccessToken(t *testing.T) {
	mockCustomerService := new(CustomerServiceMock)
	mockRefreshTokenService := new(RefreshTokenServiceMock)
	mockBlacklistService := new(BlacklistServiceMock)
	authService := NewAuthService(mockCustomerService, mockRefreshTokenService, mockBlacklistService)

	refreshToken := "refresh-token-1"

	newRefreshToken := entity.RefreshToken{
		RefreshToken: "new-refresh-token-1",
		CustomerId:   "customer-id-1",
		ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour).Format(time.RFC3339),
	}

	customer := dto.CustomerResponse{
		Id:       "customer-id-1",
		Username: "john doe",
	}

	mockRefreshTokenService.Mock.On("RotateRefreshToken", refreshToken).
		Return(newRefreshToken, nil)

	mockCustomerService.Mock.On("GetCustomerById", newRefreshToken.CustomerId).
		Return(customer, nil)

	token, err := authService.GetNewAccessToken(refreshToken)
	assert.Nil(t, err)
	assert.Equal(t, newRefreshToken.RefreshToken, token.RefreshToken)

	id, err := utils.GetCustomerIdFromClaims(token.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, customer.Id, id)
}
