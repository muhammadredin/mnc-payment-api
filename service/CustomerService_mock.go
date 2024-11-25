package service

import (
	req "PaymentAPI/dto/request"
	dto "PaymentAPI/dto/response"
	"github.com/stretchr/testify/mock"
)

type CustomerServiceMock struct {
	Mock mock.Mock
}

func (c *CustomerServiceMock) GetCustomerByUsername(username string) (dto.CustomerResponse, error) {
	args := c.Mock.Called(username)
	return args.Get(0).(dto.CustomerResponse), args.Error(1)
}

func (c *CustomerServiceMock) GetCustomerById(id string) (dto.CustomerResponse, error) {
	args := c.Mock.Called(id)
	return args.Get(0).(dto.CustomerResponse), args.Error(1)
}

func (c *CustomerServiceMock) CreateNewCustomer(request req.CreateCustomerRequest) (string, error) {
	args := c.Mock.Called(request)
	return args.String(0), args.Error(1)
}
