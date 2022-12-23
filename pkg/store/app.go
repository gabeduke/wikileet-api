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
	GetItems(c *gin.Context)
	CreateItem(c *gin.Context)
	GetItem(c *gin.Context)
	UpdateItem(c *gin.Context)
	DeleteItem(c *gin.Context)
}

//	@BasePath	/api/v1

// GetItems GET /items
// returns the list of items associated with the user
//
//	@Summary		Get items
//	@Schemes		http https
//	@Description	List items by user email.
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Param			user_email	query		string	false	"search by email address"	Format(email)
//	@Param			X-User		header		string	true	"user email"				default(api_docs@leetserve.com)
//	@Param			X-Workspace	header		string	true	"user workspace"			default(default)
//	@Success		200			{string}	model.Item
//	@Router			/items [get]
func (a *App) GetItems(c *gin.Context) {
	user := &User{}

	userEmail := c.DefaultQuery("user_email", c.GetString(contextKeyUserEmail))
	a.DB.Preload("Items").Where("email = ?", userEmail).First(&user)

	c.JSON(http.StatusOK, gin.H{"data": user.Items})
}

// CreateItem POST /items
// Create new Item
//
//	@Summary		Creat item
//	@Schemes		http https
//	@Description	Create a new item.
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Param			user_email	query		string	false	"associate item with user"		Format(email)
//	@Param			workspace	query		string	false	"associate item with workspace"	Format(email)
//	@Param			name		query		string	false	"Item name"						Format(email)
//	@Param			description	query		string	false	"Item description"				Format(email)
//	@Param			url			query		string	false	"Item help URL"					Format(email)
//	@Param			X-User		header		string	true	"user email"					default(api_docs@leetserve.com)
//	@Param			X-Workspace	header		string	true	"user workspace"				default(default)
//	@Success		200			{string}	model.Item
//	@Router			/items [post]
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

// GetItem GET /Item/:id
// Get new Item
//
//	@Summary		Get item
//	@Schemes		http https
//	@Description	Get item by id.
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Param			X-User		header		string	true	"user email"		default(api_docs@leetserve.com)
//	@Param			X-Workspace	header		string	false	"user workspace"	default(default)
//	@Success		200			{string}	model.Item
//	@Router			/item/:id [get]
func (a *App) GetItem(c *gin.Context) { // Get model if exist
	var item Item

	if err := a.DB.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// UpdateItem PATCH /Item/:id
// Update a Item
//
//	@Summary		Update item
//	@Schemes		http https
//	@Description	Get item by id.
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Param			X-User		header		string	true	"user email"		default(api_docs@leetserve.com)
//	@Param			X-Workspace	header		string	false	"user workspace"	default(default)
//	@Success		200			{string}	model.Item
//	@Router			/item/:id [patch]
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
//
//	@Summary		Delete item
//	@Schemes		http https
//	@Description	Delete item by id.
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Param			X-User		header		string	true	"user email"		default(api_docs@leetserve.com)
//	@Param			X-Workspace	header		string	false	"user workspace"	default(default)
//	@Success		200			{string}	model.Item
//	@Router			/item/:id [delete]
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

func (a *App) GetUsers(c *gin.Context) {
	var users []User
	a.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}
