package auth

import "context"

type Repo interface {
	CreateUser(ctx context.Context, params CreateUserParams) (string, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByUUID(ctx context.Context, params UserUUID) (User, error)
}
