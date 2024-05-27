package postgres

import (
	"Service/internal/core/model"
	"context"

	"go.opentelemetry.io/otel"
)

func (r *Repository) CreatePost(ctx context.Context, params model.NewPost, authorUUID string) (int, error) {
	var (
		id  int
		err error
	)
	tracer := otel.GetTracerProvider().Tracer("CreatePost")
	c, span := tracer.Start(ctx, "CreatePostRepo()")
	defer span.End()
	if id, err = r.CreatePostRepo(c, params, authorUUID); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return id, err
}

func (r *Repository) GetPostById(ctx context.Context, id int) (model.Post, error) {
	var (
		post model.Post
		err  error
	)
	tracer := otel.GetTracerProvider().Tracer("GetPostById")
	c, span := tracer.Start(ctx, "GetPostByIdRepo()")
	defer span.End()
	if post, err = r.GetPostByIdRepo(c, id); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return post, err
}

func (r *Repository) CreateComment(ctx context.Context, params model.NewComment, authorUUID string) (int, error) {
	var (
		id  int
		err error
	)
	tracer := otel.GetTracerProvider().Tracer("CreateComment")
	c, span := tracer.Start(ctx, "CreateCommentRepo()")
	defer span.End()
	if id, err = r.CreateCommentRepo(c, params, authorUUID); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return id, err
}

func (r *Repository) GetCommentById(ctx context.Context, id int) (model.Comment, error) {
	var (
		comment model.Comment
		err     error
	)
	tracer := otel.GetTracerProvider().Tracer("GetCommentById")
	c, span := tracer.Start(ctx, "GetCommentByIdRepo()")
	defer span.End()
	if comment, err = r.GetCommentByIdRepo(c, id); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return comment, err
}

func (r *Repository) GetPosts(ctx context.Context) ([]*model.Post, error) {
	var (
		response []*model.Post
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("GetPosts")
	c, span := tracer.Start(ctx, "GetPostsRepo()")
	defer span.End()
	if response, err = r.GetPostsRepo(c); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}

func (r *Repository) GetCommentsByPostId(ctx context.Context, postId int) ([]*model.Comment, error) {
	var (
		response []*model.Comment
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("GetCommentsByPostId")
	c, span := tracer.Start(ctx, "GetCommentsByPostIdRepo()")
	defer span.End()
	if response, err = r.GetCommentsByPostIdRepo(c, postId); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}

func (r *Repository) GetCommentsByParentId(ctx context.Context, parentId int) ([]*model.Comment, error) {
	var (
		response []*model.Comment
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("GetCommentsByParentId")
	c, span := tracer.Start(ctx, "GetCommentsByParentIdRepo()")
	defer span.End()
	if response, err = r.GetCommentsByParentIdRepo(c, parentId); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}
