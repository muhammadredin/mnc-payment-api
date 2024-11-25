package service

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/repository"
)

type CustomerService interface {
	GetCustomerByUsername(username string) (res.CustomerResponse, error)
	GetCustomerById(id string) (res.CustomerResponse, error)
	CreateNewCustomer(request req.CreateCustomerRequest) (string, error)
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

	return customer, nil
}

func (c *CustomerServiceImpl) GetCustomerById(id string) (res.CustomerResponse, error) {
	customer, err := c.customerRepository.GetById(id)
	if err != nil {
		return res.CustomerResponse{}, err
	}

	return customer, nil
}

func (c *CustomerServiceImpl) CreateNewCustomer(request req.CreateCustomerRequest) (string, error) {
	customer, err := c.customerRepository.Create(request)
	if err != nil {
		return "", err
	}

	err = c.walletService.CreateWallet(customer.Id)
	if err != nil {
		return "", err
	}

	return constants.CustomerCreateSuccess, nil
}
