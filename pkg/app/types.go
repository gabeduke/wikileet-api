package app

import "gorm.io/gorm"

// Workspace is a struct for a workspace
type Workspace struct {
	gorm.Model
	Name        string
	Description string
	Users       []*User `gorm:"many2many:user_workspaces;"`
}

// User is a struct for a user
type User struct {
	gorm.Model
	Name       string       `json:"name"`
	Email      string       `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	Items      []Item       `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Workspaces []*Workspace `json:"workspaces" gorm:"many2many:user_workspaces;"`
}

type Item struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url" validate:"url"`
	UserID      uint   `json:"user_id"`
	WorkspaceID uint   `json:"workspace_id"`
}

type createItemInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url" validate:"omitempty,url"`
	UserEmail   string `json:"user_email" validate:"required,email"`
	Workspace   string `json:"workspace"`
}

type updateItemInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url" validate:"omitempty,url"`
}

type createUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required" validate:"email"`
}

type updateUserInput struct {
	Name       string   `json:"name"`
	Email      string   `json:"email" validate:"omitempty,email"`
	Workspaces []string `json:"workspaces"`
}

type createWorkspaceInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
