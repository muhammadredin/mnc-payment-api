package service

import (
	dto "PaymentAPI/dto/response"
	"PaymentAPI/repository"
)

type CustomerService interface {
	GetCustomerByUsername(username string) (dto.CustomerResponse, error)
	GetCustomerById(id string) (dto.CustomerResponse, error)
}

type CustomerServiceImpl struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(customerRepository repository.CustomerRepository) CustomerService {
	return &CustomerServiceImpl{customerRepository: customerRepository}
}

func (c *CustomerServiceImpl) GetCustomerByUsername(username string) (dto.CustomerResponse, error) {
	customer, err := c.customerRepository.GetByUsername(username)
	if err != nil {
		return dto.CustomerResponse{}, err
	}

	return customer, nil
}

func (c *CustomerServiceImpl) GetCustomerById(id string) (dto.CustomerResponse, error) {
	customer, err := c.customerRepository.GetById(id)
	if err != nil {
		return dto.CustomerResponse{}, err
	}

	return customer, nil
}

//func (c *CustomerServiceImpl) CreateNewCustomer() ([]dto.CustomerResponse, error) {}
