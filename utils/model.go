package utils

import (
	"gorm.io/gorm"
	"time"
)

// Model 表实体model
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	HasNo     bool           `json:"has_no"`
	No        uint           `json:"no"`
}
