package redis

import (
	"Service/internal/auth"
	authRepo "Service/internal/auth/repository"
	"Service/internal/core/model"
	"Service/pkg/containers"
	redisStorage "Service/pkg/storage/redis"
	"context"
	"github.com/go-redis/redis/v8"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CoreRedisRepoTestSuite struct {
	suite.Suite
	redisContainer *containers.RedisContainer
	repository     *Repository
	testUserUUID   string
	ctx            context.Context
}

func (suite *CoreRedisRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	redisContainer, err := containers.NewRedisContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.redisContainer = redisContainer
	connection, err := redisContainer.ConnectionString(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	client, err := redisStorage.NewRedisClient(&redis.Options{
		Addr:     connection,
		Password: "",
		DB:       0,
	})
	if err != nil {
		log.Fatal(err)
	}

	repository := NewRepository(client)
	if err != nil {
		log.Fatal(err)
	}
	suite.repository = repository

	pgContainer, err := containers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sqlx.Connect("pgx", pgContainer.ConnectionString)
	repo := authRepo.NewRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	uuid, err := repo.CreateUserRepo(suite.ctx, auth.CreateUserParams{
		Email:     "test@test.ru",
		FirstName: "test",
		Password:  "testpass",
	})

	assert.NoError(suite.T(), err)
	suite.testUserUUID = uuid
}

func (suite *CoreRedisRepoTestSuite) TearDownSuite() {
	if err := suite.redisContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating repository container: %s", err)
	}
}

func (suite *CoreRedisRepoTestSuite) TestCreatePostRepo() {
	t := suite.T()
	postID, err := suite.repository.CreatePostRepo(suite.ctx, model.NewPost{
		Title:           "Test Post",
		Content:         "Test Content",
		CommentsEnabled: true,
	}, suite.testUserUUID)
	assert.NoError(t, err)
	assert.NotZero(t, postID)
}

func (suite *CoreRedisRepoTestSuite) TestGetPostByIdRepo() {
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

func (suite *CoreRedisRepoTestSuite) TestCreateCommentRepo() {
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

func (suite *CoreRedisRepoTestSuite) TestGetCommentByIdRepo() {
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

func (suite *CoreRedisRepoTestSuite) TestGetPostsRepo() {
	t := suite.T()

	posts, err := suite.repository.GetPostsRepo(suite.ctx)
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.NotZero(t, len(posts))
}

func (suite *CoreRedisRepoTestSuite) TestGetCommentsByPostIdRepo() {
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

func (suite *CoreRedisRepoTestSuite) TestGetCommentsByParentIdRepo() {
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
	suite.Run(t, new(CoreRedisRepoTestSuite))
}
