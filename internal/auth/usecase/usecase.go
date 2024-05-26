package usecase

import (
	"Service/internal/auth"
	"Service/pkg/paseto"
	"Service/pkg/secure"
	valkeyStorage "Service/pkg/storage/valkey"
	"context"
	"errors"
	"go.uber.org/zap"
)

type UC struct {
	lg           *zap.Logger
	paseto       *paseto.Paseto
	tokenStorage *valkeyStorage.Storage
	repo         auth.Repo
}

func NewUseCase(logger *zap.Logger, paseto *paseto.Paseto,
	tokenStorage *valkeyStorage.Storage, repo auth.Repo) *UC {
	return &UC{
		lg:           logger,
		paseto:       paseto,
		tokenStorage: tokenStorage,
		repo:         repo,
	}
}

func (u *UC) SignUpUC(ctx context.Context, params auth.SignUpParams) (auth.SignUpResponse, auth.User, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	var (
		response auth.SignUpResponse
		user     auth.User
	)
	passwordHash, err := secure.CalculateHash(params.Password)
	if err != nil {
		logger.Error("calculating hash", zap.Error(err))
		return response, user, err
	}
	userUUID, err := u.repo.CreateUser(ctx, auth.CreateUserParams{
		Email:     params.Email,
		FirstName: params.FirstName,
		Password:  passwordHash,
	})
	if err != nil {
		logger.Error("creating user db", zap.Error(err))
		return response, user, err
	}

	access, err := u.paseto.GenerateAccessToken(userUUID)
	if err != nil {
		logger.Error("generating access token", zap.Error(err))
		return response, user, err
	}

	refresh, err := u.paseto.GenerateRefreshToken(userUUID)
	if err != nil {
		logger.Error("generating refresh token", zap.Error(err))
		return response, user, err
	}

	err = u.tokenStorage.Set(userUUID+"_refresh", []byte(refresh), paseto.RefreshTTL)
	if err != nil {
		logger.Error("local saving refresh token", zap.Error(err))
		return response, user, err
	}

	response.Access = access
	response.Refresh = refresh

	user, err = u.repo.GetUserByUUID(ctx, auth.UserUUID{UUID: userUUID})
	if err != nil {
		logger.Error("getting user by uuid", zap.Error(err))
		return response, auth.User{}, err
	}

	return response, user, nil
}

func (u *UC) SignInUC(ctx context.Context, params auth.SignInParams) (auth.SignInResponse, auth.User, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	var (
		response auth.SignInResponse
	)

	user, err := u.repo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		logger.Error("getting user from db", zap.Error(err))
		return response, auth.User{}, err
	}

	compare := secure.CompareHash(user.Password, params.Password)
	if !compare {
		logger.Error("validating password")
		return response, auth.User{}, err
	}

	access, err := u.paseto.GenerateAccessToken(user.UUID)
	if err != nil {
		logger.Error("generating access token", zap.Error(err))
		return response, auth.User{}, err
	}

	refresh, err := u.paseto.GenerateRefreshToken(user.UUID)
	if err != nil {
		logger.Error("generating refresh token", zap.Error(err))
		return response, auth.User{}, err
	}

	err = u.tokenStorage.Set(user.UUID+"_refresh", []byte(refresh), paseto.RefreshTTL)
	if err != nil {
		logger.Error("saving refresh token", zap.Error(err))
		return response, auth.User{}, err
	}

	response.Access = access
	response.Refresh = refresh

	return response, user, nil
}

func (u *UC) GetUserByUUIDUC(ctx context.Context, params auth.UserUUID) (auth.User, error) {
	return u.repo.GetUserByUUID(ctx, params)
}

func (u *UC) SignOutUC(ctx context.Context, params auth.SignOutParams) error {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	if params.UUID == "" {
		err := errors.New("empty uuid")
		logger.Error("validating params")
		return err
	}

	err := u.tokenStorage.Set(params.Access+"_blacklist", []byte(params.Access), paseto.AccessTTL)
	if err != nil {
		logger.Error("saving access token to blacklist", zap.Error(err))
		return err
	}
	err = u.tokenStorage.Delete(params.UUID + "_refresh")
	if err != nil {
		logger.Error("deleting refresh token from tokenStorage", zap.Error(err))
		return err
	}
	return nil
}

func (u *UC) RefreshUC(ctx context.Context, params auth.RefreshParams) (auth.RefreshResponse, auth.User, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	var response auth.RefreshResponse

	user, err := u.repo.GetUserByUUID(ctx, params.UserUUID)
	if err != nil {
		logger.Error("getting user by uuid", zap.Error(err))
		return response, auth.User{}, err
	}

	val, err := u.tokenStorage.Get(params.UUID + "_refresh")
	if err != nil {
		logger.Error("getting refresh token from tokenStorage", zap.Error(err))
		return response, auth.User{}, err
	}
	if string(val) != params.Refresh {
		logger.Error("comparing refresh tokens")
		return response, auth.User{}, err
	}
	access, err := u.paseto.GenerateAccessToken(params.UUID)
	if err != nil {
		logger.Error("generating access token", zap.Error(err))
		return response, auth.User{}, err
	}
	response.Access = access
	return response, user, nil
}
