package graph

import (
	"Service/internal/core"
	"Service/internal/core/model"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

const commentLoaderKey = "commentLoader"

type Loaders struct {
	PostCommentLoader CommentLoader
	SubCommentLoader  CommentLoader
}

func DataLoader(coreUC core.UseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		commentLoader := &Loaders{
			PostCommentLoader: CommentLoader{
				maxBatch: 100,
				wait:     1 * time.Microsecond,
				fetch: func(postIds []int) ([][]*model.Comment, []error) {
					result := make([][]*model.Comment, 0, len(postIds))
					for _, id := range postIds {
						comments, err := coreUC.GetCommentsByPostId(ctx, id)
						if err != nil {
							return nil, []error{err}
						}
						result = append(result, comments)
					}
					return result, nil
				},
			},
			SubCommentLoader: CommentLoader{
				maxBatch: 100,
				wait:     1 * time.Microsecond,
				fetch: func(parentIds []int) ([][]*model.Comment, []error) {
					result := make([][]*model.Comment, 0, len(parentIds))
					for _, id := range parentIds {
						comments, err := coreUC.GetCommentsByParentId(ctx, id)
						if err != nil {
							return nil, []error{err}
						}
						result = append(result, comments)
					}
					return result, nil
				},
			},
		}
		c := context.WithValue(ctx.Request.Context(), commentLoaderKey, commentLoader)
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(commentLoaderKey).(*Loaders)
}
