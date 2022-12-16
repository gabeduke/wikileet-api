package config

import (
	"github.com/gabeduke/wikileet-api/pkg/store"
	"github.com/kelseyhightower/envconfig"
)

// New processes and returns a new config
func New() (*App, error) {
	var config App
	err := envconfig.Process("cfg", &config)
	return &config, err
}

// App is the configuration for the application
type App struct {
	Debug            bool
	Port             int
	DatabaseHost     string `split_words:"true" default:"localhost"`
	DatabasePort     int    `split_words:"true" default:"5432"`
	DatabaseUser     string `split_words:"true" default:"postgres"`
	DatabasePassword string `split_words:"true" default:"postgres"`
	DatabaseName     string `split_words:"true" default:"gorm"`
}

// AppI is an interface for the App struct
type AppI interface {
	ParseDialerConfig() *store.Dialer
}

// ParseDialerConfig parses the config into a Dialer struct
func (a *App) ParseDialerConfig() *store.Dialer {
	dialer := &store.Dialer{
		Host:     a.DatabaseHost,
		Port:     a.DatabasePort,
		User:     a.DatabaseUser,
		Password: a.DatabasePassword,
		Name:     a.DatabaseName,
		SslMode:  "disable",
	}
	return dialer
}
