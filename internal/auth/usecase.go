package auth

import "context"

type UseCase interface {
	SignUp(ctx context.Context, params SignUpParams) (SignUpResponse, User, error)
	SignIn(ctx context.Context, params SignInParams) (SignInResponse, User, error)
	GetUserByUUID(ctx context.Context, params UserUUID) (User, error)
	SignOut(ctx context.Context, params SignOutParams) error
	Refresh(ctx context.Context, params RefreshParams) (RefreshResponse, User, error)
}
