package services

import (
	"github.com/patrickchagastavares/rinha-backend/internal/repositories"
	"github.com/patrickchagastavares/rinha-backend/internal/services/people"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type (
	Container struct {
		People people.IService
	}

	Options struct {
		Repo *repositories.Container
		Log  logger.Logger
	}
)

func New(opts Options) *Container {
	return &Container{
		People: people.New(opts.Repo, opts.Log),
	}
}
