package main

import (
	"github.com/gabeduke/wikileet-api/controllers"
	"github.com/gabeduke/wikileet-api/models"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug    bool
	Port     int
	Database *models.DatabaseConfig
}

func main() {
	// Load config
	var config Config
	envconfig.Process("wikileet", &config)

	r := gin.Default()

	models.ConnectDatabase(config.Database)

	v1 := r.Group("/api/v1")
	v1.GET("/items", controllers.FindItems)
	v1.POST("/items", controllers.CreateItem)
	v1.GET("/items/:id", controllers.FindItem)
	v1.PATCH("/items/:id", controllers.UpdateItem)
	v1.DELETE("/items/:id", controllers.DeleteItem)

	r.Run()
}
