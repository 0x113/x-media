package handler

import (
	"net/http"

	"github.com/0x113/x-media/user/models"
	"github.com/0x113/x-media/user/service"

	"github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo"
)

type userHandler struct {
	userService service.UserService
}

// NewUserHandler initiates user handlers
func NewUserHandler(router *echo.Echo, userService service.UserService) {
	handler := &userHandler{userService}
	// swagger
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	router.File("/swagger.yaml", "./docs/swagger.yaml")
	router.GET("/docs", echo.WrapHandler(sh))

	router.POST("/api/v1/user/create", handler.CreateUser)
	router.POST("/api/v1/user/validate", handler.ValidateUser)
}

// @Summary Create user
// @Description Creates new user in the database
// @ID create-new-user
// @Accept  json
// @Produce  json
// @Param name body userPayload true "User credentials"
// @Success 201 {object} userCreateResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /create [post]
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

// @Summary Validate user
// @Description Check if user credentials are correct
// @ID validate-user
// @Accept  json
// @Produce  json
// @Param name body userValidatePayload true "User credentials"
// @Success 200 {object} models.TokenClaims
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /validate [post]
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
