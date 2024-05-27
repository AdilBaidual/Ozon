package repository

import (
	"Service/internal/auth"
	"context"

	"go.opentelemetry.io/otel"
)

func (r *Repository) CreateUser(ctx context.Context, params auth.CreateUserParams) (string, error) {
	var (
		uuid string
		err  error
	)
	tracer := otel.GetTracerProvider().Tracer("CreateUser")
	c, span := tracer.Start(ctx, "CreateUserRepo()")
	defer span.End()
	if uuid, err = r.CreateUserRepo(c, params); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return uuid, err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (auth.User, error) {
	var (
		data auth.User
		err  error
	)
	tracer := otel.GetTracerProvider().Tracer("GetUserByEmail")
	c, span := tracer.Start(ctx, "GetUserByEmailRepo()")
	defer span.End()
	if data, err = r.GetUserByEmailRepo(c, email); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return data, err
}

func (r *Repository) GetUserByUUID(ctx context.Context, params auth.UserUUID) (auth.User, error) {
	var (
		data auth.User
		err  error
	)
	tracer := otel.GetTracerProvider().Tracer("GetUserByUUID")
	c, span := tracer.Start(ctx, "GetUserByUUIDRepo()")
	defer span.End()
	if data, err = r.GetUserByUUIDRepo(c, params); err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	return data, err
}
