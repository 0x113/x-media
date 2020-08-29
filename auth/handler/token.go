package handler

import (
	"net/http"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/models"
	"github.com/0x113/x-media/auth/service"

	"github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo"
)

type authHandler struct {
	authService service.AuthService
}

// NewAuthHandler initiates authentication handlers
func NewAuthHandler(router *echo.Echo, authService service.AuthService) {
	handler := &authHandler{authService}
	// swagger
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	router.File("/swagger.yaml", "./docs/swagger.yaml")
	router.GET("/docs", echo.WrapHandler(sh))

	router.POST("/api/v1/auth/token/generate", handler.GenerateToken)
	router.POST("/api/v1/auth/token/validate", handler.GetTokenMetadata)
	router.POST("/api/v1/auth/token/refresh", handler.RefreshToken)
}

// @Summary Generate token
// @Description Generates new access and refresh token for the user
// @ID generate-new-token
// @Accept  json
// @Produce  json
// @Param name body generateTokenPayload true "User credentials"
// @Success 200 {object} models.TokenDetails
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /generate [post]
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

// @Summary Refresh token
// @Description Generates new access and refresh token for the user
// @ID refresh-token
// @Accept  json
// @Produce  json
// @Param name body models.TokenString true "Refresh token"
// @Success 200 {object} models.TokenDetails
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /refresh [post]
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
