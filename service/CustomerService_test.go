package service

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/repository"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetCustomerByUsername(t *testing.T) {
	customer := entity.Customer{
		Id:       "id-1",
		Username: "johndoe",
		Password: "password",
	}

	expectedResponse := res.CustomerResponse{
		Id:       "id-1",
		Username: "johndoe",
	}

	t.Run("ShouldReturnCustomer", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		mockCustomerRepository.Mock.On("GetByUsername", "johndoe").
			Return(customer, nil)

		result, err := customerService.GetCustomerByUsername("johndoe")
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, result)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		mockCustomerRepository.Mock.On("GetByUsername", mock.Anything).
			Return(entity.Customer{}, errors.New(constants.CustomerNotFound))

		result, err := customerService.GetCustomerByUsername("")
		assert.NotNil(t, err)
		assert.Equal(t, res.CustomerResponse{}, result)
	})
}

func TestGetCustomerByUsernameAuth(t *testing.T) {
	customer := entity.Customer{
		Id:       "id-1",
		Username: "johndoe",
		Password: "password",
	}

	t.Run("ShouldReturnCustomer", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		mockCustomerRepository.Mock.On("GetByUsername", "johndoe").
			Return(customer, nil)

		result, err := customerService.GetCustomerByUsernameAuth("johndoe")
		assert.Nil(t, err)
		assert.Equal(t, customer, result)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		mockCustomerRepository.Mock.On("GetByUsername", mock.Anything).
			Return(entity.Customer{}, errors.New(constants.CustomerNotFound))

		result, err := customerService.GetCustomerByUsernameAuth("")
		assert.NotNil(t, err)
		assert.Equal(t, entity.Customer{}, result)
	})
}

func TestGetCustomerById(t *testing.T) {
	customer := entity.Customer{
		Id:       "id-1",
		Username: "johndoe",
		Password: "password",
	}

	expectedResponse := res.CustomerResponse{
		Id:       "id-1",
		Username: "johndoe",
	}

	t.Run("ShouldReturnCustomer", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		mockCustomerRepository.Mock.On("GetById", customer.Id).
			Return(customer, nil)

		result, err := customerService.GetCustomerById(customer.Id)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, result)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		mockCustomerRepository.Mock.On("GetById", mock.Anything).
			Return(entity.Customer{}, errors.New(constants.CustomerNotFound))

		result, err := customerService.GetCustomerById("")
		assert.NotNil(t, err)
		assert.Equal(t, res.CustomerResponse{}, result)
	})
}

func TestGetCustomerByIdAuth(t *testing.T) {
	customer := entity.Customer{
		Id:       "id-1",
		Username: "johndoe",
		Password: "password",
	}

	t.Run("ShouldReturnCustomer", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		mockCustomerRepository.Mock.On("GetById", customer.Id).
			Return(customer, nil)

		result, err := customerService.GetCustomerByIdAuth(customer.Id)
		assert.Nil(t, err)
		assert.Equal(t, customer, result)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		mockCustomerRepository.Mock.On("GetById", mock.Anything).
			Return(entity.Customer{}, errors.New(constants.CustomerNotFound))

		result, err := customerService.GetCustomerByIdAuth("")
		assert.NotNil(t, err)
		assert.Equal(t, entity.Customer{}, result)
	})
}

func TestCreateNewCustomer(t *testing.T) {
	t.Run("ShouldCreateCustomer", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		request := req.CustomerRequest{
			Username: "customer-1",
			Password: "password",
		}

		mappedRequest := entity.Customer{
			Id:       "id-1",
			Username: "customer-1",
			Password: "password",
		}

		mockCustomerRepository.Mock.On("Create", mock.Anything).
			Return(mappedRequest, nil)

		mockWalletService.Mock.On("CreateWallet", mappedRequest.Id).
			Return(nil)

		result, err := customerService.CreateNewCustomer(request)
		assert.Nil(t, err)
		assert.Equal(t, constants.CustomerCreateSuccess, result)
	})

	t.Run("ShouldReturnErrorOnCustomer", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		request := req.CustomerRequest{
			Username: "customer-1",
			Password: "password",
		}

		mockCustomerRepository.Mock.On("Create", mock.Anything).
			Return(entity.Customer{}, errors.New(constants.UsernameDuplicateError))

		result, err := customerService.CreateNewCustomer(request)
		assert.Equal(t, constants.UsernameDuplicateError, err.Error())
		assert.Equal(t, "", result)
	})

	t.Run("ShouldReturnErrorOnWallet", func(t *testing.T) {
		mockCustomerRepository := new(repository.CustomerRepositoryMock)
		mockWalletService := new(WalletServiceMock)
		customerService := NewCustomerService(mockCustomerRepository, mockWalletService)

		request := req.CustomerRequest{
			Username: "customer-1",
			Password: "password",
		}

		mappedRequest := entity.Customer{
			Id:       "id-1",
			Username: "customer-1",
			Password: "password",
		}

		response := res.CustomerResponse{
			Id:       "id-1",
			Username: "customer-1",
		}

		mockCustomerRepository.Mock.On("Create", mock.Anything).
			Return(mappedRequest, nil)

		mockWalletService.Mock.On("CreateWallet", response.Id).
			Return(errors.New(constants.WalletDuplicateError))

		result, err := customerService.CreateNewCustomer(request)
		assert.Equal(t, constants.WalletDuplicateError, err.Error())
		assert.Equal(t, "", result)
	})
}
