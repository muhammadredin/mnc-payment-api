package service

import (
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
