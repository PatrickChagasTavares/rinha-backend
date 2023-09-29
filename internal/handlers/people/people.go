package people

import (
	"github.com/patrickchagastavares/rinha-backend/internal/controllers"
	"github.com/patrickchagastavares/rinha-backend/pkg/httpRouter"
)

func New(router httpRouter.Router, ctrl *controllers.Container) {

	router.Post("/pessoas", ctrl.People.Create)
	// router.Get("/pessoas", nil)
	// router.Get("/pessoas/:id", nil)
	// router.Get("/contagem-pessoas", nil)
}
