package repository

import (
	"PaymentAPI/constants"
	dto "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"errors"
	"github.com/google/uuid"
	"time"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(customer dto.CustomerResponse) (string, error)
	GetRefreshToken(refreshToken string) (entity.RefreshToken, error)
	DeleteRefreshToken(refreshToken string) error
}

type refreshTokenRepository struct {
	JsonStorage storage.JsonFileHandler[entity.RefreshToken]
}

func NewRefreshTokenRepository(jsonStorage storage.JsonFileHandler[entity.RefreshToken]) RefreshTokenRepository {
	return &refreshTokenRepository{JsonStorage: jsonStorage}
}

func (r *refreshTokenRepository) CreateRefreshToken(customer dto.CustomerResponse) (string, error) {
	refreshToken := entity.RefreshToken{
		RefreshToken: uuid.New().String(),
		CustomerId:   customer.Id,
		ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour).String(),
	}

	data, err := r.JsonStorage.ReadFile(constants.RefreshTokenJsonPath)
	if err != nil {
		return "", err
	}

	data = append(data, refreshToken)

	_, err = r.JsonStorage.WriteFile(data, constants.RefreshTokenJsonPath)
	if err != nil {
		return "", err
	}

	return refreshToken.RefreshToken, nil
}

func (r *refreshTokenRepository) GetRefreshToken(refreshToken string) (entity.RefreshToken, error) {
	data, err := r.JsonStorage.ReadFile(constants.RefreshTokenJsonPath)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	for _, token := range data {
		if token.RefreshToken == refreshToken {
			return token, nil
		}
	}

	return entity.RefreshToken{}, errors.New(constants.RefreshTokenNotFoundError)
}

func (r *refreshTokenRepository) DeleteRefreshToken(refreshToken string) error {
	data, err := r.JsonStorage.ReadFile(constants.RefreshTokenJsonPath)
	if err != nil {
		return err
	}

	indexToDelete := -1
	for i, token := range data {
		if token.RefreshToken == refreshToken {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		return errors.New(constants.RefreshTokenNotFoundError)
	}

	data = append(data[:indexToDelete], data[indexToDelete+1:]...)
	_, err = r.JsonStorage.WriteFile(data, constants.RefreshTokenJsonPath)
	if err != nil {
		return err
	}
	return nil
}
