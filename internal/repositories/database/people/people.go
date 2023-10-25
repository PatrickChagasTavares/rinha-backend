package people

import (
	"context"
	"errors"
	"strings"

	"github.com/patrickchagastavares/rinha-backend/internal/entities"
)

type IRepository interface {
	Create(ctx context.Context, person entities.PersonRequest) error
	FindByID(ctx context.Context, id string) (entities.Person, error)
	FindBySearch(ctx context.Context, query string) ([]entities.Person, error)
	Count(ctx context.Context) (uint, error)

	IsErrNotFound(err error) bool
	IsErrDuplicate(err error) bool
}

var (
	createPersonQuery = `INSERT INTO people 
		(id, name, nick_name, birth_date, stack, search, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	findByIdQuery = `
	SELECT id, name, nick_name, birth_date, stack
	FROM people
	WHERE id = $1;`

	nicknameExistQuery = `SELECT EXISTS (SELECT true FROM people WHERE nick_name=$1);`

	findBySearchQuery = `
	SELECT id, name, nick_name, birth_date, stack
	FROM people
	WHERE search LIKE '%' || $1 || '%' LIMIT 50;`

	countQuery = `SELECT count(*) AS count FROM people;`

	// Err
	errNotFound            = errors.New("person is not found or deleted")
	errQueryCreate         = errors.New("problem to create person")
	errNickNameQueryCreate = errors.New("failed create because nick_name already used")
	errQueryFindId         = errors.New("failed to found person")
	errQueryFind           = errors.New("problem to find people")
	errQueryCount          = errors.New("problem to count people")
)

func makeQueryTsvector(q string) string {
	var query strings.Builder
	qSplit := strings.Split(strings.ReplaceAll(q, "  ", ""), " ")
	for idx := range qSplit {
		if len(qSplit)-1 == idx {
			query.WriteString(qSplit[idx] + ":*")
			continue
		}
		query.WriteString(qSplit[idx] + ":* & ")
	}
	return query.String()
}
