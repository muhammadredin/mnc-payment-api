package service

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"PaymentAPI/utils"
	"fmt"
	"github.com/google/uuid"
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

func NewCustomerService(customerRepository repository.CustomerRepository, walletService WalletService) CustomerService {
	return &CustomerServiceImpl{customerRepository: customerRepository, walletService: walletService}
}

func (c *CustomerServiceImpl) GetCustomerByUsername(username string) (res.CustomerResponse, error) {
	customer, err := c.customerRepository.GetByUsername(username)
	if err != nil {
		return res.CustomerResponse{}, err
	}

	return mapCustomerToCustomerResponse(customer), nil
}

func (c *CustomerServiceImpl) GetCustomerByUsernameAuth(username string) (entity.Customer, error) {
	customer, err := c.customerRepository.GetByUsername(username)
	if err != nil {
		return entity.Customer{}, err
	}

	return customer, nil
}

func (c *CustomerServiceImpl) GetCustomerById(id string) (res.CustomerResponse, error) {
	customer, err := c.customerRepository.GetById(id)
	if err != nil {
		return res.CustomerResponse{}, err
	}

	return mapCustomerToCustomerResponse(customer), nil
}

func (c *CustomerServiceImpl) GetCustomerByIdAuth(id string) (entity.Customer, error) {
	customer, err := c.customerRepository.GetById(id)
	if err != nil {
		return entity.Customer{}, err
	}

	return customer, nil
}

func (c *CustomerServiceImpl) CreateNewCustomer(request req.CustomerRequest) (string, error) {
	customerRequest := mapCreateCustomerToCustomer(request)
	customer, err := c.customerRepository.Create(customerRequest)
	fmt.Println(err)
	if err != nil {
		return "", err
	}

	err = c.walletService.CreateWallet(customer.Id)
	if err != nil {
		return "", err
	}

	return constants.CustomerCreateSuccess, nil
}

func mapCreateCustomerToCustomer(request req.CustomerRequest) entity.Customer {
	// Encrypt the password with Bcrypt
	encryptedPassword := utils.BCryptEncoder(request.Password)

	// Map to customer entity
	customer := entity.Customer{
		Id:       uuid.New().String(),
		Username: request.Username,
		Password: encryptedPassword,
	}

	return customer
}

func mapCustomerToCustomerResponse(customer entity.Customer) res.CustomerResponse {
	return res.CustomerResponse{
		Id:       customer.Id,
		Username: customer.Username,
	}
}
