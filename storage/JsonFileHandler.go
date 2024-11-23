package storage

import (
	"PaymentAPI/constants"
	"encoding/json"
	"errors"
	"os"
)

type JsonFileHandler[T any] interface {
	ReadFile(path string) ([]T, error)
	WriteFile(data []T, path string) (string, error)
}

type JsonFileHandlerImpl[T any] struct{}

func NewJsonFileHandler[T any]() JsonFileHandler[T] {
	return &JsonFileHandlerImpl[T]{}
}

func (j *JsonFileHandlerImpl[T]) ReadFile(path string) ([]T, error) {
	// Open Json File
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New(constants.JsonFileNotFound)
	}
	defer file.Close()

	// Map the json file to golang slice
	var data []T
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, errors.New(constants.JsonMappingError)
	}

	return data, nil
}

func (j *JsonFileHandlerImpl[T]) WriteFile(data []T, path string) (string, error) {
	// Rewrite the file
	file, err := os.Create(path)
	if err != nil {
		return constants.JsonCreateError, err
	}
	defer file.Close()

	// Encode golang slice to json string
	updatedData, err := json.Marshal(data)
	if err != nil {
		return constants.JsonMarshalError, err
	}

	_, err = file.Write(updatedData)
	if err != nil {
		return constants.JsonWriteError, err
	}

	return constants.JsonWriteSuccess, nil
}
