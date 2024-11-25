package service

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"errors"
	"github.com/google/uuid"
	"time"
)

type TransactionService interface {
	CreateNewTransaction(request req.CreateTransactionRequest) (entity.Transaction, error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	walletService         WalletService
}

func NewTransactionService(transactionRepository repository.TransactionRepository, walletService WalletService) TransactionService {
	return &transactionService{transactionRepository, walletService}
}

func (t *transactionService) CreateNewTransaction(request req.CreateTransactionRequest) (entity.Transaction, error) {
	fromWallet, err := t.walletService.GetWalletById(request.FromWalletId)
	if err != nil {
		return entity.Transaction{}, err
	}

	// Check is balance sufficient
	if fromWallet.Balance-request.Amount < 0 {
		return entity.Transaction{}, errors.New(constants.TransactionInsufficientError)
	}

	toWallet, err := t.walletService.GetWalletById(request.ToWalletId)
	if err != nil {
		return entity.Transaction{}, err
	}

	transaction := entity.Transaction{
		Id:           uuid.New().String(),
		FromWalletId: fromWallet.Id,
		ToWalletId:   toWallet.Id,
		CreatedAt:    time.Now().Format(time.RFC3339),
		Amount:       request.Amount,
		Message:      request.Message,
	}

	err = t.transactionRepository.Create(transaction)
	if err != nil {
		return entity.Transaction{}, err
	}

	err = t.walletService.UpdateWallet(fromWallet.Id, request.Amount*-1)
	if err != nil {
		return entity.Transaction{}, err
	}

	err = t.walletService.UpdateWallet(toWallet.Id, request.Amount)
	if err != nil {
		return entity.Transaction{}, err
	}

	return transaction, err
}
