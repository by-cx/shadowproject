package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"log"
)

// Config structure
// Every configuration directive has to go here
type NodeProxyConfig struct {
	MasterHost  string `env:"MASTER_HOST" envDefault:"localhost"`
	MasterPort  int    `env:"MASTER_PORT" envDefault:"8080"`
	MasterProto string `env:"MASTER_PROTO" envDefault:"http"`
	Port        int    `env:"PORT" envDefault:"8081"`
}

// Handle the config
func ProcessEnvironmentVariables(config *NodeProxyConfig) {
	err := env.Parse(config)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", config)
}
