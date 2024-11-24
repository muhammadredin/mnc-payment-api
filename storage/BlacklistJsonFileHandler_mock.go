package storage

import (
	"PaymentAPI/entity"
	"github.com/stretchr/testify/mock"
)

type BlacklistJsonFileHandlerMock[T entity.Blacklist] struct {
	Mock mock.Mock
}

func (j *BlacklistJsonFileHandlerMock[T]) ReadFile(path string) ([]T, error) {
	arguments := j.Mock.Called(path)

	if args := arguments.Get(0); args != nil {
		return args.([]T), arguments.Error(1)
	}
	return nil, arguments.Error(1)
}

func (j *BlacklistJsonFileHandlerMock[T]) WriteFile(data []T, path string) (string, error) {
	arguments := j.Mock.Called(data, path)

	return arguments.String(0), arguments.Error(1)
}
