package main

import (
	docs "github.com/gabeduke/wikileet-api/docs"
	"github.com/gabeduke/wikileet-api/pkg/config"
	"github.com/gin-contrib/static"
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

//	@host		wikileet.leetserve.com
//	@host		dev.wikileet.leetserve.com
//	@host		test.wikileet.leetserve.com
//	@host		localhost:8080
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

	auth, err := app.GetAuthMiddleware(config.GetSessionSecret(), config.GetDomain(), config.GetZone(), false)
	if err != nil {
		log.Fatal(err)
	}

	// Configure router
	docs.SwaggerInfo.BasePath = "/api/v1"
	r := gin.New()
	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/live", "/ready"),
		gin.Recovery(),
		CORSMiddleware(),
	)

	r.GET("/live", gin.WrapF(health.LiveEndpoint))
	r.GET("/ready", gin.WrapF(health.ReadyEndpoint))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/login", auth.LoginHandler)
	r.GET("/refresh_token", auth.RefreshHandler)

	v1 := r.Group("/api/v1")
	v1.Use(app.GetUserMiddleware(), auth.MiddlewareFunc())

	// Register items routes
	v1.GET("/items", app.GetItems)
	v1.POST("/items", app.CreateItem)
	v1.GET("/items/:item", app.GetItem)
	v1.PATCH("/items/:item", app.UpdateItem)
	v1.DELETE("/items/:item", app.DeleteItem)

	// Register workspace routes
	workspaces := v1.Group("/workspaces")
	workspaces.GET("/", app.GetWorkspaces)
	workspaces.POST("/", app.CreateWorkspace)
	workspaces.GET("/:workspace", app.GetWorkspace)

	// Register user routes
	users := v1.Group("/users")
	users.GET("/", app.GetWorkspaceUsers)
	users.POST("/", app.CreateUser)
	users.GET("/:user", app.GetUser)

	userItems := users.Group("/:user/items")
	userItems.GET("/", app.GetUserItems)
	userItems.POST("/", app.CreateUserItem)
	userItems.GET("/:item", app.GetItem)
	userItems.PATCH("/:item", app.UpdateItem)
	userItems.DELETE("/:item", app.DeleteItem)

	r.Use(static.Serve("/", static.LocalFile("./frontend/wikileet-ui/dist", false)))

	// Start and run the server
	r.Run()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User, X-Workspace")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
