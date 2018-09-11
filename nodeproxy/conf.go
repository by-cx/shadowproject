package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"log"
)

// Config structure
// Every configuration directive has to go here
type NodeProxyConfig struct {
	MasterHost string `env:"Host" envDefault:"localhost"`
	MasterPort int    `env:"Port" envDefault:"8080"`
}

// Handle the config
func ProcessEnvironmentVariables(config *NodeProxyConfig) {
	err := env.Parse(&config)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", config)
}
