package models

import "time"

// TransactionType defines the type of transaction: 'income' or 'expense'
type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

// Transaction represents an income or expense record
type Transaction struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Description string          `gorm:"type:text" json:"description,omitempty"`
	Amount      float64         `gorm:"type:numeric(10,2);not null" json:"amount" validate:"required,gt=0"`
	Type        TransactionType `gorm:"type:varchar(7);not null" json:"type" validate:"required,oneof=income expense"`
	Date        time.Time       `gorm:"not null" json:"date" validate:"required"`
	CategoryID  uint            `json:"categoryId" validate:"required"`
	Category    Category        `gorm:"foreignKey:CategoryID" json:"category"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
