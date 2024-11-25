package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
)

type TransactionRepository interface {
	GetAll() ([]entity.Transaction, error)
	Create(entity.Transaction) error
}

type transactionRepository struct {
	JsonStorage storage.JsonFileHandler[entity.Transaction]
}

func NewTransactionRepository(jsonStorage storage.JsonFileHandler[entity.Transaction]) TransactionRepository {
	return &transactionRepository{jsonStorage}
}

func (w *transactionRepository) GetAll() ([]entity.Transaction, error) {
	data, err := w.JsonStorage.ReadFile(constants.TransactionJsonPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *transactionRepository) Create(transaction entity.Transaction) error {
	data, err := w.JsonStorage.ReadFile(constants.TransactionJsonPath)
	if err != nil {
		return err
	}

	data = append(data, transaction)

	_, err = w.JsonStorage.WriteFile(data, constants.TransactionJsonPath)
	if err != nil {
		return err
	}

	return nil
}
