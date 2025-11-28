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
	hasNo     bool
	no        uint
}

func (m *Model) SetNo(no uint) {
	m.no = no
	m.hasNo = true
}

func (m *Model) GetHasNo() bool {
	return m.hasNo
}
