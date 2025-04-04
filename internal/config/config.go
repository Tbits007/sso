package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string     	 `yaml:"env" env-default:"local"`
	StoragePath    string     	 `yaml:"storage_path" env-required:"true"`
	MigrationsPath string 	  	 `yaml:"migrations_path" env-required:"true"`
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPC           GRPCConfig 	 `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}


func MustLoad() *Config {
    configPath := os.Getenv("CONFIG_PATH")
    if configPath == "" {
        log.Fatal("CONFIG_PATH environment variable is not set")
    }

    if _, err := os.Stat(configPath); err != nil {
        log.Fatalf("error opening config file: %s", err)
    }

    var cfg Config

    err := cleanenv.ReadConfig(configPath, &cfg)
    if err != nil {
        log.Fatalf("error reading config file: %s", err)
    }

    return &cfg
}
