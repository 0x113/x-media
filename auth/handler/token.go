package handler

import (
	"net/http"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/models"
	"github.com/0x113/x-media/auth/service"

	"github.com/labstack/echo"
)

type authHandler struct {
	authService service.AuthService
}

// NewAuthHandler initiates authentication handlers
func NewAuthHandler(router *echo.Echo, authService service.AuthService) {
	handler := &authHandler{authService}
	router.POST("/api/v1/auth/token/generate", handler.GenerateToken)
	router.POST("/api/v1/auth/token/validate", handler.GetTokenMetadata)
	router.POST("/api/v1/auth/token/refresh", handler.RefreshToken)
}

// GenerateToken calls the service layer and generates new JSON Web Token
func (h *authHandler) GenerateToken(c echo.Context) error {
	errMsg := new(models.Error)
	creds := new(models.Credentials)
	if err := c.Bind(creds); err != nil {
		errMsg.Code = http.StatusBadRequest
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg.Message)
		return err
	}

	token, err := h.authService.Login(creds)
	if err != nil {
		errMsg.Code = http.StatusInternalServerError
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	return c.JSON(http.StatusOK, token)
}

// RefreshToken calls the service layer to generate new access and refresh token
func (h *authHandler) RefreshToken(c echo.Context) error {
	errMsg := new(models.Error)
	ts := new(models.TokenString)
	if err := c.Bind(ts); err != nil {
		errMsg.Code = http.StatusBadRequest
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	token, err := h.authService.Refresh(ts.Token)
	if err != nil {
		errMsg.Code = http.StatusInternalServerError
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	return c.JSON(http.StatusOK, token)
}

// GetTokenMetadata calls the service layer to validate provided token and
// to get metadata if token is valid
func (h *authHandler) GetTokenMetadata(c echo.Context) error {
	errMsg := new(models.Error)
	ts := new(models.TokenString)
	if err := c.Bind(ts); err != nil {
		errMsg.Code = http.StatusBadRequest
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg.Message)
		return err
	}

	accessDetails, err := h.authService.ExtractTokenMetadata(ts.Token, common.Config.AccessSecret)
	if err != nil {
		errMsg.Code = http.StatusInternalServerError
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	return c.JSON(http.StatusOK, accessDetails)
}
