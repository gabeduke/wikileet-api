package store

import "gorm.io/gorm"

const (
	contextKeyUserEmail   = "email"
	contextKeyUserID      = "user_id"
	contextKeyWorkspace   = "workspace"
	contextKeyWorkspaceID = "workspace_id"
	headerUser            = "X-User"
	headerWorkspace       = "X-Workspace"
)

type Workspace struct {
	gorm.Model
	Name        string
	Description string
	Users       []*User `gorm:"many2many:user_workspaces;"`
}

type User struct {
	gorm.Model
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	Items      []Item       `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Workspaces []*Workspace `json:"workspaces" gorm:"many2many:user_workspaces;"`
}

type Item struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	UserID      uint   `json:"user_id"`
	WorkspaceID uint   `json:"workspace_id"`
}

type createItemInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url"`
	UserEmail   string `json:"user_email"`
	Workspace   string `json:"workspace"`
}

type updateItemInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
