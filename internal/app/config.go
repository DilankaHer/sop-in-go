package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func GetConfig() (*Config, error) {
	env := os.Getenv("ENV")

	switch env {
	case "sit":
		err := godotenv.Load(".env.sit")
		if err != nil {
			return nil, err
		}
	case "prod":
		err := godotenv.Load(".env.prod")
		if err != nil {
			return nil, err
		}
	default:
		err := godotenv.Load(".env.local")
		if err != nil {
			return nil, err
		}
	}

	err := validateVars()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port: os.Getenv("PORT"),
	}, nil
}

func validateVars() error {
	unset := []string{}
	empty := []string{}
	invalid := []map[string]string{}
	vars := []string{"PORT"}
	var errors string

	for _, v := range vars {
		val, ok := os.LookupEnv(v)
		if !ok {
			unset = append(unset, v)
		} else if val == "" {
			empty = append(empty, v)
		}
		if v == "PORT" {
			valStr, err := strconv.Atoi(val)
			if err != nil {
				invalid = append(invalid, map[string]string{v: fmt.Sprintf("%s is not a valid port", val)})
			} else if valStr < 3000 || valStr > 65535 {
				invalid = append(invalid, map[string]string{v: fmt.Sprintf("%s is not a valid port", val)})
			}
		}
	}

	if len(unset) > 0 {
		errors = "unset variables => " + strings.Join(unset, ", ") + "\n"
	}

	if len(empty) > 0 {
		errors = errors + "empty variables => " + strings.Join(empty, ", ") + "\n"
	}

	if len(invalid) > 0 {
		var e string
		for _, v := range invalid {
			e = "invalid variables => "
			for k, err := range v {
				e = e + fmt.Sprintf("%s: %s", k, err)
			}
		}
		if e != "" {
			errors = errors + e + "\n"
		}
	}

	if errors != "" {
		return fmt.Errorf("%s", errors)
	}

	return nil
}
