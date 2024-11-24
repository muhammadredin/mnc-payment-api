package service

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"errors"
	"fmt"
	"time"
)

type RefreshTokenService interface {
	GenerateRefreshToken(customerId string) (entity.RefreshToken, error)
	RotateRefreshToken(refreshToken string) (entity.RefreshToken, error)
}

type refreshTokenService struct {
	refreshTokenRepository repository.RefreshTokenRepository
}

func NewRefreshTokenService(refreshTokenRepository repository.RefreshTokenRepository) RefreshTokenService {
	return &refreshTokenService{refreshTokenRepository: refreshTokenRepository}
}

func (r *refreshTokenService) GenerateRefreshToken(customerId string) (entity.RefreshToken, error) {
	refreshToken, err := r.refreshTokenRepository.GetAllRefreshToken()
	if err != nil {
		return entity.RefreshToken{}, err
	}

	for _, token := range refreshToken {
		if token.CustomerId == customerId {
			r.refreshTokenRepository.DeleteRefreshToken(token.RefreshToken)
			break
		}
	}

	token, err := r.refreshTokenRepository.CreateRefreshToken(customerId)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return token, nil
}

func (r *refreshTokenService) RotateRefreshToken(refreshToken string) (entity.RefreshToken, error) {
	token, err := r.refreshTokenRepository.GetRefreshToken(refreshToken)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	parsedTime, err := time.Parse(time.RFC3339, token.ExpiresAt)
	fmt.Println(parsedTime)

	r.refreshTokenRepository.DeleteRefreshToken(token.RefreshToken)
	if time.Now().After(parsedTime) {
		return entity.RefreshToken{}, errors.New(constants.RefreshTokenExpiredError)
	}

	newRefreshToken, err := r.GenerateRefreshToken(token.CustomerId)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return newRefreshToken, nil
}
