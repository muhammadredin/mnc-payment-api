package handler

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionHandler interface {
	HandleCreateTransaction(c *gin.Context)
}

type transactionHandler struct {
	transactionService service.TransactionService
	walletService      service.WalletService
}

func NewTransactionHandler(transactionService service.TransactionService, walletService service.WalletService) TransactionHandler {
	return &transactionHandler{transactionService, walletService}
}

func (t transactionHandler) HandleCreateTransaction(c *gin.Context) {
	var request req.CreateTransactionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: constants.InvalidRequestBodyError,
		})
		return
	}

	user, exists := c.Get("authenticatedUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: constants.AuthenticatedUserNotFoundError,
		})
	}

	// If logged in customer Id not equals to wallet customer Id then return error
	wallet, err := t.walletService.GetWalletById(request.FromWalletId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{})
		return
	}

	if wallet.CustomerId != user {
		c.JSON(http.StatusForbidden, res.ErrorResponse{
			StatusCode:   http.StatusForbidden,
			ErrorMessage: constants.WalletForbiddenAccess,
		})
		return
	}

	// Create new transaction
	response, err := t.transactionService.CreateNewTransaction(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, res.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    constants.TransactionSuccess,
		Data:       response,
	})
	return
}
