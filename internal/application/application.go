package application

import "github.com/Irurnnen/calc-worker/internal/config"

type Application struct {
	Config config.Config
}

func (a *Application) Run() {
	// TODO: add waitgroup
	// TODO: Add function to add gourutines
}

func New() *Application {
	return &Application{
		Config: *config.NewConfigFromEnv(),
	}
}
