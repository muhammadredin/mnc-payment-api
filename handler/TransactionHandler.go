package handler

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type TransactionHandler interface {
	HandleCreateTransaction(c *gin.Context)
}

type transactionHandler struct {
	transactionService service.TransactionService
	walletService      service.WalletService
}

// NewTransactionHandler creates a new instance of TransactionHandler.
func NewTransactionHandler(transactionService service.TransactionService, walletService service.WalletService) TransactionHandler {
	return &transactionHandler{transactionService, walletService}
}

// HandleCreateTransaction handles the request to create a new transaction.
func (t transactionHandler) HandleCreateTransaction(c *gin.Context) {
	var request req.CreateTransactionRequest

	// Bind JSON request body to struct and handle errors
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Warn("Invalid request body for transaction creation")
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: constants.InvalidRequestBodyError,
		})
		return
	}

	// Retrieve authenticated user from the context
	user, exists := c.Get("authenticatedUser")
	if !exists {
		logrus.Warn("Authenticated user not found in context")
		c.JSON(http.StatusUnauthorized, res.ErrorResponse{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: constants.AuthenticatedUserNotFoundError,
		})
		return
	}

	logrus.Infof("Authenticated user: %v is creating a transaction from wallet ID: %s", user, request.FromWalletId)

	// Fetch wallet details and validate ownership
	wallet, err := t.walletService.GetWalletById(request.FromWalletId)
	if err != nil {
		logrus.Errorf("Failed to fetch wallet with ID: %s, error: %v", request.FromWalletId, err)
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		})
		return
	}

	// Ensure the wallet belongs to the authenticated user
	if wallet.CustomerId != user {
		logrus.Warnf("User %v attempted unauthorized access to wallet ID: %s", user, request.FromWalletId)
		c.JSON(http.StatusForbidden, res.ErrorResponse{
			StatusCode:   http.StatusForbidden,
			ErrorMessage: constants.WalletForbiddenAccess,
		})
		return
	}

	logrus.Infof("Creating transaction for wallet ID: %s by user: %v", request.FromWalletId, user)

	// Create the transaction and handle potential errors
	response, err := t.transactionService.CreateNewTransaction(request)
	if err != nil {
		logrus.Errorf("Failed to create transaction, error: %v", err)
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: err.Error(),
		})
		return
	}

	// Log and send successful response
	logrus.Infof("Transaction successfully created: %v", response)
	c.JSON(http.StatusCreated, res.CommonResponse{
		StatusCode: http.StatusCreated,
		Message:    constants.TransactionSuccess,
		Data:       response,
	})
	return
}
