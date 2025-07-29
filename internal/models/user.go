package models

import "time"

// User represents a user of the personal finance tracker
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:100;not null;unique" json:"username" validate:"required,min=3,max=50"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"` // Store hashed password, omit from JSON
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
