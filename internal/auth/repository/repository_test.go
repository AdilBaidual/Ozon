package repository

import (
	"Service/internal/auth"
	"context"
	"log"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"Service/pkg/containers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthRepoTestSuite struct {
	suite.Suite
	pgContainer  *containers.PostgresContainer
	repository   *Repository
	testUserUUID string
	ctx          context.Context
}

func (suite *AuthRepoTestSuite) SetupSuite() {
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
}

func (suite *AuthRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating repository container: %s", err)
	}
}

func (suite *AuthRepoTestSuite) TestCreateUser() {
	t := suite.T()

	userUUID, err := suite.repository.CreateUserRepo(suite.ctx, auth.CreateUserParams{
		Email:     "test@test.ru",
		FirstName: "FistNameTest",
		Password:  "testpass",
	})
	assert.NoError(t, err)
	suite.testUserUUID = userUUID
}

func (suite *AuthRepoTestSuite) TestGetUserByEmail() {
	t := suite.T()

	user, err := suite.repository.GetUserByEmailRepo(suite.ctx, "test@test.ru")
	assert.NoError(t, err)
	assert.Equal(t, suite.testUserUUID, user.UUID)
	assert.Equal(t, "test@test.ru", user.Email)
	assert.Equal(t, "FistNameTest", user.FirstName)
	assert.Equal(t, "testpass", user.Password)
}

func (suite *AuthRepoTestSuite) TestGetUserByUUID() {
	t := suite.T()

	user, err := suite.repository.GetUserByUUIDRepo(suite.ctx, auth.UserUUID{UUID: suite.testUserUUID})
	assert.NoError(t, err)
	assert.Equal(t, suite.testUserUUID, user.UUID)
	assert.Equal(t, "test@test.ru", user.Email)
	assert.Equal(t, "FistNameTest", user.FirstName)
	assert.Equal(t, "testpass", user.Password)
}

func TestAuthRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRepoTestSuite))
}
