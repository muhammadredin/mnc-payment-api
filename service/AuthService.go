package service

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	dto "PaymentAPI/dto/response"
	"PaymentAPI/utils"
	"errors"
)

type AuthService interface {
	Login(request req.CustomerRequest) (dto.AuthResponse, error)
	Logout(accessToken string, refreshToken string) error
	GetNewAccessToken(refreshToken string) (dto.AuthResponse, error)
}

type authService struct {
	customerService     CustomerService
	refreshTokenService RefreshTokenService
	blacklistService    BlacklistService
}

func NewAuthService(customerService CustomerService, refreshTokenService RefreshTokenService, blacklistService BlacklistService) AuthService {
	return &authService{customerService: customerService, refreshTokenService: refreshTokenService, blacklistService: blacklistService}
}

func (a *authService) Login(request req.CustomerRequest) (dto.AuthResponse, error) {
	customer, err := a.customerService.GetCustomerByUsernameAuth(request.Username)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if !utils.BCryptCompare(request.Password, []byte(customer.Password)) {
		return dto.AuthResponse{}, errors.New(constants.LoginUnauthorizedError)
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

func (a *authService) Logout(accessToken string, refreshToken string) error {
	err := a.blacklistService.BlacklistToken(accessToken)
	if err != nil {
		return err
	}

	err = a.refreshTokenService.DeleteRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) GetNewAccessToken(refreshToken string) (dto.AuthResponse, error) {
	newRefreshToken, err := a.refreshTokenService.RotateRefreshToken(refreshToken)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	customer, err := a.customerService.GetCustomerByIdAuth(newRefreshToken.CustomerId)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	accessToken, err := utils.GenerateAccessToken(customer)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.RefreshToken,
	}, nil
}
