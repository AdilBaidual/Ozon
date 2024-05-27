package redis

import (
	"Service/internal/core/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) CreatePostRepo(ctx context.Context, params model.NewPost, authorUUID string) (int, error) {
	id, err := r.client.Incr(ctx, "post_id").Result()
	if err != nil {
		return 0, err
	}

	post := model.Post{
		ID:              int(id),
		Title:           params.Title,
		Content:         params.Content,
		CommentsEnabled: params.CommentsEnabled,
		AuthorUUID:      authorUUID,
		CreatedAt:       time.Now().Format(time.RFC3339),
	}

	postData, err := json.Marshal(post)
	if err != nil {
		return 0, err
	}

	err = r.client.Set(ctx, fmt.Sprintf("post:%d", id), postData, 0).Err()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *Repository) GetPostByIdRepo(ctx context.Context, id int) (model.Post, error) {
	var post model.Post
	postData, err := r.client.Get(ctx, fmt.Sprintf("post:%d", id)).Result()
	if err != nil {
		return post, err
	}

	err = json.Unmarshal([]byte(postData), &post)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (r *Repository) CreateCommentRepo(ctx context.Context, params model.NewComment, authorUUID string) (int, error) {
	id, err := r.client.Incr(ctx, "comment_id").Result()
	if err != nil {
		return 0, err
	}

	parentID := 0
	if params.ParentID != nil {
		parentID = *params.ParentID
	}

	comment := model.Comment{
		ID:         int(id),
		PostID:     params.PostID,
		ParentID:   &parentID,
		AuthorUUID: authorUUID,
		Content:    params.Content,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	commentData, err := json.Marshal(comment)
	if err != nil {
		return 0, err
	}

	err = r.client.Set(ctx, fmt.Sprintf("comment:%d", id), commentData, 0).Err()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *Repository) GetCommentByIdRepo(ctx context.Context, id int) (model.Comment, error) {
	var comment model.Comment
	commentData, err := r.client.Get(ctx, fmt.Sprintf("comment:%d", id)).Result()
	if err != nil {
		return comment, err
	}

	err = json.Unmarshal([]byte(commentData), &comment)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

func (r *Repository) GetPostsRepo(ctx context.Context) ([]*model.Post, error) {
	keys, err := r.client.Keys(ctx, "post:*").Result()
	if err != nil {
		return nil, err
	}

	var posts []*model.Post
	for _, key := range keys {
		postData, err := r.client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var post model.Post
		err = json.Unmarshal([]byte(postData), &post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *Repository) GetCommentsByPostIdRepo(ctx context.Context, postId int) ([]*model.Comment, error) {
	keys, err := r.client.Keys(ctx, "comment:*").Result()
	if err != nil {
		return nil, err
	}

	var comments []*model.Comment
	for _, key := range keys {
		commentData, err := r.client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var comment model.Comment
		err = json.Unmarshal([]byte(commentData), &comment)
		if err != nil {
			return nil, err
		}

		if comment.PostID == postId && (comment.ParentID == nil || *comment.ParentID == 0) {
			comments = append(comments, &comment)
		}
	}

	return comments, nil
}

func (r *Repository) GetCommentsByParentIdRepo(ctx context.Context, parentId int) ([]*model.Comment, error) {
	keys, err := r.client.Keys(ctx, "comment:*").Result()
	if err != nil {
		return nil, err
	}

	var comments []*model.Comment
	for _, key := range keys {
		commentData, err := r.client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var comment model.Comment
		err = json.Unmarshal([]byte(commentData), &comment)
		if err != nil {
			return nil, err
		}

		if comment.ParentID != nil && *comment.ParentID == parentId {
			comments = append(comments, &comment)
		}
	}

	return comments, nil
}
