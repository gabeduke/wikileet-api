package controllers

import (
	"net/http"

	"github.com/gabeduke/wikileet-api/models"
	"github.com/gin-gonic/gin"
)

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

// GET /items
// Get all items
func FindItems(c *gin.Context) {
	var items []models.Item
	models.DB.Find(&items)

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// POST /items
// Create new item
func CreateItem(c *gin.Context) {
	// Validate input
	var input CreateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create item
	item := models.Item{Name: input.Name, Description: input.Description}
	models.DB.Create(&item)

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// GET /item/:id
// Find a item
func FindItem(c *gin.Context) { // Get model if exist
	var item models.Item

	if err := models.DB.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// PATCH /item/:id
// Update a item
func UpdateItem(c *gin.Context) {
	// Get model if exist
	var item models.Item
	if err := models.DB.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&item).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// DELETE /items/:id
// Delete a item
func DeleteItem(c *gin.Context) {
	// Get model if exist
	var item models.Item
	if err := models.DB.Where("id = ?", c.Param("id")).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&item)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
