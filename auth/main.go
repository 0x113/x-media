// @title Authentication service API
// @version 1.0.0
// @description The main purpose of the API is to authenticate user
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes http
// @host localhost:8003
// @BasePath /api/v1/auth/token
package main

import (
	"net/http"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/data"
	"github.com/0x113/x-media/auth/databases"
	"github.com/0x113/x-media/auth/handler"
	"github.com/0x113/x-media/auth/service"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	router *echo.Echo
}

func (srv *Server) initServer() error {
	// load config
	if err := common.LoadConfig(); err != nil {
		return err
	}

	// connect to the redis
	if err := databases.Database.Init(); err != nil {
		return err
	}

	// set up router
	srv.router = echo.New()
	srv.router.HideBanner = true
	return nil
}

func main() {
	srv := &Server{}

	if err := srv.initServer(); err != nil {
		log.Fatalf("Couldn't initialize server: %v", err)
	}

	httpClient := &http.Client{}
	authRepository := data.NewRedisAuthRepository()
	authService := service.NewAuthService(httpClient, authRepository)
	handler.NewAuthHandler(srv.router, authService)

	srv.router.Start(":" + common.Config.Port)
}
