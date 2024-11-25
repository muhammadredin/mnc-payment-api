package service

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	t.Run("ShouldCreateWallet", func(t *testing.T) {
		mockWalletRepository := new(repository.WalletRepositoryMock)
		walletService := NewWalletService(mockWalletRepository)

		customerId := "customer-1"

		mockWalletRepository.Mock.On("Create", customerId).
			Return(nil)

		err := walletService.CreateWallet(customerId)
		assert.Nil(t, err)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockWalletRepository := new(repository.WalletRepositoryMock)
		walletService := NewWalletService(mockWalletRepository)

		customerId := "customer-1"
		mockWalletRepository.Mock.On("Create", customerId).
			Return(errors.New(constants.WalletDuplicateError))

		err := walletService.CreateWallet(customerId)
		assert.NotNil(t, err)
	})
}

func TestGetWallet(t *testing.T) {
	t.Run("ShouldGetWallet", func(t *testing.T) {
		mockWalletRepository := new(repository.WalletRepositoryMock)
		walletService := NewWalletService(mockWalletRepository)

		customerId := "customer-1"

		walletResponse := entity.Wallet{
			Id:         "id-1",
			CustomerId: "customer-1",
			Balance:    0,
		}

		mockWalletRepository.Mock.On("GetByCustomerId", customerId).
			Return(walletResponse, nil)

		wallet, err := walletService.GetWalletByCustomerId(customerId)
		assert.Nil(t, err)
		assert.Equal(t, customerId, wallet.CustomerId)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockWalletRepository := new(repository.WalletRepositoryMock)
		walletService := NewWalletService(mockWalletRepository)

		customerId := "customer-1"

		mockWalletRepository.Mock.On("GetByCustomerId", customerId).
			Return(entity.Wallet{}, errors.New(constants.WalletNotFoundError))

		wallet, err := walletService.GetWalletByCustomerId(customerId)
		assert.Equal(t, constants.WalletNotFoundError, err.Error())
		assert.Equal(t, entity.Wallet{}, wallet)
	})
}

func TestUpdateWallet(t *testing.T) {
	t.Run("ShouldUpdateWallet", func(t *testing.T) {
		mockWalletRepository := new(repository.WalletRepositoryMock)
		walletService := NewWalletService(mockWalletRepository)

		walletId := "wallet-1"
		var balance float64 = 5000

		mockWalletRepository.Mock.On("Update", walletId, balance).
			Return(nil)

		err := walletService.UpdateWallet(walletId, balance)
		assert.Nil(t, err)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockWalletRepository := new(repository.WalletRepositoryMock)
		walletService := NewWalletService(mockWalletRepository)

		walletId := "wallet-1"
		var balance float64 = 5000

		mockWalletRepository.Mock.On("Update", walletId, balance).
			Return(errors.New(constants.WalletNotFoundError))

		err := walletService.UpdateWallet(walletId, balance)
		assert.Equal(t, constants.WalletNotFoundError, err.Error())
	})
}
