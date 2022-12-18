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
	v1.GET("/items", app.FindItems)
	v1.POST("/items", app.CreateItem)
	v1.GET("/items/:id", app.FindItem)
	v1.PATCH("/items/:id", app.UpdateItem)
	v1.DELETE("/items/:id", app.DeleteItem)

	// Start and run the server
	r.Run()
}
