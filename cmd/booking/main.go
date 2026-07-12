package main

import (
	"log"
	"os"
	"path"

	"github.com/abozorov/cinema/cmd/booking/internal/app"
	"github.com/abozorov/cinema/cmd/booking/internal/config"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error in path determination: ", err)
		return
	}
	conf, err := config.NewConfig(path.Join(dir, "cmd", "booking", "internal", "config", "booking_config.env"))

	if err != nil {
		log.Println("Error when load config file: ", err)
		return
	}

	app.Run(conf)
}
