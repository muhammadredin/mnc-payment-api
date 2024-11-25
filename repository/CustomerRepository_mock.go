package repository

import (
	"PaymentAPI/entity"
	"github.com/stretchr/testify/mock"
)

type CustomerRepositoryMock struct {
	Mock mock.Mock
}

func (c *CustomerRepositoryMock) Create(customer entity.Customer) (entity.Customer, error) {
	args := c.Mock.Called(customer)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (c *CustomerRepositoryMock) GetByUsername(username string) (entity.Customer, error) {
	args := c.Mock.Called(username)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (c *CustomerRepositoryMock) GetById(id string) (entity.Customer, error) {
	args := c.Mock.Called(id)
	return args.Get(0).(entity.Customer), args.Error(1)
}
