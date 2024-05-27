package model

type Comment struct {
	ID         int    `json:"id" db:"id"`
	PostID     int    `json:"postId" db:"post_id"`
	ParentID   *int   `json:"parentId" db:"parent_id,omitempty"`
	AuthorUUID string `json:"authorUuid" db:"author_uuid"`
	Content    string `json:"content" db:"content"`
	CreatedAt  string `json:"createdAt" db:"created_at"`
}
