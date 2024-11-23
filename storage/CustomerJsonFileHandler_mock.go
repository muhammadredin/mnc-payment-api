package storage

import (
	"PaymentAPI/entity"
	"github.com/stretchr/testify/mock"
)

type CustomerJsonFileHandlerMock[T entity.Customer] struct {
	Mock mock.Mock
}

func (j *CustomerJsonFileHandlerMock[T]) ReadFile(path string) ([]T, error) {
	arguments := j.Mock.Called(path)

	// Safely cast the return value to the generic type slice []T
	if args := arguments.Get(0); args != nil {
		return args.([]T), arguments.Error(1)
	}
	return nil, arguments.Error(1)
}

func (j *CustomerJsonFileHandlerMock[T]) WriteFile(data []T, path string) (string, error) {
	arguments := j.Mock.Called(data, path)

	// Return mocked response
	return arguments.String(0), arguments.Error(1)
}
