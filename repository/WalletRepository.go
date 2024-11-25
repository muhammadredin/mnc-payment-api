package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"errors"
	"github.com/google/uuid"
)

type WalletRepository interface {
	GetAll() ([]entity.Wallet, error)
	GetByCustomerId(customerId string) (entity.Wallet, error)
	GetById(id string) (entity.Wallet, error)
	Create(customerId string) error
	Update(id string, balance float64) error
}

type walletRepository struct {
	JsonStorage storage.JsonFileHandler[entity.Wallet]
}

func NewWalletRepository(jsonStorage storage.JsonFileHandler[entity.Wallet]) WalletRepository {
	return &walletRepository{jsonStorage}
}

func (w *walletRepository) GetAll() ([]entity.Wallet, error) {
	data, err := w.JsonStorage.ReadFile(constants.WalletJsonPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *walletRepository) GetByCustomerId(customerId string) (entity.Wallet, error) {
	data, err := w.GetAll()
	if err != nil {
		return entity.Wallet{}, err
	}

	for _, wallet := range data {
		if wallet.CustomerId == customerId {
			return wallet, nil
		}
	}

	return entity.Wallet{}, errors.New(constants.WalletNotFoundError)
}

func (w *walletRepository) GetById(id string) (entity.Wallet, error) {
	data, err := w.GetAll()
	if err != nil {
		return entity.Wallet{}, err
	}

	for _, wallet := range data {
		if wallet.Id == id {
			return wallet, nil
		}
	}

	return entity.Wallet{}, errors.New(constants.WalletNotFoundError)
}

func (w *walletRepository) Create(customerId string) error {
	_, err := w.GetByCustomerId(customerId)
	if err == nil {
		return errors.New(constants.WalletDuplicateError)
	}

	wallet := entity.Wallet{
		Id:         uuid.New().String(),
		CustomerId: customerId,
		Balance:    0,
	}

	data, err := w.JsonStorage.ReadFile(constants.WalletJsonPath)
	if err != nil {
		return err
	}

	data = append(data, wallet)

	_, err = w.JsonStorage.WriteFile(data, constants.WalletJsonPath)
	if err != nil {
		return err
	}

	return nil
}

func (w *walletRepository) Update(id string, balance float64) error {
	data, err := w.GetAll()
	if err != nil {
		return err
	}

	customerFound := false

	for i := range data {
		if data[i].Id == id {
			data[i].Balance += balance
			customerFound = true
			break
		}
	}

	if !customerFound {
		return errors.New(constants.WalletNotFoundError)
	}
	_, err = w.JsonStorage.WriteFile(data, constants.WalletJsonPath)
	if err != nil {
		return err
	}

	return nil
}
