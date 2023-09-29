package people

import (
	"net/http"

	"github.com/patrickchagastavares/rinha-backend/internal/entities"
	"github.com/patrickchagastavares/rinha-backend/internal/services"
	"github.com/patrickchagastavares/rinha-backend/pkg/httpRouter"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type (
	IController interface {
		Create(c httpRouter.Context)
	}

	controllers struct {
		srv *services.Container
		log logger.Logger
	}
)

func New(srv *services.Container, log logger.Logger) IController {
	return &controllers{
		srv: srv,
		log: log,
	}
}

// people swagger document
// @Description Create one person
// @Tags people
// @Accept json
// @Produce json
// @Param house body entities.PersonRequest true "create new person"
// @Success 201 {object} entities.Person
// @Failure 400 {object} entities.HttpErr
// @Failure 409 {object} entities.HttpErr
// @Failure 500
// @Security ApiKeyAuth
// @Router /pessoas [post]
func (ctrl *controllers) Create(c httpRouter.Context) {
	var newPerson entities.PersonRequest
	if err := c.Decode(&newPerson); err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrDecode)
		return
	}

	if err := c.Validate(newPerson); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

}
