package auth

type User struct {
	UUID      string `db:"uuid"`
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	Password  string `db:"password_hash"`
}

type SignUpParams struct {
	Email     string
	FirstName string
	Password  string
}

type SignUpResponse struct {
	Access  string
	Refresh string
}

type CreateUserParams struct {
	Email     string
	FirstName string
	Password  string
}

type UserUUID struct {
	UUID string
}

type SignInParams struct {
	Email    string
	Password string
}

type SignInResponse struct {
	Access  string
	Refresh string
}

type SignOutParams struct {
	UserUUID
	Access string
}

type RefreshParams struct {
	UserUUID
	Refresh string
}

type RefreshResponse struct {
	Access string
}
