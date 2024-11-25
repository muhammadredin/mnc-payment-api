package middleware

import (
	"PaymentAPI/constants"
	res "PaymentAPI/dto/response"
	"PaymentAPI/service"
	"PaymentAPI/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus" // Importing logrus for structured logging
	"net/http"
)

func AuthMiddleware(blacklistService service.BlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.WithFields(logrus.Fields{
			"clientIP": c.ClientIP(),
		})

		// Extract token from the request header
		accessToken, err := utils.ExtractTokenFromHeader(c)
		if err != nil {
			logger.Warn("Failed to extract token from request header", "error", err)
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		// Parse and verify the token
		_, err = utils.ParseAndVerifyAccessToken(accessToken)
		if err != nil {
			logger.Warn("Invalid access token", "token", accessToken, "error", err)
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		// Check if the token is blacklisted
		blacklisted, err := blacklistService.IsBlacklisted(accessToken)
		if err != nil {
			logger.Error("Error checking token blacklist", "token", accessToken, "error", err)
			c.JSON(http.StatusInternalServerError, res.ErrorResponse{
				StatusCode:   http.StatusInternalServerError,
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		if blacklisted {
			logger.Warn("Blacklisted token used", "token", accessToken)
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: constants.JwtTokenInvalidError,
			})
			c.Abort()
			return
		}

		// Extract customer ID from token claims
		id, err := utils.GetCustomerIdFromClaims(accessToken)
		if err != nil {
			logger.Warn("Failed to extract customer ID from token", "token", accessToken, "error", err)
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: err.Error(),
			})
			c.Abort()
			return
		}

		// Successfully authenticated, set user ID in the context
		logger.Info("Authentication successful", "customerId", id)
		c.Set("authenticatedUser", id)
		c.Next()
	}
}
