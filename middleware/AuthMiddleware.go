package middleware

import (
	"PaymentAPI/constants"
	res "PaymentAPI/dto/response"
	"PaymentAPI/service"
	"PaymentAPI/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(blacklistService service.BlacklistService) gin.HandlerFunc {
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

		_, err = utils.ParseAndVerifyAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		blacklisted, err := blacklistService.IsBlacklisted(accessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, res.ErrorResponse{
				StatusCode:   http.StatusInternalServerError,
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		if blacklisted {
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: constants.JwtTokenInvalidError,
			})
			c.Abort()
			return
		}

		id, err := utils.GetCustomerIdFromClaims(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("authenticatedUser", id)
		c.Next()
	}
}
