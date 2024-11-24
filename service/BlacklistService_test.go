package service

import (
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsBlacklisted(t *testing.T) {
	t.Run("ShouldReturnTrue", func(t *testing.T) {
		mockBlacklistRepository := new(repository.BlacklistRepositoryMock)
		blacklistService := NewBlacklistService(mockBlacklistRepository)

		accessToken := "access-token-1"
		blacklists := []entity.Blacklist{
			{
				AccessToken: "access-token-1",
				ExpiresAt:   time.Now().Add(time.Duration(5) * time.Minute).Format(time.RFC3339),
			},
		}

		mockBlacklistRepository.Mock.On("GetAll").
			Return(blacklists, nil)

		blacklisted, err := blacklistService.IsBlacklisted(accessToken)
		assert.Nil(t, err)
		assert.True(t, blacklisted)
	})

	t.Run("ShouldReturnFalse", func(t *testing.T) {
		mockBlacklistRepository := new(repository.BlacklistRepositoryMock)
		blacklistService := NewBlacklistService(mockBlacklistRepository)

		accessToken := "access-token-1"

		mockBlacklistRepository.Mock.On("GetAll").
			Return([]entity.Blacklist{}, nil)

		blacklisted, err := blacklistService.IsBlacklisted(accessToken)
		assert.Nil(t, err)
		assert.False(t, blacklisted)
	})
}

func TestBlacklistToken(t *testing.T) {
	t.Run("ShouldBlacklistToken", func(t *testing.T) {
		mockBlacklistRepository := new(repository.BlacklistRepositoryMock)
		blacklistService := NewBlacklistService(mockBlacklistRepository)

		accessToken := "access-token-1"

		mockBlacklistRepository.Mock.On("GetAll").
			Return([]entity.Blacklist{}, nil)

		mockBlacklistRepository.Mock.On("CreateBlacklist", accessToken).
			Return(nil)

		err := blacklistService.BlacklistToken(accessToken)
		assert.Nil(t, err)
	})

	t.Run("ShouldNotBlacklistToken", func(t *testing.T) {
		mockBlacklistRepository := new(repository.BlacklistRepositoryMock)
		blacklistService := NewBlacklistService(mockBlacklistRepository)

		accessToken := "access-token-1"

		blacklists := []entity.Blacklist{
			{
				AccessToken: "access-token-1",
				ExpiresAt:   time.Now().Add(time.Duration(5) * time.Minute).Format(time.RFC3339),
			},
		}

		mockBlacklistRepository.Mock.On("GetAll").
			Return(blacklists, nil)

		err := blacklistService.BlacklistToken(accessToken)
		assert.Nil(t, err)
	})
}
