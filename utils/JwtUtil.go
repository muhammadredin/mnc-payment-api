package utils

import (
	dto "PaymentAPI/dto/response"
	"PaymentAPI/entity"
	"PaymentAPI/enums"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type M map[string]interface{}

var ApplicationName = "Payment API App"
var LoginExpirationDuration = time.Duration(5) * time.Minute
var JwtSigningMethod = jwt.SigningMethodHS256
var JwtSignatureKey = []byte("the secret of kalimdor")

func GenerateAccessToken(customer dto.CustomerResponse) (string, error) {
	claims := entity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ApplicationName,
			Subject:   customer.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(LoginExpirationDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Username: customer.Username,
		Role:     enums.ROLE_USER,
	}

	token := jwt.NewWithClaims(JwtSigningMethod, claims)

	signedToken, err := token.SignedString(JwtSignatureKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
