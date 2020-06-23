package main

import (
	"github.com/0x113/x-media/user/common"
	"github.com/0x113/x-media/user/data"
	"github.com/0x113/x-media/user/databases"
	"github.com/0x113/x-media/user/handler"
	"github.com/0x113/x-media/user/service"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	router *echo.Echo
}

func (srv *Server) initServer() error {
	// load config
	if err := common.LoadConfig(); err != nil {
		log.Errorf("Unable to load config file: %v", err)
		return err
	}

	// initialize the database
	if err := databases.Database.Init(); err != nil {
		return err
	}
	log.Infof("Successfully connected to the MySQL database")

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

	userRepository := data.NewMySQLUserRepository()
	userService := service.NewUserService(userRepository)
	handler.NewUserHandler(srv.router, userService)

	// run user service
	srv.router.Start(common.Config.Port)

	// close db connection
	defer databases.Database.DB.Close()

}
