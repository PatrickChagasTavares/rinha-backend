package people

import (
	"context"

	"github.com/patrickchagastavares/rinha-backend/internal/entities"
)

type IRepository interface {
	Create(ctx context.Context, person entities.PersonRequest) error
	FindByID(ctx context.Context, id string) (entities.Person, error)
	FindBySearch(ctx context.Context, query string) ([]entities.Person, error)
	FindNickNameExist(ctx context.Context, nickName string) (bool, error)
	Count(ctx context.Context) (uint, error)

	IsErrNotFound(err error) bool
}
