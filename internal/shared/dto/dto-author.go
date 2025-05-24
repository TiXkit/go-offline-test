package dto

type Author struct {
	ID         int      `json:"id"`
	AuthorName string   `json:"author"`
	Quotes     []*Quote `json:"quotes"`
}
