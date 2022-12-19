package main

import (
	docs "github.com/gabeduke/wikileet-api/docs"
	"github.com/gabeduke/wikileet-api/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

//	@title			Wikileet API
//	@version		1.0
//	@description	Wikileet gift lists API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Wikileet Support
//	@contact.url	http://www.leetserve.com/support
//	@contact.email	support@leetserve.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		https://wikileet.leetserve.com
//	@host		https://dev.wikileet.leetserve.com
//	@host		http://localhost:8080
//	@BasePath	/api/v1

//	@securitydefinitions.oauth2.application					OAuth2Application
//	@tokenUrl												https://oauth2.googleapis.com/token
//	@scope.https://www.googleapis.com/auth/userinfo.email	See your primary Google Account email address
func main() {
	logrus.Info("Starting Wikileet API")

	// Load config
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logrus.Infof("%+v", config)

	// Connect to database
	app, err := config.ParseDialerConfig().Dial()
	if err != nil {
		log.Fatal(err)
	}

	// Configure healthcheck
	health, err := app.GetHealthHandler()
	if err != nil {
		log.Fatal(err)
	}

	// Configure router
	docs.SwaggerInfo.BasePath = "/api/v1"
	r := gin.New()
	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/live", "/ready"),
		gin.Recovery(),
		app.GetUserMiddleware(),
	)

	r.GET("/live", gin.WrapF(health.LiveEndpoint))
	r.GET("/ready", gin.WrapF(health.ReadyEndpoint))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Register routes
	v1 := r.Group("/api/v1")

	// Register routes
	v1.GET("/items", app.GetItems)
	v1.POST("/items", app.CreateItem)
	v1.GET("/items/:id", app.GetItem)
	v1.PATCH("/items/:id", app.UpdateItem)
	v1.DELETE("/items/:id", app.DeleteItem)

	// Start and run the server
	r.Run()
}
