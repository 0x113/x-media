package handler

import (
	"net/http"

	"github.com/0x113/x-media/user/models"
	"github.com/0x113/x-media/user/service"

	"github.com/labstack/echo"
)

type userHandler struct {
	userService service.UserService
}

// NewUserHandler initiates user handlers
func NewUserHandler(router *echo.Echo, userService service.UserService) {
	handler := &userHandler{userService}
	router.POST("/api/v1/user/create", handler.CreateUser)
	router.POST("/api/v1/user/validate", handler.ValidateUser)
}

// CreateUser calls service layer to create a new user in the database
func (h *userHandler) CreateUser(c echo.Context) error {
	errMsg := new(models.Error)
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		errMsg.Code = http.StatusBadRequest
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	if err := h.userService.CreateUser(u); err != nil {
		errMsg.Code = http.StatusInternalServerError
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	msg := &models.Message{"Successfully create new user"}
	return c.JSON(http.StatusCreated, msg)
}

// ValidateUser calls the service to check if provided credentials matches with
// the user in the database
func (h *userHandler) ValidateUser(c echo.Context) error {
	errMsg := new(models.Error)
	creds := new(models.Credentials)
	if err := c.Bind(creds); err != nil {
		errMsg.Code = http.StatusBadRequest
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	claims, err := h.userService.ValidateUser(creds)
	if err != nil {
		errMsg.Code = http.StatusInternalServerError
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	return c.JSON(http.StatusOK, claims)

}
