package app

type App struct {
	Config *Config
}

func NewApp() (*App, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	return &App{Config: config}, nil
}
