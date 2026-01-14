package model

import (
	"time"

	"gorm.io/gorm"
)

// Testproduct represents the testproduct entity in the database
type Testproduct struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Testproduct model
func (Testproduct) TableName() string {
	return "testproducts"
}
