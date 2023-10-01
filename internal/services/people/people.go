package people

import (
	"context"

	"github.com/patrickchagastavares/rinha-backend/internal/entities"
	"github.com/patrickchagastavares/rinha-backend/internal/repositories"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type (
	IService interface {
		Create(ctx context.Context, person entities.PersonRequest) (string, error)
		FindByID(ctx context.Context, id string) (entities.Person, error)
		FindByText(ctx context.Context, query string) ([]entities.Person, error)
		Count(ctx context.Context) (uint, error)
	}

	services struct {
		repositories *repositories.Container
		log          logger.Logger
	}
)

func New(repo *repositories.Container, log logger.Logger) IService {
	return &services{repositories: repo, log: log}
}

func (srv *services) Create(ctx context.Context, person entities.PersonRequest) (string, error) {
	nicknameUsed, err := srv.repositories.Database.People.FindNickNameExist(ctx, person.NickName)
	if err != nil {
		return "", err
	}

	if nicknameUsed {
		return "", ErrNicknameAlreadyUsed
	}

	person.PreSave()

	if err := srv.repositories.Database.People.Create(ctx, person); err != nil {
		return "", err
	}

	return person.ID, nil
}

func (srv *services) FindByID(ctx context.Context, id string) (person entities.Person, err error) {
	person, err = srv.repositories.Database.People.FindByID(ctx, id)
	if err != nil {
		if srv.repositories.Database.People.IsErrNotFound(err) {
			err = ErrPersonNotFound
			return
		}
		return
	}
	return
}

func (srv *services) FindByText(ctx context.Context, query string) (people []entities.Person, err error) {
	people, err = srv.repositories.Database.People.FindBySearch(ctx, query)
	return
}

func (srv *services) Count(ctx context.Context) (count uint, err error) {
	count, err = srv.repositories.Database.People.Count(ctx)
	return
}
