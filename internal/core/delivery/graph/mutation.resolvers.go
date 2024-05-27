package graph

import (
	"Service/internal/auth"
	"Service/internal/core/model"
	"context"
	"errors"
)

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	data, user, err := r.authUC.SignUp(ctx, auth.SignUpParams{
		Email:     input.Email,
		FirstName: input.FirstName,
		Password:  input.Password,
	})
	if err != nil {
		return nil, err
	}
	response := &model.AuthResponse{
		AuthToken: &model.AuthToken{AccessToken: data.Access, RefreshToken: data.Refresh},
		User:      &model.User{UUID: user.UUID, Email: user.Email, FirstName: user.FirstName},
	}
	return response, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error) {
	data, user, err := r.authUC.SignIn(ctx, auth.SignInParams{Email: input.Email, Password: input.Password})
	if err != nil {
		return nil, err
	}
	response := &model.AuthResponse{
		AuthToken: &model.AuthToken{AccessToken: data.Access, RefreshToken: data.Refresh},
		User:      &model.User{UUID: user.UUID, Email: user.Email, FirstName: user.FirstName},
	}
	return response, nil
}

// Refresh is the resolver for the refresh field.
func (r *mutationResolver) Refresh(ctx context.Context, input model.RefreshInput) (*model.AuthResponse, error) {
	uuid, ok := ctx.Value("uuid").(string)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	data, user, err := r.authUC.Refresh(ctx, auth.RefreshParams{UserUUID: auth.UserUUID{UUID: uuid}, Refresh: input.Refresh})
	if err != nil {
		return nil, err
	}
	response := &model.AuthResponse{
		AuthToken: &model.AuthToken{AccessToken: data.Access, RefreshToken: input.Refresh},
		User:      &model.User{UUID: user.UUID, Email: user.Email, FirstName: user.FirstName},
	}
	return response, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	uuid, ok := ctx.Value("uuid").(string)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	post, err := r.coreUC.CreatePost(ctx, input, uuid)
	return &post, err
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	uuid, ok := ctx.Value("uuid").(string)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	post, err := r.coreUC.CreateComment(ctx, input, uuid)
	return &post, err
}
