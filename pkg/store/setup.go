package store

import (
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/heptiolabs/healthcheck"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
)

const identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
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

func (a *App) GetAuthMiddleware(sessionSecret, domain, zone string, authEnabled bool) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:        zone,
		Key:          []byte(sessionSecret),
		Timeout:      time.Hour,
		MaxRefresh:   time.Hour,
		IdentityKey:  identityKey,
		SendCookie:   true,
		SecureCookie: false,
		CookieName:   "wikileet", // default jwt
		CookieDomain: domain,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Email: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			email := c.GetString(contextKeyUserEmail)
			if email == "" {
				return nil, errors.New("email not found in context")
			}

			if authEnabled {
				logrus.Info("Auth enabled")
				var loginVals login
				if email == "" {
					if err := c.ShouldBind(&loginVals); err != nil {
						return "", jwt.ErrMissingLoginValues
					}
					userID := loginVals.Username
					password := loginVals.Password

					if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
						return &User{
							Email: "admin@leetserve.com",
						}, nil
					}
					return nil, jwt.ErrFailedAuthentication
				}
			}
			logrus.Info("Auth disabled")

			// Get first matched user record
			u := &User{Email: email}
			dberr := a.DB.Where("email = ?", email).First(&u).Error
			if errors.Is(dberr, gorm.ErrRecordNotFound) {
				a.DB.Create(&u)
			}

			return u, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.Email == c.GetString(contextKeyUserEmail) {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: wikileet",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		return nil, err
	}

	return authMiddleware, nil

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
