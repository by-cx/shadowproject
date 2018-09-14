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
	ProxyTarget string `env:"PROXY_TARGET" envDefault:"localhost"`
	S3Endpoint  string `env:"S3_ENDPOINT" envDefault:"127.0.0.1:9000"`
	S3Bucket    string `env:"S3_BUCKET" envDefault:"shadowproject"`
	S3SSL       bool   `env:"S3_SSL" envDefault:"0"`
	S3AccessKey string `env:"S3_ACCESS_KEY"`
	S3SecretKey string `env:"S3_SECRET_KEY"`
}

// Handle the config
func ProcessEnvironmentVariables(config *NodeProxyConfig) {
	err := env.Parse(config)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", config)
}
