package services

import (
	"github.com/patrickchagastavares/rinha-backend/internal/repositories"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type (
	Container struct {
	}

	Options struct {
		Repo *repositories.Container
		Log  logger.Logger
	}
)

func New(opts Options) *Container {
	return &Container{

	}
}
