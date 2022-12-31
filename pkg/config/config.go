package config

import (
	"github.com/gabeduke/wikileet-api/pkg/store"
	"github.com/kelseyhightower/envconfig"
)

// New processes and returns a new config
func New() (*AppConfig, error) {
	var config AppConfig
	err := envconfig.Process("cfg", &config)
	return &config, err
}

// AppConfig is the configuration for the application
type AppConfig struct {
	Debug            bool
	Port             int
	AuthEnabled      bool   `default:"false"`
	SessionSecret    string `split_words:"true" default:"secret"`
	Domain           string `default:"localhost"`
	Zone             string `default:"test"`
	DatabaseHost     string `split_words:"true" default:"localhost"`
	DatabasePort     int    `split_words:"true" default:"5432"`
	DatabaseUser     string `split_words:"true" default:"postgres"`
	DatabasePassword string `split_words:"true" default:"postgres"`
	DatabaseName     string `split_words:"true" default:"gorm"`
}

// AppConfigI is an interface for the AppConfig struct
type AppConfigI interface {
	ParseDialerConfig() *store.Dialer
}

func (a *AppConfig) GetSessionSecret() string {
	return a.SessionSecret
}

func (a *AppConfig) GetAuthEnabled() bool {
	return a.AuthEnabled
}

func (a *AppConfig) GetDomain() string {
	return a.Domain
}

func (a *AppConfig) GetZone() string {
	return a.Zone
}

// ParseDialerConfig parses the config into a Dialer struct
func (a *AppConfig) ParseDialerConfig() *store.Dialer {
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
