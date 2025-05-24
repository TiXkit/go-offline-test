package dto

type Quote struct {
	ID         int    `json:"id"`
	Text       string `json:"quote"`
	AuthorName string `json:"author"`
}
