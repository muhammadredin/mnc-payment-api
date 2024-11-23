package entity

import (
	"PaymentAPI/enums"
	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	Username string     `json:"username"`
	Role     enums.Role `json:"role"`
}
