package repository

import (
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"github.com/stretchr/testify/mock"
)

type CustomerRepositoryMock struct {
	Mock mock.Mock
}

func (c *CustomerRepositoryMock) Create(request req.CreateCustomerRequest) (res.CustomerResponse, error) {
	args := c.Mock.Called(request)
	return args.Get(0).(res.CustomerResponse), args.Error(1)
}

func (c *CustomerRepositoryMock) GetByUsername(username string) (res.CustomerResponse, error) {
	args := c.Mock.Called(username)
	return args.Get(0).(res.CustomerResponse), args.Error(1)
}

func (c *CustomerRepositoryMock) GetById(id string) (res.CustomerResponse, error) {
	args := c.Mock.Called(id)
	return args.Get(0).(res.CustomerResponse), args.Error(1)
}
