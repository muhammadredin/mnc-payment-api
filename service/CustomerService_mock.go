package service

import (
	req "PaymentAPI/dto/request"
	dto "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"github.com/stretchr/testify/mock"
)

type CustomerServiceMock struct {
	Mock mock.Mock
}

func (c *CustomerServiceMock) GetCustomerByUsername(username string) (dto.CustomerResponse, error) {
	args := c.Mock.Called(username)
	return args.Get(0).(dto.CustomerResponse), args.Error(1)
}

func (c *CustomerServiceMock) GetCustomerByUsernameAuth(username string) (entity.Customer, error) {
	args := c.Mock.Called(username)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (c *CustomerServiceMock) GetCustomerById(id string) (dto.CustomerResponse, error) {
	args := c.Mock.Called(id)
	return args.Get(0).(dto.CustomerResponse), args.Error(1)
}

func (c *CustomerServiceMock) GetCustomerByIdAuth(id string) (entity.Customer, error) {
	args := c.Mock.Called(id)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (c *CustomerServiceMock) CreateNewCustomer(request req.CustomerRequest) (string, error) {
	args := c.Mock.Called(request)
	return args.String(0), args.Error(1)
}
