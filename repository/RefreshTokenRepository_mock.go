package repository

import (
	"PaymentAPI/entity"
	"github.com/stretchr/testify/mock"
)

type RefreshTokenRepositoryMock struct {
	Mock mock.Mock
}

func (r *RefreshTokenRepositoryMock) CreateRefreshToken(customerId string) (string, error) {
	args := r.Mock.Called(customerId)
	return args.String(0), args.Error(1)
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
