package handlers

import (
	"github.com/patrickchagastavares/rinha-backend/internal/controllers"
	"github.com/patrickchagastavares/rinha-backend/internal/handlers/swagger"
	"github.com/patrickchagastavares/rinha-backend/pkg/httpRouter"
)

type (
	Options struct {
		Ctrl   *controllers.Container
		Router httpRouter.Router
	}
)

func NewRouter(opts Options) {
	swagger.New(opts.Router)
}
