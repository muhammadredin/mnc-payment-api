package service

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/logger" // Import the logger package
	"PaymentAPI/utils"
	"errors"

	"github.com/sirupsen/logrus"
)

type AuthService interface {
	Login(request req.CustomerRequest) (res.AuthResponse, error)
	Logout(accessToken string, refreshToken string) error
	GetNewAccessToken(refreshToken string) (res.AuthResponse, error)
}

type authService struct {
	customerService     CustomerService
	refreshTokenService RefreshTokenService
	blacklistService    BlacklistService
}

// Constructor for AuthService
func NewAuthService(customerService CustomerService, refreshTokenService RefreshTokenService, blacklistService BlacklistService) AuthService {
	return &authService{
		customerService:     customerService,
		refreshTokenService: refreshTokenService,
		blacklistService:    blacklistService,
	}
}

// Login handles user authentication and token generation
func (a *authService) Login(request req.CustomerRequest) (res.AuthResponse, error) {
	logger.LogInfo("Attempting to log in", logrus.Fields{
		"username": request.Username,
	})

	// Fetch customer details by username
	customer, err := a.customerService.GetCustomerByUsernameAuth(request.Username)
	if err != nil {
		logger.LogError("Failed to fetch customer by username", logrus.Fields{
			"username": request.Username,
			"error":    err.Error(),
		})
		return res.AuthResponse{}, err
	}

	// Compare provided password with stored hashed password
	if !utils.BCryptCompare(request.Password, []byte(customer.Password)) {
		logger.LogError("Unauthorized login attempt", logrus.Fields{
			"username": request.Username,
		})
		return res.AuthResponse{}, errors.New(constants.LoginUnauthorizedError)
	}

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(customer)
	if err != nil {
		logger.LogError("Failed to generate access token", logrus.Fields{
			"username": request.Username,
			"error":    err.Error(),
		})
		return res.AuthResponse{}, err
	}

	// Generate refresh token
	refreshToken, err := a.refreshTokenService.GenerateRefreshToken(customer.Id)
	if err != nil {
		logger.LogError("Failed to generate refresh token", logrus.Fields{
			"username": request.Username,
			"error":    err.Error(),
		})
		return res.AuthResponse{}, err
	}

	logger.LogInfo("Successfully logged in", logrus.Fields{
		"username":   request.Username,
		"customerId": customer.Id,
	})

	// Return authentication response with tokens
	return res.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.RefreshToken,
		CustomerId:   customer.Id,
	}, nil
}

// Logout invalidates the user's access and refresh tokens
func (a *authService) Logout(accessToken string, refreshToken string) error {
	logger.LogInfo("Attempting to log out user", logrus.Fields{
		"accessToken": accessToken,
	})

	// Add access token to the blacklist
	err := a.blacklistService.BlacklistToken(accessToken)
	if err != nil {
		logger.LogError("Failed to blacklist access token", logrus.Fields{
			"accessToken": accessToken,
			"error":       err.Error(),
		})
		return err
	}

	// Delete the refresh token
	err = a.refreshTokenService.DeleteRefreshToken(refreshToken)
	if err != nil {
		logger.LogError("Failed to delete refresh token", logrus.Fields{
			"refreshToken": refreshToken,
			"error":        err.Error(),
		})
		return err
	}

	logger.LogInfo("Successfully logged out user", logrus.Fields{
		"accessToken": accessToken,
	})
	return nil
}

// GetNewAccessToken generates a new access token and rotates the refresh token
func (a *authService) GetNewAccessToken(refreshToken string) (res.AuthResponse, error) {
	logger.LogInfo("Attempting to rotate refresh token and generate a new access token", logrus.Fields{
		"refreshToken": refreshToken,
	})

	// Rotate refresh token
	newRefreshToken, err := a.refreshTokenService.RotateRefreshToken(refreshToken)
	if err != nil {
		logger.LogError("Failed to rotate refresh token", logrus.Fields{
			"refreshToken": refreshToken,
			"error":        err.Error(),
		})
		return res.AuthResponse{}, err
	}

	// Fetch customer details by ID
	customer, err := a.customerService.GetCustomerByIdAuth(newRefreshToken.CustomerId)
	if err != nil {
		logger.LogError("Failed to fetch customer by ID", logrus.Fields{
			"customerId": newRefreshToken.CustomerId,
			"error":      err.Error(),
		})
		return res.AuthResponse{}, err
	}

	// Generate a new access token
	accessToken, err := utils.GenerateAccessToken(customer)
	if err != nil {
		logger.LogError("Failed to generate access token", logrus.Fields{
			"customerId": customer.Id,
			"error":      err.Error(),
		})
		return res.AuthResponse{}, err
	}

	logger.LogInfo("Successfully generated new access token", logrus.Fields{
		"customerId": customer.Id,
	})

	// Return authentication response with new tokens
	return res.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.RefreshToken,
		CustomerId:   customer.Id,
	}, nil
}
