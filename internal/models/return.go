package models

import "time"

type Return struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BorrowID   uint      `gorm:"not null;uniqueIndex" json:"borrow_id"`
	ReturnedAt time.Time `gorm:"not null" json:"returned_at"`
	CreatedAt  time.Time `json:"created_at"`

	Borrow Borrow `gorm:"foreignKey:BorrowID" json:"borrow,omitempty"`
}
