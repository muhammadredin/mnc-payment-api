package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"github.com/sirupsen/logrus" // Importing logrus for structured logging
)

type TransactionRepository interface {
	GetAll() ([]entity.Transaction, error)
	Create(entity.Transaction) error
}

type transactionRepository struct {
	JsonStorage storage.JsonFileHandler[entity.Transaction]
}

// NewTransactionRepository creates a new instance of TransactionRepository
func NewTransactionRepository(jsonStorage storage.JsonFileHandler[entity.Transaction]) TransactionRepository {
	return &transactionRepository{JsonStorage: jsonStorage}
}

// GetAll retrieves all transactions from storage
func (t *transactionRepository) GetAll() ([]entity.Transaction, error) {
	logger := logrus.WithFields(logrus.Fields{})

	logger.Info("Retrieving all transactions")

	// Read the transactions from storage
	data, err := t.JsonStorage.ReadFile(constants.TransactionJsonPath)
	if err != nil {
		logger.Error("Failed to read transactions file", err)
		return nil, err
	}

	logger.Info("All transactions retrieved successfully")
	return data, nil
}

// Create adds a new transaction to storage
func (t *transactionRepository) Create(transaction entity.Transaction) error {
	logger := logrus.WithFields(logrus.Fields{
		"transactionId": transaction.Id,
		"fromWalletId":  transaction.FromWalletId,
		"toWalletId":    transaction.ToWalletId,
		"amount":        transaction.Amount,
	})

	logger.Info("Creating new transaction")

	// Read the current transactions from storage
	data, err := t.JsonStorage.ReadFile(constants.TransactionJsonPath)
	if err != nil {
		logger.Error("Failed to read transactions file", err)
		return err
	}

	// Append the new transaction to the data
	data = append(data, transaction)

	// Write the updated data back to storage
	_, err = t.JsonStorage.WriteFile(data, constants.TransactionJsonPath)
	if err != nil {
		logger.Error("Failed to write updated transactions file", err)
		return err
	}

	logger.Info("New transaction created successfully")
	return nil
}
