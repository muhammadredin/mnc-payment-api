package service

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type refreshTokenTest struct {
	mockRepo *repository.RefreshTokenRepositoryMock
	service  RefreshTokenService
}

func setupRefreshTokenTest() refreshTokenTest {
	mockRepo := new(repository.RefreshTokenRepositoryMock)
	service := NewRefreshTokenService(mockRepo)
	return refreshTokenTest{
		mockRepo: mockRepo,
		service:  service,
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tests := []struct {
		name           string
		customerId     string
		existingTokens []entity.RefreshToken
		newToken       string
		expectedError  error
	}{
		{
			name:       "Should Generate New Token When User Has Existing Token",
			customerId: "user-1",
			existingTokens: []entity.RefreshToken{
				{
					RefreshToken: "refresh-token-1",
					CustomerId:   "user-1",
					ExpiresAt:    time.Now().Add(time.Hour).Format(time.RFC3339),
				},
			},
			newToken:      "new-refresh-token",
			expectedError: nil,
		},
		{
			name:           "Should Generate New Token When User Has No Existing Token",
			customerId:     "user-2",
			existingTokens: []entity.RefreshToken{},
			newToken:       "new-refresh-token",
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			test := setupRefreshTokenTest()

			// Mock GetAllRefreshToken
			test.mockRepo.Mock.On("GetAllRefreshToken").
				Return(tt.existingTokens, nil)

			// Mock DeleteRefreshToken if there are existing tokens
			if len(tt.existingTokens) > 0 {
				test.mockRepo.Mock.On("DeleteRefreshToken", tt.existingTokens[0].RefreshToken).
					Return(nil)
			}

			// Mock CreateRefreshToken
			test.mockRepo.Mock.On("CreateRefreshToken", tt.customerId).
				Return(tt.newToken, nil)

			// Execute
			token, err := test.service.GenerateRefreshToken(tt.customerId)

			// Assert
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.newToken, token)
			}
			test.mockRepo.Mock.AssertExpectations(t)
		})
	}
}

func TestRotateRefreshToken(t *testing.T) {
	tests := []struct {
		name          string
		refreshToken  string
		storedToken   entity.RefreshToken
		newToken      string
		expectedError error
	}{
		{
			name:         "Should Rotate Valid Token",
			refreshToken: "refresh-token-1",
			storedToken: entity.RefreshToken{
				RefreshToken: "refresh-token-1",
				CustomerId:   "user-1",
				ExpiresAt:    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			newToken:      "new-refresh-token",
			expectedError: nil,
		},
		{
			name:         "Should Fail For Expired Token",
			refreshToken: "expired-token",
			storedToken: entity.RefreshToken{
				RefreshToken: "expired-token",
				CustomerId:   "user-1",
				ExpiresAt:    time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			},
			newToken:      "",
			expectedError: errors.New(constants.RefreshTokenExpiredError),
		},
		{
			name:          "Should Fail For Non-existent Token",
			refreshToken:  "non-existent-token",
			storedToken:   entity.RefreshToken{},
			newToken:      "",
			expectedError: errors.New(constants.RefreshTokenNotFoundError),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			test := setupRefreshTokenTest()

			// Mock GetRefreshToken
			test.mockRepo.Mock.On("GetRefreshToken", tt.refreshToken).
				Return(tt.storedToken, tt.expectedError)

			if tt.expectedError == nil {
				// Mock DeleteRefreshToken
				test.mockRepo.Mock.On("DeleteRefreshToken", tt.refreshToken).
					Return(nil)

				// Mock GetAllRefreshToken (called by GenerateRefreshToken)
				test.mockRepo.Mock.On("GetAllRefreshToken").
					Return([]entity.RefreshToken{}, nil)

				// Mock CreateRefreshToken (called by GenerateRefreshToken)
				test.mockRepo.Mock.On("CreateRefreshToken", tt.storedToken.CustomerId).
					Return(tt.newToken, nil)
			}

			// Execute
			token, err := test.service.RotateRefreshToken(tt.refreshToken)

			// Assert
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, token)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.newToken, token)
				assert.NotEqual(t, tt.refreshToken, token)
			}
			test.mockRepo.Mock.AssertExpectations(t)
		})
	}
}
