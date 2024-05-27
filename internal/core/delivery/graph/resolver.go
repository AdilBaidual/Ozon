//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"Service/internal/auth"
	coreUC "Service/internal/core/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authUC auth.UseCase
	coreUC *coreUC.UC
}

func NewResolver(authUC auth.UseCase, coreUC *coreUC.UC) *Resolver {
	return &Resolver{
		authUC: authUC,
		coreUC: coreUC,
	}
}
