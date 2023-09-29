package controllers

import (
	"github.com/patrickchagastavares/rinha-backend/internal/controllers/people"
	"github.com/patrickchagastavares/rinha-backend/internal/services"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type (
	Container struct {
		People people.IController
	}

	Options struct {
		Srv *services.Container
		Log logger.Logger
	}
)

func New(opts Options) *Container {
	return &Container{
		People: people.New(opts.Srv, opts.Log),
	}
}
