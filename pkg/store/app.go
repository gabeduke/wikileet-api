package store

import (
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var _ AppI = (*App)(nil)

type App struct {
	DB *gorm.DB
}

type AppI interface {
	GetHealthHandler() (healthcheck.Handler, error)
	GetUserMiddleware() gin.HandlerFunc
}

//	@BasePath	/api/v1

// GetUsers GET /users
// returns the list of items associated with the user
//
//	@Summary		Get users
//	@Schemes		http https
//	@Description	List users
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user_email	query		string	false	"search by email address"	Format(email)
//	@Success		200			{string}	model.Item
//	@Router			/users [get]
func (a *App) GetUsers(c *gin.Context) {
	var users []User
	a.DB.Find(&users).Where("email = ?", c.Query("user_email"))
	c.JSON(http.StatusOK, gin.H{"data": users})
}

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
//	@Success		200			{string}	model.Item
//	@Router			/items [get]
func (a *App) GetItems(c *gin.Context) {
	user := &User{}

	userEmail := c.DefaultQuery("user_email", c.GetString(contextKeyUserEmail))
	a.DB.Preload("Items").Where("email = ?", userEmail).First(&user)

	c.JSON(http.StatusOK, gin.H{"data": user.Items})
}

// GetUserItems GET /users/:id/items
// returns the list of items associated with the user
//
//	@Summary		Get items
//	@Schemes		http https
//	@Description	List items by user email.
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	model.Item
//	@Router			/workspaces/:id/users/:id/items [get]
func (a *App) GetUserItems(c *gin.Context) {
	user := &User{}

	a.DB.Preload("Items").Where("id = ?", c.Param("user")).First(&user)

	c.JSON(http.StatusOK, gin.H{"data": user.Items})
}

// CreateItem POST /items
// Create new Item
//
//	@Summary		Create item
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
//	@Success		200			{string}	model.Item
//	@Router			/workspaces/:id/users/:id/items [post]
func (a *App) CreateItem(c *gin.Context) {
	workspaceId, found := c.Get(contextKeyWorkspaceID)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing workspace ID"})
		return
	}

	if _, ok := workspaceId.(uint); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	userEmail, found := c.GetQuery("user_email")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user"})
		return
	}

	//w64, err := strconv.ParseUint(workspace, 10, 64)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace id"})
	//	return
	//}

	user := &User{}
	if err := a.DB.Where("email = ?", userEmail).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to find User"})
		return
	}

	// Validate input
	var input createItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create item
	item := Item{
		Name:        input.Name,
		Description: input.Description,
		URL:         input.URL,
		UserID:      user.ID,
		WorkspaceID: workspaceId.(uint),
	}

	a.DB.Create(&item)

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// CreateUserItem POST /users/:id/items
// Create new Item
//
//	@Summary		Create item
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
//	@Success		200			{string}	model.Item
//	@Router			/workspaces/:id/users/:id/items [post]
func (a *App) CreateUserItem(c *gin.Context) {
	workspace := c.Param("workspace")
	user := c.Param("user")

	w64, err := strconv.ParseUint(workspace, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace id"})
		return
	}

	u64, err := strconv.ParseUint(user, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	// Validate input
	var input createItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create item
	item := Item{
		Name:        input.Name,
		Description: input.Description,
		URL:         input.URL,
		UserID:      uint(u64),
		WorkspaceID: uint(w64),
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
//	@Success		200	{string}	model.Item
//	@Router			/items/:id [get]
func (a *App) GetItem(c *gin.Context) { // Get model if exist
	var item Item

	if err := a.DB.Where("id = ?", c.Param("item")).First(&item).Error; err != nil {
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
//	@Success		200	{string}	model.Item
//	@Router			/items/:id [patch]
func (a *App) UpdateItem(c *gin.Context) {
	// Get model if exist
	var item Item
	if err := a.DB.Where("id = ?", c.Param("item")).First(&item).Error; err != nil {
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
//	@Success		200	{string}	model.Item
//	@Router			/items/:id [delete]
func (a *App) DeleteItem(c *gin.Context) {
	// Get model if exist
	var item Item
	if err := a.DB.Where("id = ?", c.Param("item")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	a.DB.Delete(&item)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (a *App) GetWorkspaceUsers(c *gin.Context) {
	var users []User
	a.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateUser POST /users
// Create new User
//
//	@Summary		Creat user
//	@Schemes		http https
//	@Description	Create a new user.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user_email	query		string	false	"associate item with user"	Format(email)
//	@Param			name		query		string	false	"Username"					Format(email)
//	@Success		200			{string}	model.Item
//	@Router			/users [post]
func (a *App) CreateUser(c *gin.Context) {
	user := &User{}
	workspace := &Workspace{Name: "default"}

	// Validate input
	var input createUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Email = c.Query("user_email")
	user.Name = c.Query("name")
	user.Workspaces = append(user.Workspaces, workspace)

	a.DB.Where("email = ?", user.Email).First(&user)

	a.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetUser GET /users/:id
// Get new User
//
//	@Summary		Get user
//	@Schemes		http https
//	@Description	Get user by id.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	model.Item
//	@Router			/users/:id [get]
func (a *App) GetUser(c *gin.Context) {
	var user User
	if err := a.DB.Where("id = ?", c.Param("user")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetWorkspaces GET /workspaces
// Get all Workspaces
//
//	@Summary		Get workspaces
//	@Schemes		http https
//	@Description	Get all workspaces.
//	@Tags			workspace
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	model.Item
//	@Router			/workspaces [get]
func (a *App) GetWorkspaces(c *gin.Context) {
	var workspaces []Workspace
	a.DB.Find(&workspaces)
	c.JSON(http.StatusOK, gin.H{"data": workspaces})
}

// CreateWorkspace POST /workspaces
// Create new Workspace
//
//	@Summary		Creat workspace
//	@Schemes		http https
//	@Description	Create a new workspace.
//	@Tags			workspace
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	model.Item
//	@Router			/workspaces [post]
func (a *App) CreateWorkspace(c *gin.Context) {
	workspace := &Workspace{}

	// Validate input
	var input createWorkspaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.DB.Create(&workspace)

	c.JSON(http.StatusOK, gin.H{"data": workspace})
}

// GetWorkspace GET /workspaces/:id
// Get new Workspace
//
//	@Summary		Get workspace
//	@Schemes		http https
//	@Description	Get workspace by id.
//	@Tags			workspace
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	model.Item
//	@Router			/workspaces/:id [get]
func (a *App) GetWorkspace(c *gin.Context) {
	var workspace Workspace
	if err := a.DB.Where("id = ?", c.Param("workspace")).First(&workspace).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": workspace})
}
