package models

// Author - Бизнес-модель автора.
type Author struct {
	ID         int
	AuthorName string
	Quotes     []Quote
}
