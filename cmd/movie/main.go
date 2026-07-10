package main

import (
	"log"
	"os"
	"path"

	"github.com/abozorov/cinema/cmd/movie/internal/app"
	"github.com/abozorov/cinema/cmd/movie/internal/config"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error in path determination: ", err)
		return
	}
	cfg, err := config.NewConfig(path.Join(dir, "cmd", "movie", "config", "movie_config.env"))

	if err != nil {
		log.Println("Error when load config file: ", err)
		return
	}

	app.Run(cfg)
}
