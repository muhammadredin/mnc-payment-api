package service

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus" // Import logrus for structured logging
	"time"
)

type TransactionService interface {
	CreateNewTransaction(request req.CreateTransactionRequest) (entity.Transaction, error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	walletService         WalletService
}

// NewTransactionService creates a new instance of TransactionService
func NewTransactionService(transactionRepository repository.TransactionRepository, walletService WalletService) TransactionService {
	return &transactionService{transactionRepository, walletService}
}

// CreateNewTransaction creates a new transaction, transferring funds between wallets
func (t *transactionService) CreateNewTransaction(request req.CreateTransactionRequest) (entity.Transaction, error) {
	logger := logrus.WithFields(logrus.Fields{
		"fromWalletId": request.FromWalletId,
		"toWalletId":   request.ToWalletId,
		"amount":       request.Amount,
	})

	logger.Info("Starting to create a new transaction")

	// Retrieve the 'from' wallet
	fromWallet, err := t.walletService.GetWalletById(request.FromWalletId)
	if err != nil {
		logger.Error("Failed to retrieve 'from' wallet", err)
		return entity.Transaction{}, err
	}

	// Check if the balance is sufficient for the transaction
	if fromWallet.Balance-request.Amount < 0 {
		logger.Error("Insufficient balance for transaction")
		return entity.Transaction{}, errors.New(constants.TransactionInsufficientError)
	}

	// Retrieve the 'to' wallet
	toWallet, err := t.walletService.GetWalletById(request.ToWalletId)
	if err != nil {
		logger.Error("Failed to retrieve 'to' wallet", err)
		return entity.Transaction{}, err
	}

	// Prepare the transaction entity
	transaction := entity.Transaction{
		Id:           uuid.New().String(),
		FromWalletId: fromWallet.Id,
		ToWalletId:   toWallet.Id,
		CreatedAt:    time.Now().Format(time.RFC3339),
		Amount:       request.Amount,
		Message:      request.Message,
	}

	// Create the transaction record in the repository
	err = t.transactionRepository.Create(transaction)
	if err != nil {
		logger.Error("Failed to create transaction in the repository", err)
		return entity.Transaction{}, err
	}

	// Update the 'from' wallet balance
	err = t.walletService.UpdateWallet(fromWallet.Id, request.Amount*-1)
	if err != nil {
		logger.Error("Failed to update 'from' wallet balance", err)
		return entity.Transaction{}, err
	}

	// Update the 'to' wallet balance
	err = t.walletService.UpdateWallet(toWallet.Id, request.Amount)
	if err != nil {
		logger.Error("Failed to update 'to' wallet balance", err)
		return entity.Transaction{}, err
	}

	logger.Info("Transaction successfully created")
	return transaction, nil
}
