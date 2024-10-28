package app

import (
	"log"

	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/app/rest"
	"raznar.id/invoice-broker/internal/pkg/database"
)

func Start(configFile string) {
	conf, err := configs.New(configFile)
	if err != nil {
		log.Fatalf(err.Error())
	}


	db := database.New(conf.DataPath)
	db.Load()

	// fiber has built-in block, so we dont need any signal block
	if err = rest.Start(conf, db); err != nil {
		log.Fatalf("An error occured when starting the bot: %s", err.Error())
	}
	
}