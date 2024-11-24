package service

import (
	dto "PaymentAPI/dto/response"
	"PaymentAPI/utils"
)

type AuthService interface {
	Login(username string) (dto.AuthResponse, error)
	Logout(accessToken string) error
}

type authService struct {
	customerService     CustomerService
	refreshTokenService RefreshTokenService
	blacklistService    BlacklistService
}

func NewAuthService(customerService CustomerService, refreshTokenService RefreshTokenService, blacklistService BlacklistService) AuthService {
	return &authService{customerService: customerService, refreshTokenService: refreshTokenService, blacklistService: blacklistService}
}

func (a *authService) Login(username string) (dto.AuthResponse, error) {
	customer, err := a.customerService.GetCustomerByUsername(username)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	accessToken, err := utils.GenerateAccessToken(customer)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	refreshToken, err := a.refreshTokenService.GenerateRefreshToken(customer.Id)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.RefreshToken,
	}, nil
}

func (a *authService) Logout(accessToken string) error {
	err := a.blacklistService.BlacklistToken(accessToken)
	if err != nil {
		return err
	}

	return nil
}

//func (a *authService) GetNewAccessToken(refreshToken string) (dto.AuthResponse, error) {
//	token, err := a.refreshTokenService.RotateRefreshToken(refreshToken)
//	if err != nil {
//		return dto.AuthResponse{}, err
//	}
//
//	a.customerService.GetCustomerByUsername(token.CustomerId)
//
//	accessToken, err := utils.GenerateAccessToken()
//}
