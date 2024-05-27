package model

type Post struct {
	ID              int    `json:"id" db:"id"`
	Title           string `json:"title" db:"title"`
	Content         string `json:"auth" db:"content"`
	CommentsEnabled bool   `json:"commentsEnabled" db:"comments_enabled"`
	AuthorUUID      string `json:"authorUuid" db:"author_uuid"`
	CreatedAt       string `json:"createdAt" db:"created_at"`
}
