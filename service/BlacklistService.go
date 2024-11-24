package service

import (
	"PaymentAPI/repository"
)

type BlacklistService interface {
	IsBlacklisted(accessToken string) (bool, error)
	BlacklistToken(accessToken string) error
}

type blacklistService struct {
	blacklistRepository repository.BlacklistRepository
}

func NewBlacklistService(blacklistRepository repository.BlacklistRepository) BlacklistService {
	return &blacklistService{blacklistRepository: blacklistRepository}
}

func (b *blacklistService) BlacklistToken(accessToken string) error {
	isBlacklisted, err := b.IsBlacklisted(accessToken)
	if err != nil {
		return err
	}

	if isBlacklisted {
		return nil
	}

	err = b.blacklistRepository.CreateBlacklist(accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (b *blacklistService) IsBlacklisted(accessToken string) (bool, error) {
	blacklists, err := b.blacklistRepository.GetAll()
	if err != nil {
		return false, err
	}

	for _, blacklist := range blacklists {
		if blacklist.AccessToken == accessToken {
			return true, nil
		}
	}
	return false, nil
}
