package database

import (
	"fmt"
	"github.com/gabeduke/wikileet-api/pkg/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

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

// Connect connects to a database
func (c *Dialer) Connect() (*app.Controller, error) {

	host := fmt.Sprintf("host=%s", c.Host)
	port := fmt.Sprintf("port=%d", c.Port)
	user := fmt.Sprintf("user=%s", c.User)
	password := fmt.Sprintf("password=%s", c.Password)
	dbname := fmt.Sprintf("dbname=%s", c.Name)
	sslmode := fmt.Sprintf("sslmode=%s", c.SslMode)

	dsn := strings.Join([]string{host, port, user, password, dbname, sslmode}, " ")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(&app.Item{}, &app.User{}, &app.Workspace{})
	if err != nil {
		return nil, err
	}

	return &app.Controller{database}, nil
}
