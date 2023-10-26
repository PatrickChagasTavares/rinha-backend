package people

import (
	"fmt"
	"net/http"

	"github.com/patrickchagastavares/rinha-backend/internal/entities"
	"github.com/patrickchagastavares/rinha-backend/internal/services"
	"github.com/patrickchagastavares/rinha-backend/pkg/httpRouter"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type (
	IController interface {
		Create(c httpRouter.Context)
		Find(c httpRouter.Context)
		FindByID(c httpRouter.Context)
		Count(c httpRouter.Context)
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
// @Failure 422 {object} entities.HttpErr
// @Failure 500
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

	c.SetHeader("Location", "/pessoas/1")
	c.JSON(http.StatusCreated, nil)

	// id, err := ctrl.srv.People.Create(c.Context(), newPerson)
	// if err != nil {
	// 	if errors.Is(err, people.ErrNicknameAlreadyUsed) {
	// 		c.JSON(
	// 			http.StatusUnprocessableEntity,
	// 			entities.NewHttpErr(http.StatusUnprocessableEntity, err.Error(), nil))
	// 		return
	// 	}

	// 	c.JSON(http.StatusInternalServerError, nil)
	// 	return
	// }

	// // newPerson.ID = id
	// c.SetHeader("Location", "/pessoas/"+id)
	// c.JSON(http.StatusCreated, nil)
	// // c.JSON(http.StatusCreated, newPerson.ToPerson())
	return
}

// people swagger document
// @Description find one person
// @Tags people
// @Accept json
// @Produce json
// @Param	t	query string true "search"
// @Success 200 {object} entities.Person
// @Failure 400 {object} entities.HttpErr
// @Failure 404 {object} entities.HttpErr
// @Failure 500
// @Router /pessoas [get]
func (ctrl *controllers) Find(c httpRouter.Context) {
	query := c.GetQuery("t")
	if len(query) == 0 {
		c.JSON(
			http.StatusBadRequest,
			nil,
			// entities.NewHttpErr(http.StatusBadRequest, "t is required", nil),
		)
		return
	}
	c.JSON(http.StatusOK, []entities.Person{})

	// people, err := ctrl.srv.People.FindByText(c.Context(), query)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, nil)
	// 	return
	// }

	// c.JSON(http.StatusOK, people)
	return
}

// people swagger document
// @Description find one person
// @Tags people
// @Accept json
// @Produce json
// @Param	id path string true "person id"
// @Success 200 {object} entities.Person
// @Failure 400 {object} entities.HttpErr
// @Failure 404 {object} entities.HttpErr
// @Failure 500
// @Router /pessoas/:id [get]
func (ctrl *controllers) FindByID(c httpRouter.Context) {
	c.JSON(http.StatusOK, entities.Person{})
	// person, err := ctrl.srv.People.FindByID(c.Context(), id)
	// if err != nil {
	// 	if errors.Is(err, people.ErrPersonNotFound) {
	// 		c.JSON(
	// 			http.StatusNotFound,
	// 			entities.NewHttpErr(http.StatusBadRequest, err.Error(), nil),
	// 		)
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, nil)
	// 	return
	// }

	// c.JSON(http.StatusOK, person)
	return
}

// people swagger document
// @Description find one person
// @Tags people
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /contagem-pessoas [get]
func (ctrl *controllers) Count(c httpRouter.Context) {
	count, err := ctrl.srv.People.Count(c.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.String(http.StatusOK, fmt.Sprint(count))
	return
}
