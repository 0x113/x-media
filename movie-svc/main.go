package main

import (
	"log"
	"net/http"

	"github.com/0x113/x-media/movie-svc/common"
	"github.com/0x113/x-media/movie-svc/data"
	"github.com/0x113/x-media/movie-svc/databases"
	"github.com/0x113/x-media/movie-svc/handler"
	"github.com/0x113/x-media/movie-svc/service"

	"github.com/labstack/echo"
)

type Server struct {
	router *echo.Echo
}

func (srv *Server) initServer() error {
	// load config from file
	if err := common.LoadConfig(); err != nil {
		return err
	}

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
		log.Fatalf("Unable to initialize server: %v", err)
	}

	movieRepository := data.NewMongoMovieRepository()
	movieService := service.NewMovieService(movieRepository, &http.Client{})
	handler.NewMovieHandler(srv.router, movieService)

	srv.router.Start(":" + common.Config.Port)
}
