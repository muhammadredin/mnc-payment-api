package service

import (
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"github.com/sirupsen/logrus" // Import logrus for structured logging
)

type WalletService interface {
	CreateWallet(customerId string) error
	GetWalletByCustomerId(customerId string) (entity.Wallet, error)
	GetWalletById(id string) (entity.Wallet, error)
	UpdateWallet(id string, balance float64) error
}

type walletService struct {
	WalletRepository repository.WalletRepository
}

// NewWalletService creates a new instance of WalletService
func NewWalletService(walletRepository repository.WalletRepository) WalletService {
	return &walletService{WalletRepository: walletRepository}
}

// CreateWallet creates a new wallet for a customer
func (w *walletService) CreateWallet(customerId string) error {
	logger := logrus.WithFields(logrus.Fields{
		"customerId": customerId,
	})

	logger.Info("Creating new wallet for customer")

	// Attempt to create a new wallet
	err := w.WalletRepository.Create(customerId)
	if err != nil {
		logger.Error("Failed to create wallet", err)
		return err
	}

	logger.Info("New wallet created successfully")
	return nil
}

// GetWalletByCustomerId retrieves a wallet by customer ID
func (w *walletService) GetWalletByCustomerId(customerId string) (entity.Wallet, error) {
	logger := logrus.WithFields(logrus.Fields{
		"customerId": customerId,
	})

	logger.Info("Retrieving wallet for customer")

	// Attempt to retrieve the wallet
	wallet, err := w.WalletRepository.GetByCustomerId(customerId)
	if err != nil {
		logger.Error("Failed to retrieve wallet by customer ID", err)
		return entity.Wallet{}, err
	}

	logger.Info("Wallet retrieved successfully")
	return wallet, nil
}

// GetWalletById retrieves a wallet by wallet ID
func (w *walletService) GetWalletById(id string) (entity.Wallet, error) {
	logger := logrus.WithFields(logrus.Fields{
		"walletId": id,
	})

	logger.Info("Retrieving wallet by ID")

	// Attempt to retrieve the wallet
	wallet, err := w.WalletRepository.GetById(id)
	if err != nil {
		logger.Error("Failed to retrieve wallet by ID", err)
		return entity.Wallet{}, err
	}

	logger.Info("Wallet retrieved successfully")
	return wallet, nil
}

// UpdateWallet updates the balance of a wallet
func (w *walletService) UpdateWallet(id string, balance float64) error {
	logger := logrus.WithFields(logrus.Fields{
		"walletId": id,
		"balance":  balance,
	})

	logger.Info("Updating wallet balance")

	// Attempt to update the wallet balance
	err := w.WalletRepository.Update(id, balance)
	if err != nil {
		logger.Error("Failed to update wallet balance", err)
		return err
	}

	logger.Info("Wallet balance updated successfully")
	return nil
}
