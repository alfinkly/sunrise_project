package repository

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	IP      string `gorm:"uniqueIndex;not null"`
	Country string `gorm:"not null"`
	City    string `gorm:"not null"`
}
