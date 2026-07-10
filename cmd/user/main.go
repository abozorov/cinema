package main

import (
	"log"
	"os"
	"path"

	"github.com/abozorov/cinema/cmd/user/config"
	"github.com/abozorov/cinema/cmd/user/internal/app"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error in path determination: ", err)
		return
	}
	cfg, err := config.NewConfig(path.Join(dir, "cmd", "user", "config", ".env"))

	if err != nil {
		log.Println("Error when load config file: ", err)
		return
	}

	app.Run(cfg)
}
