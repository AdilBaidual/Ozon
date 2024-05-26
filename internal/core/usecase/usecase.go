package usecase

import (
	"Service/internal/core"
	"go.uber.org/zap"
)

type UC struct {
	lg   *zap.Logger
	repo core.Repo
}

func NewUseCase(logger *zap.Logger, repo core.Repo) *UC {
	return &UC{
		lg:   logger,
		repo: repo,
	}
}
