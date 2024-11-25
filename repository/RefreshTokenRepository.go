package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus" // Import logrus for structured logging
	"time"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(customerId string) (entity.RefreshToken, error)
	GetRefreshToken(refreshToken string) (entity.RefreshToken, error)
	GetAllRefreshToken() ([]entity.RefreshToken, error)
	DeleteRefreshToken(refreshToken string) error
}

type refreshTokenRepository struct {
	JsonStorage storage.JsonFileHandler[entity.RefreshToken]
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository
func NewRefreshTokenRepository(jsonStorage storage.JsonFileHandler[entity.RefreshToken]) RefreshTokenRepository {
	return &refreshTokenRepository{JsonStorage: jsonStorage}
}

// CreateRefreshToken creates a new refresh token for a given customerId
func (r *refreshTokenRepository) CreateRefreshToken(customerId string) (entity.RefreshToken, error) {
	logger := logrus.WithFields(logrus.Fields{
		"customerId": customerId,
	})

	logger.Info("Creating new refresh token")

	// Create new refresh token
	refreshToken := entity.RefreshToken{
		RefreshToken: uuid.New().String(),
		CustomerId:   customerId,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	// Read existing refresh tokens from file
	data, err := r.JsonStorage.ReadFile(constants.RefreshTokenJsonPath)
	if err != nil {
		logger.Error("Failed to read refresh tokens file", err)
		return entity.RefreshToken{}, err
	}

	// Append the new token to the existing data
	data = append(data, refreshToken)

	// Write the updated data back to the file
	_, err = r.JsonStorage.WriteFile(data, constants.RefreshTokenJsonPath)
	if err != nil {
		logger.Error("Failed to write updated refresh tokens file", err)
		return entity.RefreshToken{}, err
	}

	logger.Info("Refresh token created successfully")
	return refreshToken, nil
}

// GetRefreshToken retrieves a refresh token by its value
func (r *refreshTokenRepository) GetRefreshToken(refreshToken string) (entity.RefreshToken, error) {
	logger := logrus.WithFields(logrus.Fields{
		"refreshToken": refreshToken,
	})

	logger.Info("Retrieving refresh token")

	// Read the refresh tokens from file
	data, err := r.JsonStorage.ReadFile(constants.RefreshTokenJsonPath)
	if err != nil {
		logger.Error("Failed to read refresh tokens file", err)
		return entity.RefreshToken{}, err
	}

	// Search for the matching refresh token
	for _, token := range data {
		if token.RefreshToken == refreshToken {
			logger.Info("Refresh token found")
			return token, nil
		}
	}

	logger.Error("Refresh token not found")
	return entity.RefreshToken{}, errors.New(constants.RefreshTokenNotFoundError)
}

// GetAllRefreshToken retrieves all refresh tokens
func (r *refreshTokenRepository) GetAllRefreshToken() ([]entity.RefreshToken, error) {
	logger := logrus.WithFields(logrus.Fields{})

	logger.Info("Retrieving all refresh tokens")

	// Read the refresh tokens from file
	data, err := r.JsonStorage.ReadFile(constants.RefreshTokenJsonPath)
	if err != nil {
		logger.Error("Failed to read refresh tokens file", err)
		return nil, err
	}

	logger.Info("All refresh tokens retrieved successfully")
	return data, nil
}

// DeleteRefreshToken deletes a refresh token by its value
func (r *refreshTokenRepository) DeleteRefreshToken(refreshToken string) error {
	logger := logrus.WithFields(logrus.Fields{
		"refreshToken": refreshToken,
	})

	logger.Info("Deleting refresh token")

	// Read the refresh tokens from file
	data, err := r.JsonStorage.ReadFile(constants.RefreshTokenJsonPath)
	if err != nil {
		logger.Error("Failed to read refresh tokens file", err)
		return err
	}

	// Search for the token and delete it if found
	indexToDelete := -1
	for i, token := range data {
		if token.RefreshToken == refreshToken {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		logger.Error("Refresh token not found")
		return errors.New(constants.RefreshTokenNotFoundError)
	}

	// Remove the token from the list
	data = append(data[:indexToDelete], data[indexToDelete+1:]...)

	// Write the updated list back to the file
	_, err = r.JsonStorage.WriteFile(data, constants.RefreshTokenJsonPath)
	if err != nil {
		logger.Error("Failed to write updated refresh tokens file", err)
		return err
	}

	logger.Info("Refresh token deleted successfully")
	return nil
}
