package store

import (
	"fmt"
	"github.com/heptiolabs/healthcheck"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
)

type App struct {
	db *gorm.DB
}

// DialerI is an interface for the Dialer struct
type DialerI interface {
	Dial() (*gorm.DB, error)
}

// Dialer is a struct for connecting to a database
type Dialer struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SslMode  string
}

// Dial connects to a database
func (c *Dialer) Dial() (*App, error) {

	host := fmt.Sprintf("host=%s", c.Host)
	port := fmt.Sprintf("port=%d", c.Port)
	user := fmt.Sprintf("user=%s", c.User)
	password := fmt.Sprintf("password=%s", c.Password)
	dbname := fmt.Sprintf("dbname=%s", c.Name)
	sslmode := fmt.Sprintf("sslmode=%s", c.SslMode)

	dsn := strings.Join([]string{host, port, user, password, dbname, sslmode}, " ")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(&Item{})
	if err != nil {
		return nil, err
	}

	return &App{database}, nil
}

func (a *App) GetHealthHandler() (healthcheck.Handler, error) {

	database, err := a.db.DB()
	if err != nil {
		return nil, err
	}

	health := healthcheck.NewHandler()

	health.AddReadinessCheck("database_ready", healthcheck.DatabasePingCheck(database, 5*time.Second))
	health.AddLivenessCheck("database_live", healthcheck.DatabasePingCheck(database, 5*time.Second))
	return health, nil
}
