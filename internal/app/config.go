package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server `yaml:"server" validate:"required"`
}

type Server struct {
	Port string `yaml:"port" validate:"required"`
}

func GetConfig() (*Config, error) {
	configName := os.Getenv("ENV")
	if os.Getenv("ENV") == "prod" {
		configName = ".env.prod"
	} else if os.Getenv("ENV") == "sit" {
		configName = ".env.sit"
	} else {
		configName = ".env.local"
	}
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("%s", formatValidationError(err.(validator.ValidationErrors)))
	}

	return &cfg, nil
}

func formatValidationError(err validator.ValidationErrors) string {
	errors := []string{}
	for _, e := range err {
		errors = append(errors, e.Field())
	}
	return fmt.Sprintf("missing/invalid vars: %s", strings.Join(errors, ", "))
}
