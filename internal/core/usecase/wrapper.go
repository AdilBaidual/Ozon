package usecase

import (
	"Service/internal/core/model"
	"context"

	"go.opentelemetry.io/otel"
)

func (u *UC) CreatePost(ctx context.Context, params model.NewPost, authorUUID string) (model.Post, error) {
	var (
		response model.Post
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("CreatePost")
	c, span := tracer.Start(ctx, "CreatePostUC()")
	defer span.End()
	if response, err = u.CreatePostUC(c, params, authorUUID); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}

func (u *UC) CreateComment(ctx context.Context, params model.NewComment, authorUUID string) (model.Comment, error) {
	var (
		response model.Comment
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("CreateComment")
	c, span := tracer.Start(ctx, "CreateCommentUC()")
	defer span.End()
	if response, err = u.CreateCommentUC(c, params, authorUUID); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}

func (u *UC) GetPosts(ctx context.Context) ([]*model.Post, error) {
	var (
		response []*model.Post
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("GetPosts")
	c, span := tracer.Start(ctx, "GetPostsUC()")
	defer span.End()
	if response, err = u.GetPostsUC(c); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}

func (u *UC) GetCommentsByPostId(ctx context.Context, postId int) ([]*model.Comment, error) {
	var (
		response []*model.Comment
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("GetCommentsByPostId")
	c, span := tracer.Start(ctx, "GetCommentsByPostIdUC()")
	defer span.End()
	if response, err = u.GetCommentsByPostIdUC(c, postId); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}

func (u *UC) GetCommentsByParentId(ctx context.Context, parentId int) ([]*model.Comment, error) {
	var (
		response []*model.Comment
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("GetCommentsByParentId")
	c, span := tracer.Start(ctx, "GetCommentsByParentIdUC()")
	defer span.End()
	if response, err = u.GetCommentsByParentIdUC(c, parentId); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}

func (u *UC) GetCommentById(ctx context.Context, id int) (model.Comment, error) {
	var (
		response model.Comment
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("GetCommentById")
	c, span := tracer.Start(ctx, "GetCommentByIdUC()")
	defer span.End()
	if response, err = u.GetCommentByIdUC(c, id); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}

func (u *UC) GetPostById(ctx context.Context, id int) (model.Post, error) {
	var (
		response model.Post
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("GetPostById")
	c, span := tracer.Start(ctx, "GetPostByIdUC()")
	defer span.End()
	if response, err = u.GetPostByIdUC(c, id); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}
