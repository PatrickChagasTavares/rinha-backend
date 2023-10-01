package handlers

import (
	"os"

	"github.com/patrickchagastavares/rinha-backend/internal/controllers"
	"github.com/patrickchagastavares/rinha-backend/internal/handlers/people"
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
	people.New(opts.Router, opts.Ctrl)

	if os.Getenv("ENV") != "production" {
		swagger.New(opts.Router)
	}
}
