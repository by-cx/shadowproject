package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"log"
)

// Config structure
// Every configuration directive has to go here
type MasterConfig struct {
	DatabaseFile string `env:"DBPath" envDefault:"tasks.db"`
	Port         int    `env:"PORT" envDefault:"8080"`
}

// Handle the config
func ProcessEnvironmentVariables(config *MasterConfig) {
	err := env.Parse(&config)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", config)
}
