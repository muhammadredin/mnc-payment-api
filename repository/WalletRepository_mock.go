package repository

import (
	"PaymentAPI/entity"
	"fmt"
	"github.com/stretchr/testify/mock"
)

type WalletRepositoryMock struct {
	Mock mock.Mock
}

func (w *WalletRepositoryMock) GetAll() ([]entity.Wallet, error) {
	args := w.Mock.Called()
	refreshToken, ok := args.Get(0).([]entity.Wallet)
	if !ok {
		return nil, fmt.Errorf("invalid type for wallet")
	}
	return refreshToken, args.Error(1)
}

func (w *WalletRepositoryMock) GetByCustomerId(customerId string) (entity.Wallet, error) {
	args := w.Mock.Called(customerId)
	return args.Get(0).(entity.Wallet), args.Error(1)
}

func (w *WalletRepositoryMock) GetById(id string) (entity.Wallet, error) {
	args := w.Mock.Called(id)
	return args.Get(0).(entity.Wallet), args.Error(1)
}

func (w *WalletRepositoryMock) Create(customerId string) error {
	args := w.Mock.Called(customerId)
	return args.Error(0)
}

func (w *WalletRepositoryMock) Update(customerId string, balance float64) error {
	args := w.Mock.Called(customerId, balance)
	return args.Error(0)
}
