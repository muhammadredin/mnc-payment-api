package handler

import (
	"PaymentAPI/constants"
	res "PaymentAPI/dto/response"
	"PaymentAPI/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CustomerHandler interface {
	// HandleGetCustomerById retrieves customer details by ID
	HandleGetCustomerById(c *gin.Context)
}

type customerHandler struct {
	customerService service.CustomerService
}

// NewCustomerHandler initializes a new CustomerHandler instance
func NewCustomerHandler(customerService service.CustomerService) CustomerHandler {
	return &customerHandler{customerService}
}

// HandleGetCustomerById handles the request to retrieve a customer by their ID.
// It validates the authenticated user and checks if the user has access to the requested customer.
func (ch *customerHandler) HandleGetCustomerById(c *gin.Context) {
	// Extract the customer ID from the path parameter
	customerId := c.Param("id")
	logrus.Infof("Processing request to get customer by ID: %s", customerId)

	// Retrieve the authenticated user from the context
	user, exists := c.Get("authenticatedUser")
	if !exists {
		logrus.Warn("Authenticated user not found in context")
		c.JSON(http.StatusUnauthorized, res.ErrorResponse{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: constants.AuthenticatedUserNotFoundError,
		})
		return
	}

	// Call the service to fetch the customer by ID
	customer, err := ch.customerService.GetCustomerById(customerId)
	if err != nil {
		logrus.Errorf("Error retrieving customer with ID %s: %v", customerId, err)
		c.JSON(http.StatusNotFound, res.ErrorResponse{
			StatusCode:   http.StatusNotFound,
			ErrorMessage: err.Error(),
		})
		return
	}

	// Validate if the authenticated user has access to the requested customer
	if customer.Id != user {
		logrus.Warnf("Unauthorized access attempt by user: %v to customer ID: %s", user, customerId)
		c.JSON(http.StatusUnauthorized, res.ErrorResponse{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: constants.CustomerForbiddenAccess,
		})
		return
	}

	// Respond with the customer details
	logrus.Infof("Customer retrieved successfully for ID: %s", customerId)
	c.JSON(http.StatusOK, res.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    constants.CustomerFindSuccess,
		Data:       customer,
	})
}
