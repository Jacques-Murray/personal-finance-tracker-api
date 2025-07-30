package models

import (
	"time"

	"gorm.io/gorm"
)

// TransactionType defines the type of transaction: 'income' or 'expense'
type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

// Transaction represents an income or expense record
type Transaction struct {
	gorm.Model
	Description string          `gorm:"type:text" json:"description,omitempty"`
	Amount      float64         `gorm:"type:numeric(10,2);not null" json:"amount" validate:"required,gt=0"`
	Type        TransactionType `gorm:"type:varchar(7);not null" json:"type" validate:"required,oneof=income expense"`
	Date        time.Time       `gorm:"not null" json:"date" validate:"required"`
	CategoryID  uint            `json:"categoryId" validate:"required"`
	Category    Category        `gorm:"foreignKey:CategoryID" json:"category"`
	UserID      uint            `json:"userId"`
	User        User            `gorm:"foreignKey:UserID" json:"user"`
}
