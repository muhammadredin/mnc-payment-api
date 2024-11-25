package handler

import (
	"PaymentAPI/constants"
	req "PaymentAPI/dto/request"
	res "PaymentAPI/dto/response"
	"PaymentAPI/service"
	"PaymentAPI/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AuthHandler interface {
	HandleRegister(c *gin.Context)
	HandleLogin(c *gin.Context)
	HandleLogout(c *gin.Context)
	HandleRefreshToken(c *gin.Context)
}

type authHandler struct {
	authService     service.AuthService
	customerService service.CustomerService
}

func NewAuthHandler(authService service.AuthService, customerService service.CustomerService) AuthHandler {
	return &authHandler{authService, customerService}
}

func (a *authHandler) HandleRegister(c *gin.Context) {
	var request req.CustomerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: constants.InvalidRequestBodyError,
		})
		return
	}

	if _, err := a.customerService.CreateNewCustomer(request); err != nil {
		switch err.Error() {
		case constants.UsernameDuplicateError:
			c.JSON(http.StatusBadRequest, res.ErrorResponse{StatusCode: http.StatusBadRequest, ErrorMessage: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{StatusCode: http.StatusInternalServerError, ErrorMessage: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res.CommonResponse{
		StatusCode: http.StatusCreated,
		Message:    constants.CustomerCreateSuccess,
		Data:       []interface{}{},
	})
}

func (a *authHandler) HandleLogin(c *gin.Context) {
	var request req.CustomerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: constants.InvalidRequestBodyError,
		})
		return
	}

	login, err := a.authService.Login(request)

	if err != nil {
		switch err.Error() {
		case constants.LoginUnauthorizedError:
			c.JSON(http.StatusUnauthorized, res.ErrorResponse{StatusCode: http.StatusUnauthorized, ErrorMessage: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{StatusCode: http.StatusInternalServerError, ErrorMessage: err.Error()})
		return
	}

	SetCookie(c, "refresh_token", login.RefreshToken, 24)

	c.JSON(http.StatusOK, res.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    constants.LoginSuccess,
		Data:       login,
	})
	return
}

func (a *authHandler) HandleLogout(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: err.Error(),
		})
		return
	}

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: err.Error(),
		})
		return
	}

	err = a.authService.Logout(token, cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    constants.LogoutSuccess,
		Data:       []interface{}{},
	})
	return
}

func (a *authHandler) HandleRefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: err.Error(),
		})
		return
	}

	login, err := a.authService.GetNewAccessToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.ErrorResponse{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: err.Error(),
		})
		return
	}

	SetCookie(c, "refresh_token", login.RefreshToken, 24)

	c.JSON(http.StatusOK, res.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    constants.LoginSuccess,
		Data:       login,
	})
	return
}

func SetCookie(c *gin.Context, name string, value string, duration int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Duration(duration) * time.Hour),
	}

	http.SetCookie(c.Writer, cookie)
}
