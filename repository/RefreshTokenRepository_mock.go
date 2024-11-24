package repository

import (
	"PaymentAPI/entity"
	"fmt"
	"github.com/stretchr/testify/mock"
)

type RefreshTokenRepositoryMock struct {
	Mock mock.Mock
}

func (r *RefreshTokenRepositoryMock) CreateRefreshToken(customerId string) (entity.RefreshToken, error) {
	args := r.Mock.Called(customerId)
	refreshToken, ok := args.Get(0).(entity.RefreshToken)
	if !ok {
		return entity.RefreshToken{}, fmt.Errorf("invalid type for refresh token")
	}
	return refreshToken, args.Error(1)
}

func (r *RefreshTokenRepositoryMock) GetRefreshToken(refreshToken string) (entity.RefreshToken, error) {
	args := r.Mock.Called(refreshToken)
	return args.Get(0).(entity.RefreshToken), args.Error(1)
}

func (r *RefreshTokenRepositoryMock) GetAllRefreshToken() ([]entity.RefreshToken, error) {
	args := r.Mock.Called()
	return args.Get(0).([]entity.RefreshToken), args.Error(1)
}

func (r *RefreshTokenRepositoryMock) DeleteRefreshToken(refreshToken string) error {
	args := r.Mock.Called(refreshToken)
	return args.Error(0)
}
