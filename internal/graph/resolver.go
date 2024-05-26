package graph

import "Service/internal/auth"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authUC auth.UseCase
}

func NewResolver(authUC auth.UseCase) *Resolver {
	return &Resolver{
		authUC: authUC,
	}
}
