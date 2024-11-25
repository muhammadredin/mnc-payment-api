package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetAllWallet(t *testing.T) {
	mockJsonFileHandler := new(storage.WalletJsonFileHandlerMock[entity.Wallet])
	walletRepository := NewWalletRepository(mockJsonFileHandler)

	mockJsonFileHandler.Mock.On("ReadFile", constants.WalletJsonPath).
		Return([]entity.Wallet{}, nil)

	data, err := walletRepository.GetAll()
	assert.Nil(t, err)
	assert.NotNil(t, data)
}

func TestGetWalletByCustomerId(t *testing.T) {
	t.Run("ShouldReturnWallet", func(t *testing.T) {
		mockJsonFileHandler := new(storage.WalletJsonFileHandlerMock[entity.Wallet])
		walletRepository := NewWalletRepository(mockJsonFileHandler)

		customerId := "customer-2"
		walletResponse := []entity.Wallet{
			{
				Id:         "wallet-1",
				CustomerId: "customer-1",
				Balance:    0,
			},
			{
				Id:         "wallet-2",
				CustomerId: "customer-2",
				Balance:    0,
			},
		}

		mockJsonFileHandler.Mock.On("ReadFile", constants.WalletJsonPath).
			Return(walletResponse, nil)

		wallet, err := walletRepository.GetByCustomerId(customerId)
		assert.Nil(t, err)
		assert.Equal(t, walletResponse[1], wallet)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockJsonFileHandler := new(storage.WalletJsonFileHandlerMock[entity.Wallet])
		walletRepository := NewWalletRepository(mockJsonFileHandler)

		customerId := "customer-1"
		mockJsonFileHandler.Mock.On("ReadFile", constants.WalletJsonPath).
			Return([]entity.Wallet{}, nil)

		data, err := walletRepository.GetByCustomerId(customerId)
		assert.NotNil(t, constants.WalletNotFoundError, err.Error())
		assert.Equal(t, entity.Wallet{}, data)
	})
}

func TestCreateWallet(t *testing.T) {
	t.Run("ShouldSuccessCreateWallet", func(t *testing.T) {
		mockJsonFileHandler := new(storage.WalletJsonFileHandlerMock[entity.Wallet])
		walletRepository := NewWalletRepository(mockJsonFileHandler)

		customerId := "customer-1"

		mockJsonFileHandler.Mock.On("ReadFile", constants.WalletJsonPath).
			Return([]entity.Wallet{}, nil)

		mockJsonFileHandler.Mock.On("WriteFile", mock.MatchedBy(func(wallets []entity.Wallet) bool {
			if len(wallets) == 0 {
				return false
			}

			wallet := wallets[0]
			return wallet.CustomerId == customerId
		}), constants.WalletJsonPath).
			Return(mock.Anything, nil)

		err := walletRepository.Create(customerId)
		assert.Nil(t, err)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockJsonFileHandler := new(storage.WalletJsonFileHandlerMock[entity.Wallet])
		walletRepository := NewWalletRepository(mockJsonFileHandler)

		customerId := "customer-1"

		walletResponse := []entity.Wallet{
			{
				Id:         "wallet-1",
				CustomerId: "customer-1",
				Balance:    0,
			},
		}

		mockJsonFileHandler.Mock.On("ReadFile", constants.WalletJsonPath).
			Return(walletResponse, nil)

		err := walletRepository.Create(customerId)
		assert.Equal(t, constants.WalletDuplicateError, err.Error())
	})
}

func TestUpdateWallet(t *testing.T) {
	t.Run("ShouldUpdateWallet", func(t *testing.T) {
		mockJsonFileHandler := new(storage.WalletJsonFileHandlerMock[entity.Wallet])
		walletRepository := NewWalletRepository(mockJsonFileHandler)

		walletId := "wallet-1"

		var newBalance float64 = 5000

		walletResponse := []entity.Wallet{
			{
				Id:         "wallet-1",
				CustomerId: "customer-1",
				Balance:    0,
			},
		}

		mockJsonFileHandler.Mock.On("ReadFile", constants.WalletJsonPath).
			Return(walletResponse, nil)

		mockJsonFileHandler.Mock.On("WriteFile", mock.MatchedBy(func(wallets []entity.Wallet) bool {
			if len(wallets) != 1 {
				return false
			}

			wallet := wallets[0]
			return wallet.Balance == newBalance
		}), constants.WalletJsonPath).
			Return(mock.Anything, nil)

		err := walletRepository.Update(walletId, newBalance)
		assert.Nil(t, err)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockJsonFileHandler := new(storage.WalletJsonFileHandlerMock[entity.Wallet])
		walletRepository := NewWalletRepository(mockJsonFileHandler)

		walletId := "wallet-1"

		var newBalance float64 = 5000

		mockJsonFileHandler.Mock.On("ReadFile", constants.WalletJsonPath).
			Return([]entity.Wallet{}, nil)

		err := walletRepository.Update(walletId, newBalance)
		assert.Equal(t, constants.WalletNotFoundError, err.Error())
	})
}
