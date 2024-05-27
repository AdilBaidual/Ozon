package repository

import (
	"Service/internal/auth"
	"context"

	"github.com/jmoiron/sqlx"
)

const (
	UserTable = " users "
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUserRepo(ctx context.Context, params auth.CreateUserParams) (string, error) {
	var userUUID string

	createUserQuery := `
		INSERT INTO` + UserTable + ` 
		(email, first_name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING uuid;
	`

	err := r.db.QueryRowContext(ctx, createUserQuery, params.Email, params.FirstName, params.Password).Scan(&userUUID)
	if err != nil {
		return "", err
	}
	return userUUID, nil
}

func (r *Repository) GetUserByEmailRepo(ctx context.Context, email string) (auth.User, error) {
	query := `
		SELECT 
			uuid, email, first_name, password_hash
		FROM` + UserTable + `
		WHERE email = $1;
	`
	var data auth.User

	err := r.db.GetContext(ctx, &data, query, email)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *Repository) GetUserByUUIDRepo(ctx context.Context, params auth.UserUUID) (auth.User, error) {
	query := `
		SELECT 
			uuid, email, first_name, password_hash
		FROM` + UserTable + `
		WHERE uuid = $1;
	`
	var data auth.User
	err := r.db.GetContext(ctx, &data, query, params.UUID)
	if err != nil {
		return data, err
	}
	return data, nil
}
