package graph

import (
	"Service/internal/core/model"
	"context"
	"encoding/base64"
	"errors"
	"strconv"
)

type contentResolver struct{ *Resolver }

//nolint:funlen
func (c *contentResolver) Comments(ctx context.Context, obj *model.Content,
	first *int, cursor *string, deep *bool) (*model.CommentsConnection, error) {
	var (
		response model.CommentsConnection
		data     []*model.Comment
		err      error
	)

	loaders := For(ctx)

	limit := 10
	if first != nil {
		limit = *first
	}

	var startID int
	if cursor != nil && *cursor != "" {
		decodedCursor, err := base64.StdEncoding.DecodeString(*cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		startID, err = strconv.Atoi(string(decodedCursor))
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		if deep != nil && *deep {
			data, err = loaders.SubCommentLoader.Load(startID)
			if err != nil {
				return nil, err
			}
		} else {
			comment, err := c.coreUC.GetCommentById(ctx, startID)
			if err != nil {
				return nil, err
			}
			if comment.ParentID == nil || *comment.ParentID == 0 {
				data, err = loaders.PostCommentLoader.Load(obj.ID)
				if err != nil {
					return nil, err
				}
			} else {
				data, err = loaders.SubCommentLoader.Load(startID)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		data, err = loaders.PostCommentLoader.Load(obj.ID)
		if err != nil {
			return nil, err
		}
	}

	filteredData := []*model.Comment{}
	for _, comment := range data {
		if deep != nil && !*deep {
			if comment.ID > startID {
				filteredData = append(filteredData, comment)
			}
		} else {
			filteredData = append(filteredData, comment)
		}
	}
	filteredDataLen := len(filteredData)
	if len(filteredData) > limit {
		filteredData = filteredData[:limit]
	}

	edges := make([]*model.CommentsEdge, len(filteredData))
	for i, comment := range filteredData {
		edge := &model.CommentsEdge{
			Cursor: base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(comment.ID))),
			Node:   comment,
		}
		subs, err := loaders.SubCommentLoader.Load(comment.ID)
		if err == nil && len(subs) != 0 {
			edge.HasSubComments = true
		}

		edges[i] = edge
	}

	response.Edges = edges
	if len(edges) > 0 {
		startCursor := edges[0].Cursor
		endCursor := edges[len(edges)-1].Cursor
		hasNextPage := filteredDataLen > limit

		response.PageInfo = &model.PageInfo{
			StartCursor: startCursor,
			EndCursor:   endCursor,
			HasNextPage: &hasNextPage,
		}
	} else {
		hasNextPage := false
		response.PageInfo = &model.PageInfo{
			StartCursor: "",
			EndCursor:   "",
			HasNextPage: &hasNextPage,
		}
	}

	return &response, nil
}

func (r *Resolver) Content() ContentResolver {
	return &contentResolver{r}
}
