package repository

import (
	"PaymentAPI/constants"
	dto "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"PaymentAPI/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestGetAll(t *testing.T) {
	mockJsonFileHandler := new(storage.BlacklistJsonFileHandlerMock[entity.Blacklist])
	blacklistRepository := NewBlacklistRepository(mockJsonFileHandler)

	blacklistTokenList := []entity.Blacklist{
		{
			AccessToken: "access-token-1",
			ExpiresAt:   time.Now().Add(time.Duration(5) * time.Minute).String(),
		},
	}

	mockJsonFileHandler.Mock.On("ReadFile", constants.BlacklistJsonPath).
		Return(blacklistTokenList, nil)

	blacklists, err := blacklistRepository.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(blacklists))
}

func TestCreateBlacklist(t *testing.T) {
	t.Run("ShouldSuccessCreateBlacklist", func(t *testing.T) {
		mockJsonFileHandler := new(storage.BlacklistJsonFileHandlerMock[entity.Blacklist])
		blacklistRepository := NewBlacklistRepository(mockJsonFileHandler)

		customer := dto.CustomerResponse{
			Id:       "user-1",
			Username: "John Doe",
		}
		token, _ := utils.GenerateAccessToken(customer)
		tokenExpiryDate, _ := utils.GetExpirationFromClaimsAsString(token)

		blacklist := entity.Blacklist{
			AccessToken: token,
			ExpiresAt:   tokenExpiryDate,
		}

		mockJsonFileHandler.Mock.On("ReadFile", constants.BlacklistJsonPath).
			Return([]entity.Blacklist{}, nil)

		mockJsonFileHandler.Mock.On("WriteFile", mock.MatchedBy(func(blacklists []entity.Blacklist) bool {
			if len(blacklists) != 1 {
				return false
			}
			blacklistedToken := blacklists[0]
			return blacklistedToken.AccessToken == blacklist.AccessToken
		}), constants.BlacklistJsonPath).
			Return(mock.Anything, nil)

		err := blacklistRepository.CreateBlacklist(token)
		assert.Nil(t, err)
		mockJsonFileHandler.Mock.AssertExpectations(t)
	})

	t.Run("ShouldReturnErrorOnInvalidToken", func(t *testing.T) {
		mockJsonFileHandler := new(storage.BlacklistJsonFileHandlerMock[entity.Blacklist])
		blacklistRepository := NewBlacklistRepository(mockJsonFileHandler)

		mockJsonFileHandler.Mock.On("ReadFile", constants.BlacklistJsonPath).
			Return([]entity.Blacklist{}, nil)

		err := blacklistRepository.CreateBlacklist("")
		assert.Equal(t, constants.JwtTokenInvalidError, err.Error())
	})
}
