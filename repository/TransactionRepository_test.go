package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreateTransaction(t *testing.T) {
	mockFileHandler := new(storage.TransactionJsonFileHandlerMock[entity.Transaction])
	transactionRepository := NewTransactionRepository(mockFileHandler)

	newTransaction := entity.Transaction{
		Id:           "transaction-1",
		FromWalletId: "wallet-1",
		ToWalletId:   "wallet-2",
		Amount:       5000,
		CreatedAt:    time.Now().Format(time.RFC3339),
		Message:      "transaction",
	}

	mockFileHandler.Mock.On("ReadFile", constants.TransactionJsonPath).
		Return([]entity.Transaction{}, nil)

	mockFileHandler.Mock.On("WriteFile",
		mock.MatchedBy(func(transactions []entity.Transaction) bool {
			// Check if the slice has exactly one element
			if len(transactions) != 1 {
				return false
			}

			transaction := transactions[0]
			if transaction.Id != newTransaction.Id {
				return false
			}
			return true
		}),
		constants.TransactionJsonPath,
	).Return(constants.JsonWriteSuccess, nil)

	err := transactionRepository.Create(newTransaction)

	assert.Nil(t, err)
}

func TestGetAllTransactions(t *testing.T) {
	mockFileHandler := new(storage.TransactionJsonFileHandlerMock[entity.Transaction])
	transactionRepository := NewTransactionRepository(mockFileHandler)

	transaction := entity.Transaction{
		Id:           "transaction-1",
		FromWalletId: "wallet-1",
		ToWalletId:   "wallet-2",
		Amount:       5000,
		CreatedAt:    time.Now().Format(time.RFC3339),
		Message:      "transaction",
	}

	mockFileHandler.Mock.On("ReadFile", constants.TransactionJsonPath).
		Return([]entity.Transaction{transaction}, nil)

	transactions, err := transactionRepository.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(transactions))
}
