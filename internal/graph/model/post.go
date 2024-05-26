package model

type Post struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Content         string  `json:"auth"`
	CommentsEnabled bool    `json:"comments_enabled"`
	AuthorUUID      string  `json:"author_uuid"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       *string `json:"updated_at,omitempty"`
}
