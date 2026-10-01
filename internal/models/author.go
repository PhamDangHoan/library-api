package models

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:150;not null" json:"name"`
	Bio  string `gorm:"type:text" json:"bio"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Books []Book `gorm:"many2many:book_authors;" json:"books,omitempty"`
}
