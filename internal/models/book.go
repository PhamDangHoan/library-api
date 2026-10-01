package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	ISBN        string     `gorm:"size:20;uniqueIndex;not null" json:"isbn"`
	Description string     `gorm:"type:text" json:"description"`
	PublishedAt *time.Time `json:"published_at"`

	Quantity  int `gorm:"not null;default:0" json:"quantity"`
	Available int `gorm:"not null;default:0" json:"available"`

	Authors []Author `gorm:"many2many:book_authors;" json:"authors,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
