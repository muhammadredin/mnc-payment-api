package service

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"PaymentAPI/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus" // Import the logrus package for structured logging
)

type CustomerService interface {
	GetCustomerByUsername(username string) (res.CustomerResponse, error)
	GetCustomerByUsernameAuth(username string) (entity.Customer, error)
	GetCustomerById(id string) (res.CustomerResponse, error)
	GetCustomerByIdAuth(id string) (entity.Customer, error)
	CreateNewCustomer(request req.CustomerRequest) (string, error)
}

type CustomerServiceImpl struct {
	customerRepository repository.CustomerRepository
	walletService      WalletService
}

// NewCustomerService creates a new instance of CustomerServiceImpl
func NewCustomerService(customerRepository repository.CustomerRepository, walletService WalletService) CustomerService {
	return &CustomerServiceImpl{customerRepository: customerRepository, walletService: walletService}
}

// GetCustomerByUsername retrieves a customer by username and maps to response
func (c *CustomerServiceImpl) GetCustomerByUsername(username string) (res.CustomerResponse, error) {
	logger := logrus.WithFields(logrus.Fields{"username": username})
	logger.Info("Fetching customer by username")

	// Retrieve customer from the repository
	customer, err := c.customerRepository.GetByUsername(username)
	if err != nil {
		logger.Error("Failed to retrieve customer by username", err)
		return res.CustomerResponse{}, err
	}

	// Fetch wallet information for the customer
	wallet, err := c.walletService.GetWalletByCustomerId(customer.Id)
	if err != nil {
		logger.Error("Failed to retrieve wallet for customer", err)
		return res.CustomerResponse{}, err
	}

	// Map the customer and wallet details to response
	logger.Info("Successfully fetched customer and wallet")
	return mapCustomerToCustomerResponse(customer, wallet), nil
}

// GetCustomerByUsernameAuth retrieves a customer by username for authentication
func (c *CustomerServiceImpl) GetCustomerByUsernameAuth(username string) (entity.Customer, error) {
	logger := logrus.WithFields(logrus.Fields{"username": username})
	logger.Info("Fetching customer by username for authentication")

	// Retrieve customer from repository
	customer, err := c.customerRepository.GetByUsername(username)
	if err != nil {
		logger.Error("Failed to retrieve customer by username for authentication", err)
		return entity.Customer{}, err
	}

	logger.Info("Successfully fetched customer for authentication")
	return customer, nil
}

// GetCustomerById retrieves a customer by ID and maps to response
func (c *CustomerServiceImpl) GetCustomerById(id string) (res.CustomerResponse, error) {
	logger := logrus.WithFields(logrus.Fields{"customerId": id})
	logger.Info("Fetching customer by ID")

	// Retrieve customer from the repository
	customer, err := c.customerRepository.GetById(id)
	if err != nil {
		logger.Error("Failed to retrieve customer by ID", err)
		return res.CustomerResponse{}, err
	}

	// Fetch wallet information for the customer
	wallet, err := c.walletService.GetWalletByCustomerId(customer.Id)
	if err != nil {
		logger.Error("Failed to retrieve wallet for customer", err)
		return res.CustomerResponse{}, err
	}

	// Map the customer and wallet details to response
	logger.Info("Successfully fetched customer and wallet")
	return mapCustomerToCustomerResponse(customer, wallet), nil
}

// GetCustomerByIdAuth retrieves a customer by ID for authentication
func (c *CustomerServiceImpl) GetCustomerByIdAuth(id string) (entity.Customer, error) {
	logger := logrus.WithFields(logrus.Fields{"customerId": id})
	logger.Info("Fetching customer by ID for authentication")

	// Retrieve customer from repository
	customer, err := c.customerRepository.GetById(id)
	if err != nil {
		logger.Error("Failed to retrieve customer by ID for authentication", err)
		return entity.Customer{}, err
	}

	logger.Info("Successfully fetched customer for authentication")
	return customer, nil
}

// CreateNewCustomer creates a new customer and wallet
func (c *CustomerServiceImpl) CreateNewCustomer(request req.CustomerRequest) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"username": request.Username,
	})
	logger.Info("Creating new customer")

	// Map the customer request to customer entity
	customerRequest := mapCreateCustomerToCustomer(request)

	// Create the customer in the repository
	customer, err := c.customerRepository.Create(customerRequest)
	if err != nil {
		logger.Error("Failed to create new customer", err)
		return "", err
	}

	// Create wallet for the new customer
	err = c.walletService.CreateWallet(customer.Id)
	if err != nil {
		logger.Error("Failed to create wallet for new customer", err)
		return "", err
	}

	logger.Info("Successfully created new customer and wallet")
	return constants.CustomerCreateSuccess, nil
}

// mapCreateCustomerToCustomer maps the customer request to a customer entity
func mapCreateCustomerToCustomer(request req.CustomerRequest) entity.Customer {
	// Encrypt the password with Bcrypt
	encryptedPassword := utils.BCryptEncoder(request.Password)

	// Map to customer entity
	return entity.Customer{
		Id:       uuid.New().String(),
		Username: request.Username,
		Password: encryptedPassword,
	}
}

// mapCustomerToCustomerResponse maps the customer and wallet details to response format
func mapCustomerToCustomerResponse(customer entity.Customer, wallet entity.Wallet) res.CustomerResponse {
	return res.CustomerResponse{
		Id:       customer.Id,
		Username: customer.Username,
		WalletId: wallet.Id,
		Balance:  wallet.Balance,
	}
}
