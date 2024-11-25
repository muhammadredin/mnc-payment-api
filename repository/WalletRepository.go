package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

// NewWalletRepository initializes a new instance of WalletRepository.
func NewWalletRepository(jsonStorage storage.JsonFileHandler[entity.Wallet]) WalletRepository {
	return &walletRepository{jsonStorage}
}

// GetAll retrieves all wallets from the JSON storage.
func (w *walletRepository) GetAll() ([]entity.Wallet, error) {
	logrus.Info("Fetching all wallets from storage")
	data, err := w.JsonStorage.ReadFile(constants.WalletJsonPath)
	if err != nil {
		logrus.Errorf("Error reading wallet data: %v", err)
		return nil, err
	}
	return data, nil
}

// GetByCustomerId retrieves a wallet by the customer's ID.
func (w *walletRepository) GetByCustomerId(customerId string) (entity.Wallet, error) {
	logrus.Infof("Fetching wallet for customer ID: %s", customerId)
	data, err := w.GetAll()
	if err != nil {
		return entity.Wallet{}, err
	}

	for _, wallet := range data {
		if wallet.CustomerId == customerId {
			logrus.Infof("Wallet found for customer ID: %s", customerId)
			return wallet, nil
		}
	}

	logrus.Warnf("Wallet not found for customer ID: %s", customerId)
	return entity.Wallet{}, errors.New(constants.WalletNotFoundError)
}

// GetById retrieves a wallet by its ID.
func (w *walletRepository) GetById(id string) (entity.Wallet, error) {
	logrus.Infof("Fetching wallet by ID: %s", id)
	data, err := w.GetAll()
	if err != nil {
		return entity.Wallet{}, err
	}

	for _, wallet := range data {
		if wallet.Id == id {
			logrus.Infof("Wallet found for ID: %s", id)
			return wallet, nil
		}
	}

	logrus.Warnf("Wallet not found for ID: %s", id)
	return entity.Wallet{}, errors.New(constants.WalletNotFoundError)
}

// Create adds a new wallet for the given customer ID.
func (w *walletRepository) Create(customerId string) error {
	logrus.Infof("Creating wallet for customer ID: %s", customerId)

	// Check if wallet already exists for the customer
	_, err := w.GetByCustomerId(customerId)
	if err == nil {
		logrus.Warnf("Wallet already exists for customer ID: %s", customerId)
		return errors.New(constants.WalletDuplicateError)
	}

	// Create a new wallet
	wallet := entity.Wallet{
		Id:         uuid.New().String(),
		CustomerId: customerId,
		Balance:    0,
	}

	data, err := w.JsonStorage.ReadFile(constants.WalletJsonPath)
	if err != nil {
		logrus.Errorf("Error reading wallet data: %v", err)
		return err
	}

	data = append(data, wallet)

	_, err = w.JsonStorage.WriteFile(data, constants.WalletJsonPath)
	if err != nil {
		logrus.Errorf("Error writing new wallet to storage: %v", err)
		return err
	}

	logrus.Infof("Wallet created successfully for customer ID: %s", customerId)
	return nil
}

// Update modifies the balance of an existing wallet.
func (w *walletRepository) Update(id string, balance float64) error {
	logrus.Infof("Updating wallet ID: %s with balance change: %f", id, balance)

	data, err := w.GetAll()
	if err != nil {
		return err
	}

	walletFound := false
	for i := range data {
		if data[i].Id == id {
			data[i].Balance += balance
			walletFound = true
			logrus.Infof("Wallet updated successfully. New balance: %f", data[i].Balance)
			break
		}
	}

	if !walletFound {
		logrus.Warnf("Wallet not found for ID: %s", id)
		return errors.New(constants.WalletNotFoundError)
	}

	_, err = w.JsonStorage.WriteFile(data, constants.WalletJsonPath)
	if err != nil {
		logrus.Errorf("Error writing updated wallet to storage: %v", err)
		return err
	}

	return nil
}
