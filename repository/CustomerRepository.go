package repository

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"PaymentAPI/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	GetByUsername(id string) (res.CustomerResponse, error)
	GetById(id string) (res.CustomerResponse, error)
	Create(request req.CreateCustomerRequest) (res.CustomerResponse, error)
}

type customerRepository struct {
	JsonStorage storage.JsonFileHandler[entity.Customer]
}

func NewCustomerRepository(jsonStorage storage.JsonFileHandler[entity.Customer]) CustomerRepository {
	return &customerRepository{
		JsonStorage: jsonStorage,
	}
}

func (cr *customerRepository) Create(request req.CreateCustomerRequest) (res.CustomerResponse, error) {
	// Create customer entity by calling mapper function
	customer := MapCreateCustomerToCustomer(request)

	// Get the previous json file
	data, err := cr.JsonStorage.ReadFile(constants.CustomerJsonPath)
	if err != nil {
		return res.CustomerResponse{}, err
	}
	fmt.Println(data)

	// Check if there's duplicate username
	for _, item := range data {
		if customer.Username == item.Username {
			return res.CustomerResponse{}, errors.New(constants.UsernameDuplicateError)
		}
	}

	// Append the json file
	data = append(data, customer)

	// Save the json file
	_, err = cr.JsonStorage.WriteFile(data, constants.CustomerJsonPath)
	if err != nil {
		return res.CustomerResponse{}, err
	}

	return mapCustomerToCustomerResponse(customer), nil
}

func (cr *customerRepository) GetByUsername(username string) (res.CustomerResponse, error) {
	data, err := cr.JsonStorage.ReadFile(constants.CustomerJsonPath)
	if err != nil {
		return res.CustomerResponse{}, err
	}

	for _, c := range data {
		if c.Username == username {
			return mapCustomerToCustomerResponse(c), nil
		}
	}

	return res.CustomerResponse{}, errors.New(constants.CustomerNotFound)
}

func (cr *customerRepository) GetById(id string) (res.CustomerResponse, error) {
	data, err := cr.JsonStorage.ReadFile(constants.CustomerJsonPath)
	if err != nil {
		return res.CustomerResponse{}, err
	}

	for _, c := range data {
		if c.Id == id {
			return mapCustomerToCustomerResponse(c), nil
		}
	}

	return res.CustomerResponse{}, errors.New(constants.CustomerNotFound)
}

func MapCreateCustomerToCustomer(request req.CreateCustomerRequest) entity.Customer {
	// Encrypt the password with Bcrypt
	encryptedPassword := utils.BCryptEncoder(request.Password)

	// Map to customer entity
	customer := entity.Customer{
		Id:       uuid.New().String(),
		Username: request.Username,
		Password: encryptedPassword,
	}

	return customer
}

func mapCustomerToCustomerResponse(customer entity.Customer) res.CustomerResponse {
	return res.CustomerResponse{
		Id:       customer.Id,
		Username: customer.Username,
	}
}
