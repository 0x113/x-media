// @title Tv show service API
// @version 1.0.0
// @description Tv shows API allows to get data from the third party API (TVmaze at this moment) about the tv show from the local drive.
// @description The main purpose of the API is to update data, save it to the database and return it in the JSON format.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes http
// @host localhost:8001
// @BasePath /api/v1/tvshows
package main

import (
	"net/http"

	"github.com/0x113/x-media/tvshow/common"
	"github.com/0x113/x-media/tvshow/data"
	"github.com/0x113/x-media/tvshow/databases"
	"github.com/0x113/x-media/tvshow/handler"
	"github.com/0x113/x-media/tvshow/service"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	router *echo.Echo
}

func (srv *Server) initServer() error {
	// load config from file
	if err := common.LoadConfig(); err != nil {
		log.Errorf("Unable to load config file, err: %v", err)
		return err
	}

	// initialize database
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

	client := &http.Client{}
	tvShowRepository := data.NewMongoTVShowRepository()
	tvShowService := service.NewTVShowService(client, tvShowRepository)
	handler.NewTVShowHandler(srv.router, tvShowService)

	srv.router.Start(":" + common.Config.Port)
}
