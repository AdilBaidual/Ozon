package core

import (
	"Service/internal/core/model"
	"context"
)

type UseCase interface {
	CreatePost(ctx context.Context, params model.NewPost, authorUUID string) (model.Post, error)
	CreateComment(ctx context.Context, params model.NewComment, authorUUID string) (model.Comment, error)
	GetPosts(ctx context.Context) ([]*model.Post, error)
	GetCommentsByPostId(ctx context.Context, postId int) ([]*model.Comment, error)
	GetCommentsByParentId(ctx context.Context, parentId int) ([]*model.Comment, error)
	GetCommentById(ctx context.Context, id int) (model.Comment, error)
	GetPostById(ctx context.Context, id int) (model.Post, error)
}
