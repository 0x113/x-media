package handler

import (
	"net/http"

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
}

// GenerateToken calls the service layer and generates new JSON Web Token
func (h *authHandler) GenerateToken(c echo.Context) error {
	claims := new(models.AccessDetails)
	if err := c.Bind(claims); err != nil {
		errMsg := &models.Error{
			Code:    http.StatusUnprocessableEntity,
			Message: "Provided data is invalid",
		}
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	tokenDetails, err := h.authService.GenerateJWT(claims)
	if err != nil {
		errMsg := &models.Error{
			Code:    http.StatusInternalServerError,
			Message: "Couldn't generate token for user",
		}
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	return c.JSON(http.StatusOK, tokenDetails)
}
