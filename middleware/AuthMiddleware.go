package middleware

import (
	res "PaymentAPI/dto/response"
	"PaymentAPI/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := utils.ExtractTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseAndVerifyAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("userClaims", claims)
		c.Next()
	}
}
