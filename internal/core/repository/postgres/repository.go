package postgres

import (
	"Service/internal/core/model"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	PostTable    = " posts "
	CommentTable = " comments "
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreatePostRepo(ctx context.Context, params model.NewPost, authorUUID string) (int, error) {
	var id int

	createUserQuery := `
		INSERT INTO` + PostTable + ` 
		(title, content, comments_enabled, author_uuid, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	err := r.db.QueryRowContext(ctx, createUserQuery, params.Title, params.Content, params.CommentsEnabled, authorUUID, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetPostByIdRepo(ctx context.Context, id int) (model.Post, error) {
	query := `
		SELECT 
			id, title, content, comments_enabled, author_uuid, created_at
		FROM` + PostTable + `
		WHERE id = $1;
	`
	var data model.Post

	err := r.db.GetContext(ctx, &data, query, id)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *Repository) CreateCommentRepo(ctx context.Context, params model.NewComment, authorUUID string) (int, error) {
	var id int

	createUserQuery := `
		INSERT INTO` + CommentTable + ` 
		(post_id, parent_id, author_uuid, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var parentID int
	if params.ParentID == nil {
		parentID = 0
	} else {
		parentID = *params.ParentID
	}
	err := r.db.QueryRowContext(ctx, createUserQuery, params.PostID, parentID, authorUUID, params.Content, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetCommentByIdRepo(ctx context.Context, id int) (model.Comment, error) {
	query := `
		SELECT 
			id, post_id, parent_id, author_uuid, content, created_at
		FROM` + CommentTable + `
		WHERE id = $1;
	`
	var data model.Comment

	err := r.db.GetContext(ctx, &data, query, id)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *Repository) GetPostsRepo(ctx context.Context) ([]*model.Post, error) {
	query := `
		SELECT 
			id, title, content, comments_enabled, author_uuid, created_at
		FROM` + PostTable + `
		ORDER BY id ASC;
	`
	var data []*model.Post

	err := r.db.SelectContext(ctx, &data, query)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *Repository) GetCommentsByPostIdRepo(ctx context.Context, postId int) ([]*model.Comment, error) {
	query := `
		SELECT 
			id, post_id, parent_id, author_uuid, content, created_at
		FROM` + CommentTable + `
		WHERE post_id = $1 AND parent_id = 0
		ORDER BY id ASC;
	`
	var data []*model.Comment

	err := r.db.SelectContext(ctx, &data, query, postId)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *Repository) GetCommentsByParentIdRepo(ctx context.Context, parentId int) ([]*model.Comment, error) {
	query := `
		SELECT 
			id, post_id, parent_id, author_uuid, content, created_at
		FROM` + CommentTable + `
		WHERE parent_id = $1
		ORDER BY id ASC;
	`
	var data []*model.Comment

	err := r.db.SelectContext(ctx, &data, query, parentId)
	if err != nil {
		return data, err
	}

	return data, nil
}
