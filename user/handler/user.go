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
}

// CreateUser calls service layer to create a new user in the database
func (h *userHandler) CreateUser(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		errMsg := &models.Error{
			Code:    http.StatusUnprocessableEntity,
			Message: "Provided user data is invalid",
		}
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	if err := h.userService.CreateUser(u); err != nil {
		errMsg := &models.Error{
			Code:    http.StatusInternalServerError,
			Message: "Couldn't create new user",
		}
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	msg := &models.Message{"Successfully create new user"}
	return c.JSON(http.StatusCreated, msg)
}
