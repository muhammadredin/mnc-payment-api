package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/logger" // Import the logger package
	"PaymentAPI/storage"
	"PaymentAPI/utils"

	"github.com/sirupsen/logrus"
)

type BlacklistRepository interface {
	GetAll() ([]entity.Blacklist, error)
	CreateBlacklist(accessToken string) error
}

type blacklistRepository struct {
	JsonStorage storage.JsonFileHandler[entity.Blacklist]
}

func NewBlacklistRepository(jsonStorage storage.JsonFileHandler[entity.Blacklist]) BlacklistRepository {
	return &blacklistRepository{JsonStorage: jsonStorage}
}

func (r *blacklistRepository) CreateBlacklist(accessToken string) error {
	// Log the start of the operation
	logger.LogInfo("Creating a new blacklist entry", logrus.Fields{
		"operation": "CreateBlacklist",
		"token":     accessToken,
	})

	data, err := r.JsonStorage.ReadFile(constants.BlacklistJsonPath)
	if err != nil {
		logger.LogError("Failed to read blacklist file", logrus.Fields{
			"operation": "CreateBlacklist",
			"error":     err.Error(),
		})
		return err
	}

	expStr, err := utils.GetExpirationFromClaimsAsString(accessToken)
	if err != nil {
		logger.LogError("Failed to get expiration from token claims", logrus.Fields{
			"operation": "CreateBlacklist",
			"token":     accessToken,
			"error":     err.Error(),
		})
		return err
	}

	blacklist := entity.Blacklist{
		AccessToken: accessToken,
		ExpiresAt:   expStr,
	}
	data = append(data, blacklist)

	_, err = r.JsonStorage.WriteFile(data, constants.BlacklistJsonPath)
	if err != nil {
		logger.LogError("Failed to write to blacklist file", logrus.Fields{
			"operation": "CreateBlacklist",
			"error":     err.Error(),
		})
		return err
	}

	// Log success
	logger.LogInfo("Successfully created a blacklist entry", logrus.Fields{
		"operation": "CreateBlacklist",
		"token":     accessToken,
	})
	return nil
}

func (r *blacklistRepository) GetAll() ([]entity.Blacklist, error) {
	// Log the start of the operation
	logger.LogInfo("Fetching all blacklist entries", logrus.Fields{
		"operation": "GetAll",
	})

	data, err := r.JsonStorage.ReadFile(constants.BlacklistJsonPath)
	if err != nil {
		logger.LogError("Failed to fetch blacklist entries", logrus.Fields{
			"operation": "GetAll",
			"error":     err.Error(),
		})
		return nil, err
	}

	// Log success
	logger.LogInfo("Successfully fetched blacklist entries", logrus.Fields{
		"operation": "GetAll",
		"count":     len(data),
	})

	return data, nil
}
