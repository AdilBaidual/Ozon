package repository

import (
	"api-cp/internal/auth"
	"context"
	"log"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"api-cp/pkg/containers"

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

	userUUID, err := suite.repository.CreateUser(suite.ctx, auth.CreateUserParams{
		Email:      "test@test.ru",
		FirstName:  "FistNameTest",
		LastName:   "LastNameTest",
		Password:   "testpass",
		Phone:      "+77777777777",
		Tg:         "@Test",
		TotpSecret: "test",
	})
	assert.NoError(t, err)
	suite.testUserUUID = userUUID
}

func (suite *AuthRepoTestSuite) TestGetUserConfirmations() {
	t := suite.T()

	confrimations, err := suite.repository.GetUserConfirmations(suite.ctx, auth.UserUUID{UUID: suite.testUserUUID})
	assert.NoError(t, err)
	assert.Equal(t, false, confrimations.Tg)
	assert.Equal(t, false, confrimations.Email)
	assert.Equal(t, false, confrimations.Totp)
	assert.Equal(t, suite.testUserUUID, confrimations.UserUUID)
}

func (suite *AuthRepoTestSuite) TestGetUserByEmail() {
	t := suite.T()

	user, err := suite.repository.GetUserByEmail(suite.ctx, "test@test.ru")
	assert.NoError(t, err)
	assert.Equal(t, suite.testUserUUID, user.UUID)
	assert.Equal(t, "test@test.ru", user.Email)
	assert.Equal(t, "FistNameTest", user.FirstName)
	assert.Equal(t, "LastNameTest", user.LastName)
	assert.Equal(t, "testpass", user.Password)
	assert.Equal(t, "+77777777777", user.Phone)
	assert.Equal(t, "@Test", user.Tg)
	assert.Equal(t, "test", user.TotpSecret)
}

func (suite *AuthRepoTestSuite) TestUpdateUserTotp() {
	t := suite.T()

	userUUID, err := suite.repository.CreateUser(suite.ctx, auth.CreateUserParams{
		Email:      "test2@test.ru",
		FirstName:  "FistNameTest",
		LastName:   "LastNameTest",
		Password:   "testpass",
		Phone:      "+77777777777",
		Tg:         "@Test",
		TotpSecret: "test",
	})
	assert.NoError(t, err)

	err = suite.repository.UpdateUserTotp(suite.ctx, auth.User{
		UUID:       userUUID,
		TotpSecret: "updatedtotp",
	})
	assert.NoError(t, err)

	confrimations, err := suite.repository.GetUserConfirmations(suite.ctx, auth.UserUUID{UUID: userUUID})
	assert.NoError(t, err)
	assert.Equal(t, false, confrimations.Tg)
	assert.Equal(t, false, confrimations.Email)
	assert.Equal(t, false, confrimations.Totp)
	assert.Equal(t, userUUID, confrimations.UserUUID)
}

func (suite *AuthRepoTestSuite) TestGetUserByUUID() {
	t := suite.T()

	user, err := suite.repository.GetUserByUUID(suite.ctx, auth.UserUUID{UUID: suite.testUserUUID})
	assert.NoError(t, err)
	assert.Equal(t, suite.testUserUUID, user.UUID)
	assert.Equal(t, "test@test.ru", user.Email)
	assert.Equal(t, "FistNameTest", user.FirstName)
	assert.Equal(t, "LastNameTest", user.LastName)
	assert.Equal(t, "testpass", user.Password)
	assert.Equal(t, "+77777777777", user.Phone)
	assert.Equal(t, "@Test", user.Tg)
	assert.Equal(t, "test", user.TotpSecret)
}

func (suite *AuthRepoTestSuite) TestUpdateUserConfirmations() {
	t := suite.T()

	err := suite.repository.UpdateUserConfirmations(suite.ctx, auth.UserConfirmations{UserUUID: suite.testUserUUID, Tg: true})
	assert.NoError(t, err)

	confrimations, err := suite.repository.GetUserConfirmations(suite.ctx, auth.UserUUID{UUID: suite.testUserUUID})
	assert.NoError(t, err)
	assert.Equal(t, true, confrimations.Tg)
	assert.Equal(t, false, confrimations.Email)
	assert.Equal(t, false, confrimations.Totp)
	assert.Equal(t, suite.testUserUUID, confrimations.UserUUID)
}

func TestAuthRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRepoTestSuite))
}
