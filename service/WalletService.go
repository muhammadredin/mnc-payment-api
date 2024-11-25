package service

import (
	"PaymentAPI/entity"
	"PaymentAPI/repository"
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

func NewWalletService(walletRepository repository.WalletRepository) WalletService {
	return &walletService{WalletRepository: walletRepository}
}

func (w *walletService) CreateWallet(customerId string) error {
	err := w.WalletRepository.Create(customerId)
	if err != nil {
		return err
	}
	return nil
}

func (w *walletService) GetWalletByCustomerId(customerId string) (entity.Wallet, error) {
	wallet, err := w.WalletRepository.GetByCustomerId(customerId)
	if err != nil {
		return entity.Wallet{}, err
	}

	return wallet, nil
}

func (w *walletService) GetWalletById(id string) (entity.Wallet, error) {
	wallet, err := w.WalletRepository.GetById(id)
	if err != nil {
		return entity.Wallet{}, err
	}

	return wallet, nil
}

func (w *walletService) UpdateWallet(id string, balance float64) error {
	err := w.WalletRepository.Update(id, balance)
	if err != nil {
		return err
	}
	return nil
}
