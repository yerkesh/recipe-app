package config

import (
	"time"

	"github.com/golang-migrate/migrate/v4/database/postgres"
)

var (
	Configuration *Config
	ServerProfile string
)

const (
	Dev   = "dev"
	Stage = "stage"
	Prod  = "prod"
)

type Migration struct {
	MigrationsTable       string `yaml:"migrationsTable"`
	MigrationsTableQuoted bool
	MultiStatementEnabled bool `yaml:"multiStatementEnabled"`
	DatabaseName          string `yaml:"databaseName"`
	SchemaName            string
	StatementTimeout      time.Duration `yaml:"statementTimeout"`
	MultiStatementMaxSize int
}

// Config properties. All configurations should be described here.
type Config struct {
	ServiceName string `yaml:"serviceName"`

	Version string `yaml:"version"`

	Server struct {
		Port    string `yaml:"port"`
		Profile string `yaml:"profile"`
	}

	Database struct {
		URI  string `yaml:"uri"`
		Name string `yaml:"name"`
	}

	Migration postgres.Config
}
