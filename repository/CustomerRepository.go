package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/logger" // Import the logger package
	"PaymentAPI/storage"
	"errors"

	"github.com/sirupsen/logrus"
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
	logger.LogInfo("Starting to create customer", logrus.Fields{
		"operation": "Create",
		"username":  customer.Username,
		"id":        customer.Id,
	})

	// Get the previous JSON file
	data, err := cr.JsonStorage.ReadFile(constants.CustomerJsonPath)
	if err != nil {
		logger.LogError("Failed to read customer file", logrus.Fields{
			"operation": "Create",
			"error":     err.Error(),
		})
		return entity.Customer{}, err
	}

	// Check for duplicate username
	for _, item := range data {
		if customer.Username == item.Username {
			logger.LogError("Duplicate username found", logrus.Fields{
				"operation": "Create",
				"username":  customer.Username,
			})
			return entity.Customer{}, errors.New(constants.UsernameDuplicateError)
		}
	}

	// Append the JSON file
	data = append(data, customer)

	// Save the JSON file
	_, err = cr.JsonStorage.WriteFile(data, constants.CustomerJsonPath)
	if err != nil {
		logger.LogError("Failed to write customer file", logrus.Fields{
			"operation": "Create",
			"error":     err.Error(),
		})
		return entity.Customer{}, err
	}

	logger.LogInfo("Successfully created customer", logrus.Fields{
		"operation": "Create",
		"username":  customer.Username,
		"id":        customer.Id,
	})

	return customer, nil
}

func (cr *customerRepository) GetByUsername(username string) (entity.Customer, error) {
	logger.LogInfo("Fetching customer by username", logrus.Fields{
		"operation": "GetByUsername",
		"username":  username,
	})

	data, err := cr.JsonStorage.ReadFile(constants.CustomerJsonPath)
	if err != nil {
		logger.LogError("Failed to read customer file", logrus.Fields{
			"operation": "GetByUsername",
			"error":     err.Error(),
		})
		return entity.Customer{}, err
	}

	for _, c := range data {
		if c.Username == username {
			logger.LogInfo("Successfully fetched customer", logrus.Fields{
				"operation": "GetByUsername",
				"username":  username,
				"id":        c.Id,
			})
			return c, nil
		}
	}

	logger.LogError("Customer not found", logrus.Fields{
		"operation": "GetByUsername",
		"username":  username,
	})
	return entity.Customer{}, errors.New(constants.CustomerNotFound)
}

func (cr *customerRepository) GetById(id string) (entity.Customer, error) {
	logger.LogInfo("Fetching customer by ID", logrus.Fields{
		"operation": "GetById",
		"id":        id,
	})

	data, err := cr.JsonStorage.ReadFile(constants.CustomerJsonPath)
	if err != nil {
		logger.LogError("Failed to read customer file", logrus.Fields{
			"operation": "GetById",
			"error":     err.Error(),
		})
		return entity.Customer{}, err
	}

	for _, c := range data {
		if c.Id == id {
			logger.LogInfo("Successfully fetched customer", logrus.Fields{
				"operation": "GetById",
				"username":  c.Username,
				"id":        c.Id,
			})
			return c, nil
		}
	}

	logger.LogError("Customer not found", logrus.Fields{
		"operation": "GetById",
		"id":        id,
	})
	return entity.Customer{}, errors.New(constants.CustomerNotFound)
}
