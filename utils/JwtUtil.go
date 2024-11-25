package utils

import (
	"PaymentAPI/config"
	"PaymentAPI/constants"
	"PaymentAPI/entity"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"strings"
	"time"
)

type M map[string]interface{}

func GenerateAccessToken(customer entity.Customer) (string, error) {
	claims := entity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.ApplicationName,
			Subject:   customer.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginExpirationDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(config.JwtSigningMethod, claims)

	signedToken, err := token.SignedString(config.JwtSignatureKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	// Retrieve the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New(constants.AuthorizationHeaderMissingError)
	}

	// Bearer token format: "Bearer <token>"
	// Split the header by space to extract the token part
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New(constants.AuthorizationHeaderInvalidError)
	}

	// Extract the JWT token
	return parts[1], nil
}

func ParseAndVerifyAccessToken(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if method != config.JwtSigningMethod {
			return nil, fmt.Errorf("invalid signing method: expected %v, got %v", config.JwtSigningMethod.Alg(), method.Alg())
		}

		return config.JwtSignatureKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf(constants.JwtTokenInvalidError)
	}

	if !token.Valid {
		return nil, errors.New(constants.JwtTokenInvalidError)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	return claims, nil
}

func GetExpirationFromClaimsAsString(accessToken string) (string, error) {
	claims, err := ParseAndVerifyAccessToken(accessToken)
	if err != nil {
		return "", err
	}

	exp, ok := claims["exp"]
	if !ok {
		return "", errors.New("expiration claim (exp) not found in token")
	}

	expStr := ""
	switch v := exp.(type) {
	case float64:
		expStr = strconv.FormatInt(int64(v), 10)
	case string:
		expStr = v
	default:
		return "", errors.New("unexpected type for expiration claim")
	}

	return expStr, nil
}

func GetCustomerIdFromClaims(accessToken string) (string, error) {
	claims, err := ParseAndVerifyAccessToken(accessToken)
	if err != nil {
		return "", err
	}

	id, ok := claims["sub"]
	if !ok {
		return "", errors.New("Customer Id (sub) not found in token")
	}
	return id.(string), nil
}
