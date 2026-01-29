/*
Copyright © 2024 Raznar Lab <xabhista19@raznar.id>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/app"
	"raznar.id/invoice-broker/internal/services"
	"raznar.id/invoice-broker/internal/workers"
)

func startServer() {
	conf, err := configs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Start background workers (e.g., 5 workers)
	// This will now use the GlobalLevel (Debug/Info) set in initConfig
	workers.Init(5000, 20)

	log.Info().
		Bool("debug_enabled", debugMode).
		Msg("Starting Invoice Broker server...")

	s := services.New(conf)
	if err = app.Start(conf, s); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting the service")
	}
}
