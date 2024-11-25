package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"errors"
)

type CustomerRepository interface {
	GetByUsername(id string) (entity.Customer, error)
	GetById(id string) (entity.Customer, error)
	Create(customer entity.Customer) (entity.Customer, error)
}

type customerRepository struct {
	JsonStorage storage.JsonFileHandler[entity.Customer]
}

func NewCustomerRepository(jsonStorage storage.JsonFileHandler[entity.Customer]) CustomerRepository {
	return &customerRepository{
		JsonStorage: jsonStorage,
	}
}

func (cr *customerRepository) Create(customer entity.Customer) (entity.Customer, error) {
	// Get the previous json file
	data, err := cr.JsonStorage.ReadFile(constants.CustomerJsonPath)
	if err != nil {
		return entity.Customer{}, err
	}

	// Check if there's duplicate username
	for _, item := range data {
		if customer.Username == item.Username {
			return entity.Customer{}, errors.New(constants.UsernameDuplicateError)
		}
	}

	// Append the json file
	data = append(data, customer)

	// Save the json file
	_, err = cr.JsonStorage.WriteFile(data, constants.CustomerJsonPath)
	if err != nil {
		return entity.Customer{}, err
	}

	return customer, nil
}

func (cr *customerRepository) GetByUsername(username string) (entity.Customer, error) {
	data, err := cr.JsonStorage.ReadFile(constants.CustomerJsonPath)
	if err != nil {
		return entity.Customer{}, err
	}

	for _, c := range data {
		if c.Username == username {
			return c, nil
		}
	}

	return entity.Customer{}, errors.New(constants.CustomerNotFound)
}

func (cr *customerRepository) GetById(id string) (entity.Customer, error) {
	data, err := cr.JsonStorage.ReadFile(constants.CustomerJsonPath)
	if err != nil {
		return entity.Customer{}, err
	}

	for _, c := range data {
		if c.Id == id {
			return c, nil
		}
	}

	return entity.Customer{}, errors.New(constants.CustomerNotFound)
}
