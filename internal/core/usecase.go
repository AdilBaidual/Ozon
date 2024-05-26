package core

import (
	"Service/internal/graph/model"
	"context"
)

type UseCase interface {
	CreatePost(ctx context.Context, params model.NewPost) (model.Post, error)
	CreateComment(ctx context.Context, params model.NewComment) (model.Comment, error)
	GetPosts(ctx context.Context) ([]*model.Post, error)
	GetComments(ctx context.Context, params model.GetCommentsInput) ([]*model.Comment, error)
}
