package usecase

import (
	"Service/internal/core/model"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreatePost(ctx context.Context, params model.NewPost, authorUUID string) (int, error) {
	args := m.Called(ctx, params, authorUUID)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) GetPostById(ctx context.Context, id int) (model.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Post), args.Error(1)
}

func (m *MockRepo) CreateComment(ctx context.Context, params model.NewComment, authorUUID string) (int, error) {
	args := m.Called(ctx, params, authorUUID)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) GetCommentById(ctx context.Context, id int) (model.Comment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Comment), args.Error(1)
}

func (m *MockRepo) GetPosts(ctx context.Context) ([]*model.Post, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockRepo) GetCommentsByPostId(ctx context.Context, postId int) ([]*model.Comment, error) {
	args := m.Called(ctx, postId)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockRepo) GetCommentsByParentId(ctx context.Context, parentId int) ([]*model.Comment, error) {
	args := m.Called(ctx, parentId)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func setup(t *testing.T) (*UC, *MockRepo) {
	logger, _ := zap.NewDevelopment()
	mockRepo := new(MockRepo)
	uc := NewUseCase(logger, mockRepo)
	return uc, mockRepo
}

func TestCreatePostUC(t *testing.T) {
	uc, mockRepo := setup(t)

	ctx := context.Background()
	params := model.NewPost{Title: "Test Post", Content: "Test Content", CommentsEnabled: true}
	authorUUID := "author-uuid"
	postID := 1
	expectedPost := model.Post{ID: postID, Title: params.Title, Content: params.Content, CommentsEnabled: params.CommentsEnabled, AuthorUUID: authorUUID}

	mockRepo.On("CreatePost", ctx, params, authorUUID).Return(postID, nil)
	mockRepo.On("GetPostById", ctx, postID).Return(expectedPost, nil)

	post, err := uc.CreatePostUC(ctx, params, authorUUID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)
}

func TestCreateCommentUC(t *testing.T) {
	uc, mockRepo := setup(t)

	ctx := context.Background()
	params := model.NewComment{PostID: 1, Content: "Test Comment"}
	authorUUID := "author-uuid"
	commentID := 1
	expectedComment := model.Comment{ID: commentID, PostID: params.PostID, Content: params.Content, AuthorUUID: authorUUID}

	mockRepo.On("GetPostById", ctx, params.PostID).Return(model.Post{ID: params.PostID, CommentsEnabled: true}, nil)
	mockRepo.On("CreateComment", ctx, params, authorUUID).Return(commentID, nil)
	mockRepo.On("GetCommentById", ctx, commentID).Return(expectedComment, nil)

	comment, err := uc.CreateCommentUC(ctx, params, authorUUID)
	assert.NoError(t, err)
	assert.Equal(t, expectedComment, comment)
}

func TestGetPostsUC(t *testing.T) {
	uc, mockRepo := setup(t)

	ctx := context.Background()
	expectedPosts := []*model.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1", CommentsEnabled: true},
		{ID: 2, Title: "Post 2", Content: "Content 2", CommentsEnabled: true},
	}

	mockRepo.On("GetPosts", ctx).Return(expectedPosts, nil)

	posts, err := uc.GetPostsUC(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)
}

func TestGetCommentsByPostIdUC(t *testing.T) {
	uc, mockRepo := setup(t)

	ctx := context.Background()
	postID := 1
	expectedComments := []*model.Comment{
		{ID: 1, PostID: postID, Content: "Comment 1"},
		{ID: 2, PostID: postID, Content: "Comment 2"},
	}

	mockRepo.On("GetCommentsByPostId", ctx, postID).Return(expectedComments, nil)

	comments, err := uc.GetCommentsByPostIdUC(ctx, postID)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
}

func TestGetCommentsByParentIdUC(t *testing.T) {
	uc, mockRepo := setup(t)

	ctx := context.Background()
	parentID := 1
	expectedComments := []*model.Comment{
		{ID: 1, PostID: 1, ParentID: &parentID, Content: "Comment 1"},
		{ID: 2, PostID: 1, ParentID: &parentID, Content: "Comment 2"},
	}

	mockRepo.On("GetCommentsByParentId", ctx, parentID).Return(expectedComments, nil)

	comments, err := uc.GetCommentsByParentIdUC(ctx, parentID)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
}

func TestGetCommentByIdUC(t *testing.T) {
	uc, mockRepo := setup(t)

	ctx := context.Background()
	commentID := 1
	expectedComment := model.Comment{ID: commentID, PostID: 1, Content: "Test Comment"}

	mockRepo.On("GetCommentById", ctx, commentID).Return(expectedComment, nil)

	comment, err := uc.GetCommentByIdUC(ctx, commentID)
	assert.NoError(t, err)
	assert.Equal(t, expectedComment, comment)
}

func TestGetPostByIdUC(t *testing.T) {
	uc, mockRepo := setup(t)

	ctx := context.Background()
	postID := 1
	expectedPost := model.Post{ID: postID, Title: "Test Post", Content: "Test Content", CommentsEnabled: true}

	mockRepo.On("GetPostById", ctx, postID).Return(expectedPost, nil)

	post, err := uc.GetPostByIdUC(ctx, postID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)
}

func TestCreateCommentUC_PostCommentsDisabled(t *testing.T) {
	uc, mockRepo := setup(t)

	ctx := context.Background()
	params := model.NewComment{PostID: 1, Content: "Test Comment"}
	authorUUID := "author-uuid"

	mockRepo.On("GetPostById", ctx, params.PostID).Return(model.Post{ID: params.PostID, CommentsEnabled: false}, nil)

	comment, err := uc.CreateCommentUC(ctx, params, authorUUID)
	assert.Error(t, err)
	assert.Equal(t, model.Comment{}, comment)
}
