package repository

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateCustomer(t *testing.T) {
	customer := req.CreateCustomerRequest{
		Username: "budi",
		Password: "budi123",
	}

	t.Run("ShouldReturnSuccessOnCreate", func(t *testing.T) {
		mockFileHandler := new(storage.CustomerJsonFileHandlerMock[entity.Customer])
		customerRepository := NewCustomerRepository(mockFileHandler)
		mockFileHandler.Mock.On("ReadFile", constants.CustomerJsonPath).
			Return([]entity.Customer{}, nil)

		mockFileHandler.Mock.On("WriteFile",
			mock.MatchedBy(func(customers []entity.Customer) bool {
				if len(customers) != 1 {
					return false
				}
				customer := customers[0]
				return customer.Username == "budi"
			}),
			constants.CustomerJsonPath).
			Return(constants.JsonWriteSuccess, nil)

		response, err := customerRepository.Create(customer)

		assert.Equal(t, constants.CustomerCreateSuccess, response)
		assert.Nil(t, err)
		mockFileHandler.Mock.AssertExpectations(t)
	})

	t.Run("ShouldReturnDuplicateErrorOnCreate", func(t *testing.T) {
		mockFileHandler := new(storage.CustomerJsonFileHandlerMock[entity.Customer])
		customerRepository := NewCustomerRepository(mockFileHandler)

		existingCustomers := []entity.Customer{
			{
				Id:       "id-1",
				Username: "budi",
				Password: "budi123",
			},
		}

		mockFileHandler.Mock.On("ReadFile", constants.CustomerJsonPath).
			Return(existingCustomers, nil)

		response, err := customerRepository.Create(customer)

		assert.NotNil(t, err)
		assert.Equal(t, "", response)
		mockFileHandler.Mock.AssertExpectations(t)
	})
}

func TestGetCustomerById(t *testing.T) {
	mockFileHandler := new(storage.CustomerJsonFileHandlerMock[entity.Customer])
	customerRepository := NewCustomerRepository(mockFileHandler)
	customer := entity.Customer{
		Id:       "customer-1",
		Username: "customer-1",
		Password: "customer-1",
	}

	customerResponse := res.CustomerResponse{
		Id:       "customer-1",
		Username: "customer-1",
	}

	mockFileHandler.Mock.On("ReadFile", mock.Anything).Return([]entity.Customer{customer}, nil) // Mock ReadFile to return customer when called
	mockFileHandler.Mock.On("ReadFile", "").Return([]entity.Customer{}, nil)

	t.Run("ShouldReturnCustomer", func(t *testing.T) {
		response, err := customerRepository.GetByUsername(customer.Id)
		assert.Equal(t, response, customerResponse, "Response not equal to customer")
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		response, err := customerRepository.GetByUsername("")
		assert.Equal(t, response, res.CustomerResponse{}, "Response should be empty struct")
		assert.Equal(t, constants.CustomerNotFound, err.Error(), "Error message not correct")
	})
}
