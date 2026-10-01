package models

import "time"

type Borrow struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	BookID     uint       `gorm:"not null;index" json:"book_id"`
	BorrowedAt time.Time  `gorm:"not null" json:"borrowed_at"`
	DueDate    *time.Time `json:"due_date"`
	Status     string     `gorm:"size:20;not null;default:borrowed" json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	User   User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Book   Book    `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Return *Return `gorm:"foreignKey:BorrowID" json:"return,omitempty"`
}
