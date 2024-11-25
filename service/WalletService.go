package service

import (
	"PaymentAPI/entity"
	"PaymentAPI/repository"
)

type WalletService interface {
	CreateWallet(customerId string) error
	GetWallet(customerId string) (entity.Wallet, error)
	UpdateWallet(customerId string, balance float64) error
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

func (w *walletService) GetWallet(customerId string) (entity.Wallet, error) {
	wallet, err := w.WalletRepository.GetByCustomerId(customerId)
	if err != nil {
		return entity.Wallet{}, err
	}

	return wallet, nil
}

func (w *walletService) UpdateWallet(customerId string, balance float64) error {
	err := w.WalletRepository.Update(customerId, balance)
	if err != nil {
		return err
	}
	return nil
}
