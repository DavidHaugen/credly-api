package httpservice

import (
	"github.com/DavidHaugen/golang-boilerplate/internal"
	"github.com/DavidHaugen/golang-boilerplate/internal/config"
	"github.com/DavidHaugen/golang-boilerplate/internal/marvel"
)

type applicationConfig struct {
	config *config.Config
	marvel internal.MarvelService
}

func newApplicationConfig(config *config.Config) *applicationConfig {
	applicationConfig := &applicationConfig{
		config: config,
	}

	applicationConfig.addDependencies()

	return applicationConfig
}

// addDependencies :
func (a *applicationConfig) addDependencies() *applicationConfig {
	a.marvel = marvel.NewService(a.config.Marvel.PublicAPIKey, a.config.Marvel.PrivateAPIKey)
	return a
}
