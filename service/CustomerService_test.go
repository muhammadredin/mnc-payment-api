package service

import (
	"PaymentAPI/constants"
	dto "PaymentAPI/dto/response"
	"PaymentAPI/repository"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetCustomerByUsername(t *testing.T) {
	customer := dto.CustomerResponse{
		Id:       "id-1",
		Username: "johndoe",
	}

	t.Run("ShouldReturnCustomer", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		customerService := NewCustomerService(mockCustomerRepository)

		mockCustomerRepository.Mock.On("GetByUsername", "johndoe").
			Return(customer, nil)

		result, err := customerService.GetCustomerByUsername("johndoe")
		assert.Nil(t, err)
		assert.Equal(t, customer, result)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		customerService := NewCustomerService(mockCustomerRepository)

		mockCustomerRepository.Mock.On("GetByUsername", mock.Anything).
			Return(dto.CustomerResponse{}, errors.New(constants.CustomerNotFound))

		result, err := customerService.GetCustomerByUsername("")
		assert.NotNil(t, err)
		assert.Equal(t, dto.CustomerResponse{}, result)
	})
}

func TestGetCustomerById(t *testing.T) {
	customer := dto.CustomerResponse{
		Id:       "id-1",
		Username: "johndoe",
	}

	t.Run("ShouldReturnCustomer", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		customerService := NewCustomerService(mockCustomerRepository)

		mockCustomerRepository.Mock.On("GetById", customer.Id).
			Return(customer, nil)

		result, err := customerService.GetCustomerById(customer.Id)
		assert.Nil(t, err)
		assert.Equal(t, customer, result)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		customerService := NewCustomerService(mockCustomerRepository)

		mockCustomerRepository.Mock.On("GetById", mock.Anything).
			Return(dto.CustomerResponse{}, errors.New(constants.CustomerNotFound))

		result, err := customerService.GetCustomerById("")
		assert.NotNil(t, err)
		assert.Equal(t, dto.CustomerResponse{}, result)
	})
}
