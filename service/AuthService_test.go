package service

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	dto "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	t.Run("ShouldReturnAuth", func(t *testing.T) {
		mockCustomerService := new(CustomerServiceMock)
		mockRefreshTokenService := new(RefreshTokenServiceMock)
		mockBlacklistService := new(BlacklistServiceMock)
		authService := NewAuthService(mockCustomerService, mockRefreshTokenService, mockBlacklistService)

		request := req.CustomerRequest{
			Username: "johndoe",
			Password: "password",
		}

		customer := entity.Customer{
			Id:       "id-1",
			Username: request.Username,
			Password: "$2a$10$Ghu/KGIz/UnyrAnvNG7JcODckVqHgXA/Un7/MFKqz/CqhQ2BFGJlK",
		}

		refreshToken := entity.RefreshToken{
			RefreshToken: "refresh-token-1",
			CustomerId:   customer.Id,
			ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour).Format(time.RFC3339),
		}

		mockCustomerService.Mock.On("GetCustomerByUsernameAuth", request.Username).
			Return(customer, nil)

		mockRefreshTokenService.Mock.On("GenerateRefreshToken", customer.Id).
			Return(refreshToken, nil)

		login, err := authService.Login(request)
		assert.Nil(t, err)
		assert.NotNil(t, login.AccessToken)
		assert.NotNil(t, login.RefreshToken)

		id, err := utils.GetCustomerIdFromClaims(login.AccessToken)
		assert.Nil(t, err)
		assert.Equal(t, customer.Id, id)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockCustomerService := new(CustomerServiceMock)
		mockRefreshTokenService := new(RefreshTokenServiceMock)
		mockBlacklistService := new(BlacklistServiceMock)
		authService := NewAuthService(mockCustomerService, mockRefreshTokenService, mockBlacklistService)

		request := req.CustomerRequest{
			Username: "johndoe",
			Password: "",
		}

		customer := entity.Customer{
			Id:       "id-1",
			Username: request.Username,
			Password: "$2a$10$Ghu/KGIz/UnyrAnvNG7JcODckVqHgXA/Un7/MFKqz/CqhQ2BFGJlK",
		}

		mockCustomerService.Mock.On("GetCustomerByUsernameAuth", request.Username).
			Return(customer, nil)

		login, err := authService.Login(request)
		assert.Equal(t, constants.LoginUnauthorizedError, err.Error())
		assert.Equal(t, dto.AuthResponse{}, login)
	})
}

func TestLogout(t *testing.T) {
	mockCustomerService := new(CustomerServiceMock)
	mockRefreshTokenService := new(RefreshTokenServiceMock)
	mockBlacklistService := new(BlacklistServiceMock)
	authService := NewAuthService(mockCustomerService, mockRefreshTokenService, mockBlacklistService)

	accessToken := "access-token-1"
	refreshToken := "refresh-token-1"

	mockBlacklistService.Mock.On("BlacklistToken", accessToken).
		Return(nil)

	mockRefreshTokenService.Mock.On("DeleteRefreshToken", refreshToken).
		Return(nil)

	err := authService.Logout(accessToken, refreshToken)
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

	customer := entity.Customer{
		Id:       "customer-id-1",
		Username: "johndoe",
		Password: "password",
	}

	mockRefreshTokenService.Mock.On("RotateRefreshToken", refreshToken).
		Return(newRefreshToken, nil)

	mockCustomerService.Mock.On("GetCustomerByIdAuth", newRefreshToken.CustomerId).
		Return(customer, nil)

	token, err := authService.GetNewAccessToken(refreshToken)
	assert.Nil(t, err)
	assert.NotNil(t, token.AccessToken)
	assert.NotNil(t, token.RefreshToken)

	id, err := utils.GetCustomerIdFromClaims(token.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, customer.Id, id)
}
