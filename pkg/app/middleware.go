package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
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
// Given external auth: Checks for a cookie and validates it
func (a *Controller) GetAuthMiddleware(sessionSecret, domain, zone string, authInternal bool) (*jwt.GinJWTMiddleware, error) {
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
			if authInternal {
				logrus.Info("Auth internal")
				var loginVals login
				if err := c.ShouldBind(&loginVals); err != nil {
					return "", jwt.ErrMissingLoginValues
				}
				userID := loginVals.Username
				password := loginVals.Password

				if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
					user.Email = fmt.Sprintf("%s@leetserve.com", userID)
				}
			} else {
				logrus.Info("Auth external")

				cookie, err := c.Request.Cookie("_forward_auth")
				if err != nil {
					return nil, err
				}

				user.Email, err = validateCookie(c.Request, cookie, domain, sessionSecret)
				if err != nil {
					if err.Error() == "Cookie has expired" {
						logrus.Info("Cookie has expired")
						return nil, jwt.ErrExpiredToken
					} else {
						logrus.WithField("error", err).Warn("Invalid cookie")
						return nil, jwt.ErrFailedAuthentication
					}
				}
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

// validateCookie verifies that a cookie matches the expected format of:
// Cookie = hash(secret, cookie domain, email, expires)|expires|email
func validateCookie(r *http.Request, c *http.Cookie, domain, sessionSecret string) (string, error) {
	parts := strings.Split(c.Value, "|")

	if len(parts) != 3 {
		return "", errors.New("invalid cookie format")
	}

	mac, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", errors.New("unable to decode cookie mac")
	}

	expectedSignature := cookieSignature(parts[2], parts[1], domain, sessionSecret)
	expected, err := base64.URLEncoding.DecodeString(expectedSignature)
	if err != nil {
		return "", errors.New("unable to generate mac")
	}

	// Valid token?
	if !hmac.Equal(mac, expected) {
		return "", errors.New("invalid cookie mac")
	}

	expires, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", errors.New("unable to parse cookie expiry")
	}

	// Has it expired?
	if time.Unix(expires, 0).Before(time.Now()) {
		return "", errors.New("cookie has expired")
	}

	// Looks valid
	return parts[2], nil
}

// Create cookie hmac
func cookieSignature(email, expires, domain, sessionSecret string) string {
	hash := hmac.New(sha256.New, []byte(sessionSecret))
	hash.Write([]byte(domain))
	hash.Write([]byte(email))
	hash.Write([]byte(expires))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
