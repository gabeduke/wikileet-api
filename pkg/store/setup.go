package store

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/heptiolabs/healthcheck"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
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

// Dial connects to a database
func (c *Dialer) Dial() (*App, error) {

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

	err = database.AutoMigrate(&Item{}, &User{}, &Workspace{})
	if err != nil {
		return nil, err
	}

	return &App{database}, nil
}

func (a *App) GetHealthHandler() (healthcheck.Handler, error) {

	database, err := a.DB.DB()
	if err != nil {
		return nil, err
	}

	health := healthcheck.NewHandler()

	health.AddReadinessCheck("database_ready", healthcheck.DatabasePingCheck(database, 5*time.Second))
	health.AddLivenessCheck("database_live", healthcheck.DatabasePingCheck(database, 5*time.Second))
	return health, nil
}

func (a *App) GetUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the email addr from the header
		email := c.GetHeader(headerUser)
		c.Set(contextKeyUserEmail, email)

		// Get the workspace from the header
		wrk := c.GetHeader(headerWorkspace)
		if wrk == "" {
			wrk = "default"
		}
		c.Set(contextKeyWorkspace, wrk)

		// Get first matched user record
		u := &User{Email: email}
		err := a.DB.Where("email = ?", email).First(&u).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.DB.Create(&u)
		}
		c.Set(contextKeyUserID, u.ID)

		// Get first matched workspace record
		w := &Workspace{Name: wrk}
		a.DB.Where("name = ?", wrk).First(&w)
		err = a.DB.First(&w).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.DB.Create(&w)
		}
		c.Set(contextKeyWorkspaceID, w.ID)

		//
		uuid := uuid.New()
		c.Set("uuid", uuid)
		logrus.Tracef("The request with uuid %s is started", uuid)
		c.Next()
		logrus.Tracef("The request with uuid %s is served", uuid)
	}
}
