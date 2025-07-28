package models

import "time"

// Category represents a classification for a transaction
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;unique" json:"name"`
	ParentID  *uint     `json:"parentId,omitempty"`
	Parent    *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
