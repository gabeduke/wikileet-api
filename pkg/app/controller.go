package app

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// Controller is the controller for the app
type Controller struct {
	DB *gorm.DB
}

// GetProfile GET /profile
// Get a user's profile
//
//	@Summary		Get profile
//	@Schemes		http https
//	@Description	Get logged-in user profile
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	app.User
//	@Failure		400	{string}	error
//	@Router			/profile [get]
func (a *Controller) GetProfile(c *gin.Context) {
	var user User
	jwt.ExtractClaims(c)
	claims := jwt.ExtractClaims(c)

	if err := a.DB.Where("email = ?", claims[identityKey]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
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
//	@Success		200			{object}	app.Item
//	@Failure		400			{string}	error
//	@Router			/users [get]
func (a *Controller) GetUsers(c *gin.Context) {
	var users []User
	if err := a.DB.Find(&users).Where("email = ?", c.Query("user_email")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
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
//	@Success		200			{array}		app.Item
//	@Failure		400			{string}	error
//	@Router			/items [get]
func (a *Controller) GetItems(c *gin.Context) {
	user := &User{}

	claims := jwt.ExtractClaims(c)

	userEmail := c.DefaultQuery("user_email", claims[identityKey].(string))
	if err := a.DB.Preload("Items").Where("email = ?", userEmail).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, user.Items)
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
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{array}		app.Item
//	@Failure		400	{string}	error
//	@Router			/users/{id}/items [get]
func (a *Controller) GetUserItems(c *gin.Context) {
	user := &User{}

	if err := a.DB.Preload("Items").Where("id = ?", c.Param("user")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, user.Items)
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
//	@Param			user_email	query		string			false	"associate item with user"	Format(email)
//	@Param			data		body		createItemInput	false	"Item data"
//	@Success		200			{object}	app.Item
//	@Failure		400			{string}	error
//	@Router			/items [post]
func (a *Controller) CreateItem(c *gin.Context) {

	userEmail, found := c.GetQuery("user_email")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user"})
		return
	}

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
	}

	if err := a.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
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
//	@Param			data	body		createItemInput	false	"User data"
//	@Param			user_id	path		string			true	"User ID"	Format(uuid)
//	@Success		200		{object}	app.Item
//	@Failure		400		{string}	error
//	@Router			/users/{id}/items [post]
func (a *Controller) CreateUserItem(c *gin.Context) {
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

	if err := a.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
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
//	@Success		200	{object}	app.Item
//	@Failure		400	{string}	error
//	@Router			/items/{id} [get]
func (a *Controller) GetItem(c *gin.Context) { // Get model if exist
	var item Item

	if err := a.DB.Where("id = ?", c.Param("item")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, item)
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
//	@Param			data	body		updateItemInput	false	"Item data"
//	@Success		200		{object}	app.Item
//	@Failure		400		{string}	error
//	@Router			/items/{id} [patch]
func (a *Controller) UpdateItem(c *gin.Context) {
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

	if err := a.DB.Model(&item).Updates(input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
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
//	@Param			id	path	string	true	"Item ID"
//	@Success		204
//	@Failure		400	{string}	error
//	@Router			/items/{id} [delete]
func (a *Controller) DeleteItem(c *gin.Context) {
	var item Item
	if err := a.DB.Where("id = ?", c.Param("item")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := a.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetWorkspaceUsers GET /workspaces/:id/users
// Get all users
//
//	@Summary		Get all users
//	@Schemes		http https
//	@Description	Get all users.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Workspace ID"
//	@Success		200	{array}		app.User
//	@Failure		400	{string}	error
//	@Router			/workspaces/{id}/users [get]
func (a *Controller) GetWorkspaceUsers(c *gin.Context) {
	var users []User
	a.DB.Find(&users)
	c.JSON(http.StatusOK, users)
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
//	@Param			data	body		createUserInput	false	"User data"
//	@Success		200		{object}	app.User
//	@Failure		400		{string}	error
//	@Router			/users [post]
func (a *Controller) CreateUser(c *gin.Context) {
	user := &User{}
	workspace := &Workspace{Name: "default"}

	// Validate input
	var input createUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Email = input.Email
	user.Name = input.Name
	user.Workspaces = append(user.Workspaces, workspace)

	if err := a.DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := a.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
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
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	app.User
//	@Failure		400	{string}	error
//	@Router			/users/{id} [get]
func (a *Controller) GetUser(c *gin.Context) {
	var user User
	if err := a.DB.Where("id = ?", c.Param("user")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, user)
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
//	@Success		200	{array}		app.Workspace
//	@Failure		400	{string}	error
//	@Router			/workspaces [get]
func (a *Controller) GetWorkspaces(c *gin.Context) {
	var workspaces []Workspace
	err := a.DB.Find(&workspaces).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workspaces)
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
//	@Param			data	body		createWorkspaceInput	false	"Workspace data"
//	@Success		200		{object}	app.Workspace
//	@Failure		400		{string}	error
//	@Router			/workspaces [post]
func (a *Controller) CreateWorkspace(c *gin.Context) {
	workspace := &Workspace{}

	// Validate input
	var input createWorkspaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Create(&workspace).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workspace)
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
//	@Param			id	path		string	true	"Workspace ID"
//	@Success		200	{object}	app.Workspace
//	@Failure		400	{string}	error
//	@Router			/workspaces/{id} [get]
func (a *Controller) GetWorkspace(c *gin.Context) {
	var workspace Workspace
	if err := a.DB.Where("id = ?", c.Param("workspace")).First(&workspace).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, workspace)
}
