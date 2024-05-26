package usecase

import (
	"Service/internal/auth"
	"context"

	"go.opentelemetry.io/otel"
)

func (u *UC) SignUp(ctx context.Context, params auth.SignUpParams) (auth.SignUpResponse, auth.User, error) {
	var (
		response auth.SignUpResponse
		user     auth.User
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("SignUp")
	c, span := tracer.Start(ctx, "SignUpUC()")
	defer span.End()
	if response, user, err = u.SignUpUC(c, params); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, user, err
}

func (u *UC) SignIn(ctx context.Context, params auth.SignInParams) (auth.SignInResponse, auth.User, error) {
	var (
		response auth.SignInResponse
		user     auth.User
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("SignIn")
	c, span := tracer.Start(ctx, "SignInUC()")
	defer span.End()
	if response, user, err = u.SignInUC(c, params); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, user, err
}

func (u *UC) GetUserByUUID(ctx context.Context, params auth.UserUUID) (auth.User, error) {
	var (
		response auth.User
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("GetUserByUUID")
	c, span := tracer.Start(ctx, "GetUserByUUIDUC()")
	defer span.End()
	if response, err = u.GetUserByUUIDUC(c, params); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, err
}

func (u *UC) SignOut(ctx context.Context, params auth.SignOutParams) error {
	var err error
	tracer := otel.GetTracerProvider().Tracer("SignOut")
	c, span := tracer.Start(ctx, "SignOutUC()")
	defer span.End()
	if err = u.SignOutUC(c, params); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return err
}

func (u *UC) Refresh(ctx context.Context, params auth.RefreshParams) (auth.RefreshResponse, auth.User, error) {
	var (
		response auth.RefreshResponse
		user     auth.User
		err      error
	)
	tracer := otel.GetTracerProvider().Tracer("Refresh")
	c, span := tracer.Start(ctx, "RefreshUC()")
	defer span.End()
	if response, user, err = u.RefreshUC(c, params); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return response, user, err
}
