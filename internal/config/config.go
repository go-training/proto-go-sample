package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type (
	// Config provides the system configuration.
	Config struct {
		Logging Logging
		Server  Server
		Otel    Otel
	}

	// Logging provides the logging configuration.
	Logging struct {
		Debug   bool   `envconfig:"APP_LOGS_DEBUG"`
		Level   string `envconfig:"APP_LOGS_LEVEL" default:"info"`
		NoColor bool   `envconfig:"APP_LOGS_COLOR"`
		Pretty  bool   `envconfig:"APP_LOGS_PRETTY"`
		Text    bool   `envconfig:"APP_LOGS_TEXT"`
	}

	// Server provides the server configuration.
	Server struct {
		Addr  string `envconfig:"-"`
		Host  string `envconfig:"APP_SERVER_HOST" default:"localhost:8080"`
		Port  string `envconfig:"APP_SERVER_PORT" default:":8080"`
		Proto string `envconfig:"APP_SERVER_PROTO" default:"http"`
		Pprof bool   `envconfig:"APP_PPROF_ENABLED"`
		Acme  bool   `envconfig:"APP_TLS_AUTOCERT"`
		Email string `envconfig:"APP_TLS_EMAIL"`
		Cert  string `envconfig:"APP_TLS_CERT"`
		Key   string `envconfig:"APP_TLS_KEY"`
		Debug bool   `envconfig:"APP_SERVER_DEBUG"`
	}

	Otel struct {
		ServiceName  string `envconfig:"APP_SERVICE_NAME" default:"grpcService"`
		ServiceType  string `envconfig:"APP_SERVICE_TYPE" default:"signoz"`
		CollectorURL string `envconfig:"APP_SERVICE_NAME" default:"localhost:4317"`
		TargetURL    string `envconfig:"APP_SERVICE_NAME" default:"localhost:8081"`
	}
)

// Environ returns the settings from the environment.
func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	defaultAddress(&cfg)
	return cfg, err
}

func cleanHostname(hostname string) string {
	hostname = strings.ToLower(hostname)
	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.TrimPrefix(hostname, "https://")

	return hostname
}

func defaultAddress(c *Config) {
	if c.Server.Key != "" || c.Server.Cert != "" || c.Server.Acme {
		c.Server.Port = ":443"
		c.Server.Proto = "https"
	}
	c.Server.Host = cleanHostname(c.Server.Host)
	c.Server.Addr = c.Server.Proto + "://" + c.Server.Host
}
