package graph

import (
	"Service/internal/core/model"
	"context"
)

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	return r.coreUC.GetPosts(ctx)
}

func (r *queryResolver) Content(ctx context.Context, postId int) (*model.Content, error) {
	post, err := r.coreUC.GetPostById(ctx, postId)
	if err != nil {
		return nil, err
	}
	var response model.Content
	response.ID = post.ID
	response.Title = &post.Title
	response.Content = &post.Content
	response.CommentsEnabled = &post.CommentsEnabled
	response.AuthorUUID = &post.AuthorUUID
	response.CreatedAt = &post.CreatedAt
	return &response, nil
}
