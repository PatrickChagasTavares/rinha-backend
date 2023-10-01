package people

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/patrickchagastavares/rinha-backend/internal/entities"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type repoSqlx struct {
	log    logger.Logger
	writer *sqlx.DB
	reader *sqlx.DB
}

var (
	createPersonQuery = `INSERT INTO people 
		(id, name, nick_name, birth_date, stack, created_at)
		VALUES ($1, $2, $3, $4, $5, $6);`

	findByIdQuery = `
	SELECT id, name, nick_name, birth_date, stack
	FROM people
	WHERE id = $1;`

	nicknameExistQuery = `SELECT EXISTS (SELECT true FROM people WHERE nick_name=$1);`

	findBySearchQuery = `
	SELECT id, name, nick_name, birth_date, stack
	FROM people
	WHERE fts_tokens @@ to_tsquery(unaccent($1)) LIMIT 50;`

	countQuery = `SELECT count(*) AS count FROM people;`

	// Err
	errNotFound    = errors.New("person is not found or deleted")
	errQueryCreate = errors.New("problem to create person")
	errQueryFindId = errors.New("failed to found person")
	errQueryFind   = errors.New("problem to find people")
	errQueryCount  = errors.New("problem to count people")
)

func NewSqlx(log logger.Logger, writer, reader *sqlx.DB) IRepository {
	return &repoSqlx{log: log, writer: writer, reader: reader}
}

func (repo *repoSqlx) Create(ctx context.Context, person entities.PersonRequest) (err error) {
	_, err = repo.writer.ExecContext(
		ctx, createPersonQuery,
		person.ID, person.Name, person.NickName,
		person.BirthDate, person.Stack, person.CreatedAt)
	if err != nil {
		repo.log.ErrorContext(ctx, "people.SqlxRepo.Create data: ", person, ", err: "+err.Error())
		err = errQueryCreate
		return
	}

	return nil
}

func (repo *repoSqlx) FindNickNameExist(ctx context.Context, nickName string) (exists bool, err error) {
	err = repo.reader.GetContext(ctx, &exists, nicknameExistQuery, nickName)
	if err != nil {
		repo.log.ErrorContext(ctx, "people.SqlxRepo.FindNickNameExist nick_name: "+nickName+", err: "+err.Error())
		err = errNotFound
		return
	}

	return exists, nil
}

func (repo *repoSqlx) FindByID(ctx context.Context, id string) (person entities.Person, err error) {
	err = repo.reader.GetContext(ctx, &person, findByIdQuery, id)
	if err != nil {
		repo.log.ErrorContext(ctx, "people.SqlxRepo.FindByID id: "+id+", err: "+err.Error())
		if err == sql.ErrNoRows {
			err = errNotFound
			return
		}
		err = errQueryFind
		return
	}

	return
}

func (repo *repoSqlx) FindBySearch(ctx context.Context, query string) (people []entities.Person, err error) {
	people = make([]entities.Person, 0, 50)
	query = makeQueryTsvector(query)
	err = repo.reader.SelectContext(ctx, &people, findBySearchQuery, query)
	if err != nil && err != sql.ErrNoRows {
		repo.log.ErrorContext(ctx, "people.SqlxRepo.FindBySearch query: "+query+", err: "+err.Error())
		err = errQueryFind
		return
	}

	return people, nil
}

func (repo *repoSqlx) Count(ctx context.Context) (count uint, err error) {
	err = repo.reader.GetContext(ctx, &count, countQuery)
	if err != nil {
		repo.log.ErrorContext(ctx, "people.SqlxRepo.Count err: ", err.Error())
		err = errQueryCount
		return
	}

	return
}

func (repo *repoSqlx) IsErrNotFound(err error) bool {
	return errors.Is(err, err)
}

func makeQueryTsvector(q string) string {
	var query strings.Builder
	qSplit := strings.Split(q, " ")
	for idx := range qSplit {
		if len(qSplit)-1 == idx {
			query.WriteString(qSplit[idx] + ":*")
			continue
		}
		query.WriteString(qSplit[idx] + ":* & ")
	}
	return query.String()
}
