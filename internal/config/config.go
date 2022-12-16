// Package config describe and import app config data.
package config

type (
	// Config -.
	Config struct {
		App      AppConfig
		HTTP     HTTPConfig
		Log      LogConfig
		Postgres PostgresConfig
	}

	// AppConfig -.
	AppConfig struct {
		AppName    string
		AppVersion string
	}

	// HTTPConfig -.
	HTTPConfig struct {
		Port string
	}

	// LogConfig -.
	LogConfig struct {
		LogLevel string
	}

	// Postgresconfig -.
	PostgresConfig struct {
		URI string
	}
)
