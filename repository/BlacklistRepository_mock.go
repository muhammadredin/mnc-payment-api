package repository

import (
	"PaymentAPI/entity"
	"github.com/stretchr/testify/mock"
)

type BlacklistRepositoryMock struct {
	Mock mock.Mock
}

func (b *BlacklistRepositoryMock) CreateBlacklist(accessToken string) error {
	args := b.Mock.Called(accessToken)
	return args.Error(0)
}

func (b *BlacklistRepositoryMock) GetAll() ([]entity.Blacklist, error) {
	args := b.Mock.Called()
	return args.Get(0).([]entity.Blacklist), args.Error(1)
}
