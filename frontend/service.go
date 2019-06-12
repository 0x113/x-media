package frontend

import (
	"os"

	"github.com/0x113/x-media/env"

	log "github.com/sirupsen/logrus"
)

type FrontendService interface {
	// FrontendDir checks if frontend dir exists and returns path to it
	FrontendDir() (string, error)
}

type frontendService struct{}

func NewFrontendService() FrontendService {
	return &frontendService{}
}

func (s *frontendService) FrontendDir() (string, error) {
	frontendDirPath := env.EnvString("frontend_dir")
	if _, err := os.Stat(frontendDirPath); os.IsNotExist(err) {
		log.Errorf("Cannot find frontend directory %s: %v", frontendDirPath, err)
		return "", err
	}
	return frontendDirPath, nil
}
