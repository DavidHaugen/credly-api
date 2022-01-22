package cmd

import (
	"log"

	"github.com/DavidHaugen/golang-boilerplate/internal/config"
	"github.com/DavidHaugen/golang-boilerplate/internal/httpservice"
)

func RunHTTPServer() {
	config, err := config.GetConfig()
	switch {
	case config == nil:
		log.Fatal("unable to read nil config file")
	case err != nil:
		log.Fatal("error fetching config:", err)
	}
	httpservice.ListenAndServe()
}
