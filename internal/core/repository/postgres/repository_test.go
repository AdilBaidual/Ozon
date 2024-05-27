package postgres

import (
	"Service/internal/core/model"
	"Service/pkg/containers"
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CorePostgresRepoTestSuite struct {
	suite.Suite
	pgContainer  *containers.PostgresContainer
	repository   *Repository
	testUserUUID string
	ctx          context.Context
}

func (suite *CorePostgresRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := containers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	db, err := sqlx.Connect("pgx", suite.pgContainer.ConnectionString)
	repository := NewRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	suite.repository = repository

	var userUUID string

	createUserQuery := `
		INSERT INTO users 
		(email, first_name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING uuid;
	`
	err = suite.repository.db.QueryRowContext(suite.ctx, createUserQuery, "test@test.ru", "test", "testpass").Scan(&userUUID)
	assert.NoError(suite.T(), err)
	suite.testUserUUID = userUUID
}

func (suite *CorePostgresRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating repository container: %s", err)
	}
}

func (suite *CorePostgresRepoTestSuite) TestCreatePostRepo() {
	t := suite.T()
	postID, err := suite.repository.CreatePostRepo(suite.ctx, model.NewPost{
		Title:           "Test Post",
		Content:         "Test Content",
		CommentsEnabled: true,
	}, suite.testUserUUID)
	assert.NoError(t, err)
	assert.NotZero(t, postID)
}

func (suite *CorePostgresRepoTestSuite) TestGetPostByIdRepo() {
	t := suite.T()

	postID, err := suite.repository.CreatePostRepo(suite.ctx, model.NewPost{
		Title:           "Test Post",
		Content:         "Test Content",
		CommentsEnabled: true,
	}, suite.testUserUUID)
	assert.NoError(t, err)

	post, err := suite.repository.GetPostByIdRepo(suite.ctx, postID)
	assert.NoError(t, err)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "Test Content", post.Content)
	assert.True(t, post.CommentsEnabled)
}

func (suite *CorePostgresRepoTestSuite) TestCreateCommentRepo() {
	t := suite.T()

	postID, err := suite.repository.CreatePostRepo(suite.ctx, model.NewPost{
		Title:           "Test Post",
		Content:         "Test Content",
		CommentsEnabled: true,
	}, suite.testUserUUID)
	assert.NoError(t, err)

	commentID, err := suite.repository.CreateCommentRepo(suite.ctx, model.NewComment{
		PostID:  postID,
		Content: "Test Content",
	}, suite.testUserUUID)
	assert.NoError(t, err)
	assert.NotZero(t, postID)
	assert.NotZero(t, commentID)
}

func (suite *CorePostgresRepoTestSuite) TestGetCommentByIdRepo() {
	t := suite.T()

	postID, err := suite.repository.CreatePostRepo(suite.ctx, model.NewPost{
		Title:           "Test Post",
		Content:         "Test Content",
		CommentsEnabled: true,
	}, suite.testUserUUID)
	assert.NoError(t, err)

	commentID, err := suite.repository.CreateCommentRepo(suite.ctx, model.NewComment{
		PostID:  postID,
		Content: "Test Content",
	}, suite.testUserUUID)
	assert.NoError(t, err)
	assert.NotZero(t, postID)
	assert.NotZero(t, commentID)

	comment, err := suite.repository.GetCommentByIdRepo(suite.ctx, commentID)
	assert.NoError(t, err)
	assert.Equal(t, commentID, comment.ID)
	assert.Equal(t, "Test Content", comment.Content)
}

func (suite *CorePostgresRepoTestSuite) TestGetPostsRepo() {
	t := suite.T()

	posts, err := suite.repository.GetPostsRepo(suite.ctx)
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.NotZero(t, len(posts))
}

func (suite *CorePostgresRepoTestSuite) TestGetCommentsByPostIdRepo() {
	t := suite.T()

	postID, err := suite.repository.CreatePostRepo(suite.ctx, model.NewPost{
		Title:           "Test Post",
		Content:         "Test Content",
		CommentsEnabled: true,
	}, suite.testUserUUID)
	assert.NoError(t, err)

	commentID, err := suite.repository.CreateCommentRepo(suite.ctx, model.NewComment{
		PostID:  postID,
		Content: "Test Content",
	}, suite.testUserUUID)
	assert.NoError(t, err)
	assert.NotZero(t, postID)
	assert.NotZero(t, commentID)

	comments, err := suite.repository.GetCommentsByPostIdRepo(suite.ctx, postID)
	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.NotZero(t, len(comments))
}

func (suite *CorePostgresRepoTestSuite) TestGetCommentsByParentIdRepo() {
	t := suite.T()

	postID, err := suite.repository.CreatePostRepo(suite.ctx, model.NewPost{
		Title:           "Test Post",
		Content:         "Test Content",
		CommentsEnabled: true,
	}, suite.testUserUUID)
	assert.NoError(t, err)

	parentID, err := suite.repository.CreateCommentRepo(suite.ctx, model.NewComment{
		PostID:  postID,
		Content: "Test Comment",
	}, suite.testUserUUID)
	assert.NoError(t, err)

	_, err = suite.repository.CreateCommentRepo(suite.ctx, model.NewComment{
		PostID:   postID,
		ParentID: &parentID,
		Content:  "Test Comment",
	}, suite.testUserUUID)
	assert.NoError(t, err)

	comments, err := suite.repository.GetCommentsByParentIdRepo(suite.ctx, parentID)
	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.NotZero(t, len(comments))
}

func TestCoreRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CorePostgresRepoTestSuite))
}
