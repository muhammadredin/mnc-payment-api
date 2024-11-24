package repository

import (
	"PaymentAPI/constants"
	dto "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreateNewRefreshToken(t *testing.T) {
	mockFileHandler := new(storage.RefreshTokenJsonFileHandlerMock[entity.RefreshToken])
	refreshTokenRepository := NewRefreshTokenRepository(mockFileHandler)

	customer := dto.CustomerResponse{
		Id:       "id-1",
		Username: "johndoe",
	}

	mockFileHandler.Mock.On("ReadFile", constants.RefreshTokenJsonPath).
		Return([]entity.RefreshToken{}, nil)

	mockFileHandler.Mock.On("WriteFile",
		mock.MatchedBy(func(refreshTokens []entity.RefreshToken) bool {
			// Check if the slice has exactly one element
			if len(refreshTokens) != 1 {
				return false
			}

			refreshToken := refreshTokens[0]
			if refreshToken.CustomerId != customer.Id {
				return false
			}

			if _, err := uuid.Parse(refreshToken.RefreshToken); err != nil {
				return false
			}

			return true
		}),
		constants.RefreshTokenJsonPath,
	).Return(constants.JsonWriteSuccess, nil)

	token, err := refreshTokenRepository.CreateRefreshToken(customer.Id)

	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	_, err = uuid.Parse(token)
	assert.Nil(t, err)
}

func TestGetRefreshToken(t *testing.T) {
	mockFileHandler := new(storage.RefreshTokenJsonFileHandlerMock[entity.RefreshToken])
	refreshTokenRepository := NewRefreshTokenRepository(mockFileHandler)

	refreshToken := "token-1"
	expectedRefreshToken := entity.RefreshToken{
		RefreshToken: "token-1",
		CustomerId:   "user-1",
		ExpiresAt:    time.Now().Add(24 * time.Hour).String(),
	}

	t.Run("ShouldReturnRefreshToken", func(t *testing.T) {
		mockFileHandler.Mock.On("ReadFile", mock.Anything).
			Return([]entity.RefreshToken{expectedRefreshToken}, nil)

		token, err := refreshTokenRepository.GetRefreshToken(refreshToken)
		assert.Nil(t, err)
		assert.Equal(t, expectedRefreshToken, token)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockFileHandler.Mock.On("ReadFile", "").
			Return([]entity.RefreshToken{}, nil)

		token, err := refreshTokenRepository.GetRefreshToken("")
		assert.Equal(t, constants.RefreshTokenNotFoundError, err.Error())
		assert.Equal(t, entity.RefreshToken{}, token)
	})
}

func TestDeleteRefreshToken(t *testing.T) {
	mockFileHandler := new(storage.RefreshTokenJsonFileHandlerMock[entity.RefreshToken])
	refreshTokenRepository := NewRefreshTokenRepository(mockFileHandler)

	refreshToken := "token-1"

	storedRefreshToken := entity.RefreshToken{
		RefreshToken: "token-1",
		CustomerId:   "user-1",
		ExpiresAt:    time.Now().Add(24 * time.Hour).String(),
	}

	t.Run("ShouldDeleteRefreshToken", func(t *testing.T) {
		mockFileHandler.Mock.On("ReadFile", constants.RefreshTokenJsonPath).
			Return([]entity.RefreshToken{storedRefreshToken}, nil)

		mockFileHandler.Mock.On("WriteFile",
			mock.MatchedBy(func(refreshTokens []entity.RefreshToken) bool {
				if len(refreshTokens) != 0 {
					return false
				}

				return true
			}),
			constants.RefreshTokenJsonPath,
		).Return(constants.JsonWriteSuccess, nil)

		err := refreshTokenRepository.DeleteRefreshToken(refreshToken)

		assert.Nil(t, err)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockFileHandler.Mock.On("ReadFile", constants.RefreshTokenJsonPath).
			Return([]entity.RefreshToken{}, nil)

		err := refreshTokenRepository.DeleteRefreshToken("")

		assert.Equal(t, constants.RefreshTokenNotFoundError, err.Error())
	})
}
