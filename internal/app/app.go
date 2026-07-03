package app

import (
	"github.com/DilankaHer/sop-in-go/internal/logger"
	"github.com/DilankaHer/sop-in-go/internal/middleware"
)

type App struct {
	Config     *Config
	Logger     *logger.Logger
	Middleware middleware.Middleware
}

func NewApp() (*App, error) {
	config, err := GetConfig()
	logger := logger.NewLogger()
	middleware := middleware.NewMiddleware(logger)
	if err != nil {
		return nil, err
	}
	return &App{Config: config, Logger: logger, Middleware: middleware}, nil
}
