package people

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/patrickchagastavares/rinha-backend/internal/entities"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type repoSqlx struct {
	log logger.Logger
	db  *sqlx.DB
}

func NewSqlx(log logger.Logger, db *sqlx.DB) IRepository {
	return &repoSqlx{log: log, db: db}
}

func (repo *repoSqlx) Create(ctx context.Context, person entities.PersonRequest) (err error) {
	_, err = repo.db.ExecContext(
		ctx, createPersonQuery,
		person.ID, person.Name, person.NickName,
		person.BirthDate, person.StackDB, person.SearchStr(), person.CreatedAt)
	if err != nil {
		if repo.validDuplicate(err) {
			err = errNickNameQueryCreate
			return
		}
		repo.log.ErrorContext(ctx, "people.SqlxRepo.Create data: ", person, ", err: "+err.Error())
		err = errQueryCreate
		return
	}

	return nil
}

func (repo *repoSqlx) FindByID(ctx context.Context, id string) (person entities.Person, err error) {
	err = repo.db.GetContext(ctx, &person, findByIdQuery, id)
	if err != nil {
		repo.log.ErrorContext(ctx, "people.SqlxRepo.FindByID id: "+id+", err: "+err.Error())
		if err == sql.ErrNoRows {
			err = errNotFound
			return
		}
		err = errQueryFind
		return
	}

	person.ConvertStackDB()
	return
}

func (repo *repoSqlx) FindBySearch(ctx context.Context, query string) (people []entities.Person, err error) {
	people = make([]entities.Person, 0, 50)
	query = makeQueryTsvector(query)
	err = repo.db.SelectContext(ctx, &people, findBySearchQuery, query)
	if err != nil && err != sql.ErrNoRows {
		repo.log.ErrorContext(ctx, "people.SqlxRepo.FindBySearch query: "+query+", err: "+err.Error())
		err = errQueryFind
		return
	}

	for idx := range people {
		people[idx].ConvertStackDB()
	}

	return people, nil
}

func (repo *repoSqlx) Count(ctx context.Context) (count uint, err error) {
	err = repo.db.GetContext(ctx, &count, countQuery)
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

func (repo *repoSqlx) IsErrDuplicate(err error) bool {
	return errors.Is(err, errNickNameQueryCreate)
}

func (repo *repoSqlx) validDuplicate(err error) bool {
	return err.Error() == "pq: duplicate key value violates unique constraint \"people_nick_name_key\""
}
