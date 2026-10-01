package services

import "errors"

var (
	ErrBookNotFound    = errors.New("book not found")
	ErrDuplicateISBN   = errors.New("book with this ISBN already exists")
	ErrInvalidBookData = errors.New("invalid book data")
)
