package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"errors"
	"github.com/google/uuid"
	"time"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(customerId string) (entity.RefreshToken, error)
	GetRefreshToken(refreshToken string) (entity.RefreshToken, error)
	GetAllRefreshToken() ([]entity.RefreshToken, error)
	DeleteRefreshToken(refreshToken string) error
}

type refreshTokenRepository struct {
	JsonStorage storage.JsonFileHandler[entity.RefreshToken]
}

func NewRefreshTokenRepository(jsonStorage storage.JsonFileHandler[entity.RefreshToken]) RefreshTokenRepository {
	return &refreshTokenRepository{JsonStorage: jsonStorage}
}

func (r *refreshTokenRepository) CreateRefreshToken(customerId string) (entity.RefreshToken, error) {
	refreshToken := entity.RefreshToken{
		RefreshToken: uuid.New().String(),
		CustomerId:   customerId,
		ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour).Format(time.RFC3339),
	}

	data, err := r.JsonStorage.ReadFile(constants.RefreshTokenJsonPath)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	data = append(data, refreshToken)

	_, err = r.JsonStorage.WriteFile(data, constants.RefreshTokenJsonPath)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return refreshToken, nil
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

func (r *refreshTokenRepository) GetAllRefreshToken() ([]entity.RefreshToken, error) {
	data, err := r.JsonStorage.ReadFile(constants.RefreshTokenJsonPath)
	if err != nil {
		return nil, err
	}

	return data, nil
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
