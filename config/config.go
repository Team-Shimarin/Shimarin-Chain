package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port           string `default:"50051"`
	MinorAccountID string `default:"root"` //NOTE: Root Account is generated GENESIS_BLOCK
	RedisPort      string `default:"6379"`
	RedisHost      string `default:"localhost"`
	RedisNodeCount int    `default:4`
}

var config Config

func init() {
	_ = envconfig.Process("anzu", &config)
}

func GetConfig() *Config {
	return &config
}
