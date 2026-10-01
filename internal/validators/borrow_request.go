package validators

import "time"

type BorrowBookRequest struct {
	BookID  uint       `json:"book_id" binding:"required"`
	DueDate *time.Time `json:"due_date"`
}
