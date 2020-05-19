package goinfluxdb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

// Config is the configuration to connect to an influxdb
type Config struct {
	Database string `envconfig:"database" required:"true"`
	Username string `envconfig:"username"`
	Password string `envconfig:"password"`
	Host     string `envconfig:"host" default:"127.0.0.1"`
	Port     int    `envconfig:"port" default:"8086"`
	Scheme   string `envconfig:"scheme" default:"http"`
}

const defaultPrefix = "INFLUXDB"

// NewInfluxDBConfigFromEnv returns a new config based on env variables
// prefix is without "_". For example INFLUXDB will match
// INFLUXDB_DATABASE, INFLUXDB_USERNAME, INFLUXDB_PASSWORD, INFLUXDB_HOST, INFLUXDB_PORT
func NewInfluxDBConfigFromEnv(prefix string) (*Config, error) {
	var fromEnv Config
	err := envconfig.Process(prefix, &fromEnv)

	if err != nil {
		return nil, err
	}

	return &fromEnv, nil
}

// NewDefaultInfluxDBConfig is similar to NewInfluxDBConfigFromEnv
// with `prefix` set to INFLUXDB
func NewDefaultInfluxDBConfig() (*Config, error) {
	return NewInfluxDBConfigFromEnv(defaultPrefix)
}

// NewInfluxClient returns a new influxdb client based on the provided config
func NewInfluxClient(config *Config) influxdb2.Client {
	url := fmt.Sprintf("%s://%s:%d", config.Scheme, config.Host, config.Port)
	auth := ""
	if config.Username != "" && config.Password != "" {
		auth = fmt.Sprintf("%s:%s", config.Username, config.Password)
	}
	return influxdb2.NewClient(url, auth)
}

// NewInfluxClientFromEnvConfig returns a new influxdb client using the env variables
func NewInfluxClientFromEnvConfig(prefix string) (influxdb2.Client, *Config) {
	config, err := NewInfluxDBConfigFromEnv(prefix)
	if err != nil {
		panic(err)
	}
	return NewInfluxClient(config), config
}

// NewDefaultInfluxClient is similar to NewInfluxClientFromEnvConfig but always uses INFLUXDB as the prefix
func NewDefaultInfluxClient() (influxdb2.Client, *Config) {
	return NewInfluxClientFromEnvConfig(defaultPrefix)
}
