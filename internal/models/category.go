package models

import "gorm.io/gorm"

// Category represents a classification for a transaction
type Category struct {
	gorm.Model
	Name     string    `gorm:"size:100;not null;unique" json:"name" validate:"required,min=2,max=100"`
	ParentID *uint     `json:"parentId,omitempty"`
	Parent   *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	UserID   uint      `json:"userId"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
}
