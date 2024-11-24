package service

import (
	"github.com/stretchr/testify/mock"
)

type RefreshTokenServiceMock struct {
	mock.Mock
}

func (m *RefreshTokenServiceMock) GenerateRefreshToken(customerId string) (string, error) {
	args := m.Called(customerId)
	return args.String(0), args.Error(1)
}

func (m *RefreshTokenServiceMock) RotateRefreshToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}
