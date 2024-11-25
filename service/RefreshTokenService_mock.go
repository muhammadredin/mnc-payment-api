package service

import (
	"PaymentAPI/entity"
	"fmt"
	"github.com/stretchr/testify/mock"
)

type RefreshTokenServiceMock struct {
	mock.Mock
}

func (m *RefreshTokenServiceMock) GenerateRefreshToken(customerId string) (entity.RefreshToken, error) {
	args := m.Called(customerId)

	refreshToken, ok := args.Get(0).(entity.RefreshToken)
	if !ok {
		return entity.RefreshToken{}, fmt.Errorf("invalid type for refresh token")
	}
	return refreshToken, args.Error(1)
}

func (m *RefreshTokenServiceMock) RotateRefreshToken(refreshToken string) (entity.RefreshToken, error) {
	args := m.Called(refreshToken)

	newRefreshToken, ok := args.Get(0).(entity.RefreshToken)
	if !ok {
		return entity.RefreshToken{}, fmt.Errorf("invalid type for refresh token")
	}
	return newRefreshToken, args.Error(1)
}

func (m *RefreshTokenServiceMock) DeleteRefreshToken(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}
