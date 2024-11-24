package repository

import (
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"PaymentAPI/storage"
	"PaymentAPI/utils"
)

type BlacklistRepository interface {
	GetAll() ([]entity.Blacklist, error)
	CreateBlacklist(accessToken string) error
}

type blacklistRepository struct {
	JsonStorage storage.JsonFileHandler[entity.Blacklist]
}

func NewBlacklistRepository(jsonStorage storage.JsonFileHandler[entity.Blacklist]) BlacklistRepository {
	return &blacklistRepository{JsonStorage: jsonStorage}
}

func (r *blacklistRepository) CreateBlacklist(accessToken string) error {
	data, err := r.JsonStorage.ReadFile(constants.BlacklistJsonPath)
	if err != nil {
		return err
	}

	expStr, err := utils.GetExpirationFromClaimsAsString(accessToken)
	if err != nil {
		return err
	}

	blacklist := entity.Blacklist{
		AccessToken: accessToken,
		ExpiresAt:   expStr,
	}
	data = append(data, blacklist)

	_, err = r.JsonStorage.WriteFile(data, constants.BlacklistJsonPath)
	if err != nil {
		return err
	}

	return nil
}

func (r *blacklistRepository) GetAll() ([]entity.Blacklist, error) {
	data, err := r.JsonStorage.ReadFile(constants.BlacklistJsonPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
