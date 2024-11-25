package service

import (
	"PaymentAPI/logger" // Import the logger package
	"PaymentAPI/repository"
	"github.com/sirupsen/logrus"
)

type BlacklistService interface {
	IsBlacklisted(accessToken string) (bool, error)
	BlacklistToken(accessToken string) error
}

type blacklistService struct {
	blacklistRepository repository.BlacklistRepository
}

// Constructor for BlacklistService
func NewBlacklistService(blacklistRepository repository.BlacklistRepository) BlacklistService {
	return &blacklistService{blacklistRepository: blacklistRepository}
}

// BlacklistToken checks if the token is blacklisted, and if not, adds it to the blacklist
func (b *blacklistService) BlacklistToken(accessToken string) error {
	logger.LogInfo("Attempting to blacklist token", logrus.Fields{
		"accessToken": accessToken,
	})

	// Check if the token is already blacklisted
	isBlacklisted, err := b.IsBlacklisted(accessToken)
	if err != nil {
		logger.LogError("Failed to check if token is blacklisted", logrus.Fields{
			"accessToken": accessToken,
			"error":       err.Error(),
		})
		return err
	}

	// If the token is already blacklisted, return nil (no action needed)
	if isBlacklisted {
		logger.LogInfo("Token is already blacklisted", logrus.Fields{
			"accessToken": accessToken,
		})
		return nil
	}

	// Add the token to the blacklist repository
	err = b.blacklistRepository.CreateBlacklist(accessToken)
	if err != nil {
		logger.LogError("Failed to add token to blacklist", logrus.Fields{
			"accessToken": accessToken,
			"error":       err.Error(),
		})
		return err
	}

	logger.LogInfo("Successfully blacklisted token", logrus.Fields{
		"accessToken": accessToken,
	})

	return nil
}

// IsBlacklisted checks if the token is blacklisted
func (b *blacklistService) IsBlacklisted(accessToken string) (bool, error) {
	logger.LogInfo("Checking if token is blacklisted", logrus.Fields{
		"accessToken": accessToken,
	})

	// Retrieve all blacklisted tokens from the repository
	blacklists, err := b.blacklistRepository.GetAll()
	if err != nil {
		logger.LogError("Failed to retrieve blacklists", logrus.Fields{
			"error": err.Error(),
		})
		return false, err
	}

	// Check if the access token is in the blacklist
	for _, blacklist := range blacklists {
		if blacklist.AccessToken == accessToken {
			logger.LogInfo("Token is blacklisted", logrus.Fields{
				"accessToken": accessToken,
			})
			return true, nil
		}
	}

	// Token is not blacklisted
	logger.LogInfo("Token is not blacklisted", logrus.Fields{
		"accessToken": accessToken,
	})

	return false, nil
}
