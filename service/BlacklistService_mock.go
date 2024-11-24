package service

import (
	"github.com/stretchr/testify/mock"
)

type BlacklistServiceMock struct {
	mock.Mock
}

func (b *BlacklistServiceMock) BlacklistToken(accessToken string) error {
	args := b.Called(accessToken)
	return args.Error(0)
}

func (b *BlacklistServiceMock) IsBlacklisted(accessToken string) (bool, error) {
	args := b.Called(accessToken)
	return args.Bool(0), args.Error(1)
}
