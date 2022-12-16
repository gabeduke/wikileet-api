package store

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type CreateItemInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type UpdateItemInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// FindItems GET /items
// Get all items
func (a *App) FindItems(c *gin.Context) {
	var items []Item
	a.db.Find(&items)

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// CreateItem POST /items
// Create new item
func (a *App) CreateItem(c *gin.Context) {
	// Validate input
	var input CreateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create item
	item := Item{Name: input.Name, Description: input.Description}
	a.db.Create(&item)

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// FindItem GET /item/:id
// Find a item
func (a *App) FindItem(c *gin.Context) { // Get model if exist
	var item Item

	if err := a.db.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// UpdateItem PATCH /item/:id
// Update a item
func (a *App) UpdateItem(c *gin.Context) {
	// Get model if exist
	var item Item
	if err := a.db.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.db.Model(&item).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// DeleteItem DELETE /items/:id
// Delete a item
func (a *App) DeleteItem(c *gin.Context) {
	// Get model if exist
	var item Item
	if err := a.db.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	a.db.Delete(&item)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
