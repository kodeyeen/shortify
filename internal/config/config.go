package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

const (
	PersistenceTypeInmemory = "inmemory"
	PersistenceTypePostgres = "postgres"
)

type Config struct {
	Env             string           `yaml:"env" env:"ENV" env-default:"dev"`
	PersistenceType string           `yaml:"persistence_type" env:"PERSISTENCE_TYPE"`
	Alias           AliasConfig      `yaml:"alias"`
	HTTPServer      HTTPServerConfig `yaml:"http_server"`
	Postgres        PostgresConfig   `yaml:"postgres"`
}

type AliasConfig struct {
	Length  int    `yaml:"length" env:"LENGTH" env-default:"10"`
	Charset string `yaml:"charset" env:"CHARSET" env-required:"true"`
}

type HTTPServerConfig struct {
	Port            int           `yaml:"port" env:"HTTP_SERVER_PORT" env-required:"true"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"HTTP_SERVER_READ_TIMEOUT" env-default:"3s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"HTTP_SERVER_WRITE_TIMEOUT" env-default:"3s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"30s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SERVER_SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Database string `yaml:"database" env:"POSTGRES_DB"`
	Username string `yaml:"username" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
}

func MustLoad() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", cfgPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	return &cfg
}
