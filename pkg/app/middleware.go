package app

import (
	"errors"
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

// login is a struct for the login form
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// GetAuthMiddleware returns the auth middleware
//
// Authenticated users will receive a cookie with the following format:
// Cookie = hash(secret, cookie domain, email, expires)|expires|email
//
// Given internal auth: Checks if username and password are correct
// Given external auth: Checks for a _forward_auth cookie and validates it
func (a *Controller) GetAuthMiddleware(sessionSecret, domain, zone string) (*jwt.GinJWTMiddleware, error) {
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
			user := &User{}
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				user.Email = fmt.Sprintf("%s@leetserve.com", userID)
			}

			// Get first matched user record
			err := a.DB.Where("email = ?", user.Email).First(&user).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				a.DB.Create(&user)
			}

			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// TODO: Implement authorizator
			claims := jwt.ExtractClaims(c)
			if v, ok := data.(*User); ok && v.Email == claims[identityKey].(string) {
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
		TokenLookup:   "header: Authorization, query: token, cookie: wikileet",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		return nil, err
	}

	return authMiddleware, nil

}
