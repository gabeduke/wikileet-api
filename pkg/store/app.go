package store

import (
	"github.com/heptiolabs/healthcheck"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

var _ AppI = (*App)(nil)

type App struct {
	DB *gorm.DB
}

type AppI interface {
	GetHealthHandler() (healthcheck.Handler, error)
	GetUserMiddleware() gin.HandlerFunc
	FindItems(c *gin.Context)
	CreateItem(c *gin.Context)
	FindItem(c *gin.Context)
	UpdateItem(c *gin.Context)
	DeleteItem(c *gin.Context)
}

//	@BasePath	/api/v1

// FindItems godoc
//	@Summary	get items
//	@Schemes
//	@Description	get items
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Param			user_email	query		string	false	"search by email address"	Format(email)
//	@Success		200			{string}	Helloworld
//	@Router			/items [get]
func (a *App) FindItems(c *gin.Context) {
	user := &User{}

	userEmail := c.DefaultQuery("user_email", c.GetString(contextKeyUserEmail))
	a.DB.Preload("Items").Where("email = ?", userEmail).First(&user)

	c.JSON(http.StatusOK, gin.H{"data": user.Items})
}

// CreateItem POST /items
// Create new Item
func (a *App) CreateItem(c *gin.Context) {
	workspace := &Workspace{}
	user := &User{}
	var inputUserEmail string
	var inputWorkspaceName string

	// Validate input
	var input createItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.UserEmail == "" {
		inputUserEmail = c.DefaultQuery("user_email", c.GetString(contextKeyUserEmail))
	}

	if input.Workspace == "" {
		inputWorkspaceName = c.DefaultQuery("workspace", "default")
	}

	a.DB.Where("email = ?", inputUserEmail).First(&user)
	a.DB.Where("name = ?", inputWorkspaceName).First(&workspace)

	// Create item
	item := Item{
		Name:        input.Name,
		Description: input.Description,
		URL:         input.URL,
		UserID:      user.ID,
		WorkspaceID: workspace.ID,
	}

	a.DB.Create(&item)

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// FindItem GET /Item/:id
// Find a Item
func (a *App) FindItem(c *gin.Context) { // Get model if exist
	var item Item

	if err := a.DB.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// UpdateItem PATCH /Item/:id
// Update a Item
func (a *App) UpdateItem(c *gin.Context) {
	// Get model if exist
	var item Item
	if err := a.DB.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input updateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.DB.Model(&item).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// DeleteItem DELETE /items/:id
// Delete a Item
func (a *App) DeleteItem(c *gin.Context) {
	// Get model if exist
	var item Item
	if err := a.DB.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	a.DB.Delete(&item)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
