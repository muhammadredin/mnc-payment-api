package service

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"errors"
	"github.com/sirupsen/logrus" // Import logrus for structured logging
	"time"
)

type RefreshTokenService interface {
	GenerateRefreshToken(customerId string) (entity.RefreshToken, error)
	RotateRefreshToken(refreshToken string) (entity.RefreshToken, error)
	DeleteRefreshToken(refreshToken string) error
}

type refreshTokenService struct {
	refreshTokenRepository repository.RefreshTokenRepository
}

// NewRefreshTokenService creates a new instance of RefreshTokenService
func NewRefreshTokenService(refreshTokenRepository repository.RefreshTokenRepository) RefreshTokenService {
	return &refreshTokenService{refreshTokenRepository: refreshTokenRepository}
}

// GenerateRefreshToken generates a new refresh token for a customer
func (r *refreshTokenService) GenerateRefreshToken(customerId string) (entity.RefreshToken, error) {
	logger := logrus.WithFields(logrus.Fields{
		"customerId": customerId,
	})
	logger.Info("Generating refresh token")

	// Retrieve all refresh tokens to check if there's an existing one for the customer
	refreshTokens, err := r.refreshTokenRepository.GetAllRefreshToken()
	if err != nil {
		logger.Error("Failed to retrieve refresh tokens", err)
		return entity.RefreshToken{}, err
	}

	// Remove the existing refresh token for the customer if present
	for _, token := range refreshTokens {
		if token.CustomerId == customerId {
			logger.Info("Deleting existing refresh token for customer")
			r.refreshTokenRepository.DeleteRefreshToken(token.RefreshToken)
			break
		}
	}

	// Create a new refresh token
	token, err := r.refreshTokenRepository.CreateRefreshToken(customerId)
	if err != nil {
		logger.Error("Failed to create new refresh token", err)
		return entity.RefreshToken{}, err
	}

	logger.Info("Successfully generated new refresh token")
	return token, nil
}

// RotateRefreshToken rotates the refresh token by verifying and generating a new one
func (r *refreshTokenService) RotateRefreshToken(refreshToken string) (entity.RefreshToken, error) {
	logger := logrus.WithFields(logrus.Fields{
		"refreshToken": refreshToken,
	})
	logger.Info("Rotating refresh token")

	// Retrieve the existing refresh token
	token, err := r.refreshTokenRepository.GetRefreshToken(refreshToken)
	if err != nil {
		logger.Error("Failed to retrieve refresh token", err)
		return entity.RefreshToken{}, err
	}

	// Parse the expiration time from the refresh token
	parsedTime, err := time.Parse(time.RFC3339, token.ExpiresAt)
	if err != nil {
		logger.Error("Failed to parse expiration time", err)
		return entity.RefreshToken{}, err
	}

	// If the refresh token is expired, return an error
	if time.Now().After(parsedTime) {
		logger.Error("Refresh token has expired")
		return entity.RefreshToken{}, errors.New(constants.RefreshTokenExpiredError)
	}

	// Delete the expired refresh token
	r.refreshTokenRepository.DeleteRefreshToken(token.RefreshToken)

	// Generate a new refresh token for the customer
	newRefreshToken, err := r.GenerateRefreshToken(token.CustomerId)
	if err != nil {
		logger.Error("Failed to generate new refresh token", err)
		return entity.RefreshToken{}, err
	}

	logger.Info("Successfully rotated refresh token")
	return newRefreshToken, nil
}

// DeleteRefreshToken deletes a refresh token
func (r *refreshTokenService) DeleteRefreshToken(refreshToken string) error {
	logger := logrus.WithFields(logrus.Fields{
		"refreshToken": refreshToken,
	})
	logger.Info("Deleting refresh token")

	// Delete the refresh token from the repository
	err := r.refreshTokenRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		logger.Error("Failed to delete refresh token", err)
		return err
	}

	logger.Info("Successfully deleted refresh token")
	return nil
}
