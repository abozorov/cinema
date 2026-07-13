package main

import (
	"log"
	"os"
	"path"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/config"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/app"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error in path determination: ", err)
		return
	}
	conf, err := config.NewConfig(path.Join(dir, "cmd", "api_gateway", "internal", "config", "api_gateway_config.env"))

	if err != nil {
		log.Println("Error when load config file: ", err)
		return
	}

	app.Run(conf)
}
