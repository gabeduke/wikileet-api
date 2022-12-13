package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

var DB *gorm.DB

type DatabaseConfig struct {
	Host     string
	Port     int    `default:"5432"`
	User     string `default:"postgres"`
	Password string
	DBName   string `default:"gorm"`
}

func ConnectDatabase(cfg *DatabaseConfig) {

	host := fmt.Sprintf("host=%s", cfg.Host)
	port := fmt.Sprintf("port=%d", cfg.Port)
	user := fmt.Sprintf("user=%s", cfg.User)
	password := fmt.Sprintf("password=%s", cfg.Password)
	dbname := fmt.Sprintf("dbname=%s", cfg.DBName)
	sslmode := "sslmode=disable"

	dsn := strings.Join([]string{host, port, user, password, dbname, sslmode}, " ")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&Item{})
	if err != nil {
		return
	}

	DB = database
}
