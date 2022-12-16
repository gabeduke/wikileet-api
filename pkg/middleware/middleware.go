package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func User() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the header
		user := c.GetHeader("X-User")
		c.Set("user", user)

		//
		uuid := uuid.New()
		c.Set("uuid", uuid)
		logrus.Tracef("The request with uuid %s is started", uuid)
		c.Next()
		logrus.Tracef("The request with uuid %s is served", uuid)
	}
}
