package people

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickchagastavares/rinha-backend/internal/entities"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type repoPgxpool struct {
	log logger.Logger
	db  *pgxpool.Pool
}

func NewPgxPool(log logger.Logger, db *pgxpool.Pool) IRepository {
	return &repoPgxpool{log: log, db: db}
}

func (repo *repoPgxpool) Create(ctx context.Context, person entities.PersonRequest) (err error) {
	_, err = repo.db.Exec(
		ctx, createPersonQuery,
		person.ID, person.Name, person.NickName,
		person.BirthDate, person.StackDB, person.SearchStr(), person.CreatedAt)
	if err != nil {
		if repo.validDuplicate(err) {
			err = errNickNameQueryCreate
			return
		}
		repo.log.ErrorContext(ctx, "people.repoPgxpool.Create data: ", person, ", err: "+err.Error())
		err = errQueryCreate
		return
	}
	return
}

func (repo *repoPgxpool) FindByID(ctx context.Context, id string) (person entities.Person, err error) {
	err = repo.db.
		QueryRow(ctx, findByIdQuery, id).
		Scan(&person.ID, &person.Name, &person.NickName, &person.BirthDate, &person.StackDB)
	if err != nil {
		repo.log.ErrorContext(ctx, "people.repoPgxpool.FindByID id: "+id+", err: "+err.Error())
		if err.Error() == sql.ErrNoRows.Error() {
			err = errNotFound
			return
		}
		err = errQueryFind
		return
	}

	person.ConvertStackDB()
	return
}

func (repo *repoPgxpool) FindBySearch(ctx context.Context, query string) (people []entities.Person, err error) {
	people = make([]entities.Person, 0, 50)
	query = makeQueryTsvector(query)
	rows, err := repo.db.Query(ctx, findBySearchQuery, query)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		repo.log.ErrorContext(ctx, "people.repoPgxpool.FindBySearch query: "+query+", err: "+err.Error())
		err = errQueryFind
		return
	}

	return scanRows(ctx, rows, repo.scanPerson)
}

func (repo *repoPgxpool) Count(ctx context.Context) (count uint, err error) {
	err = repo.db.
		QueryRow(ctx, countQuery).
		Scan(&count)
	if err != nil {
		repo.log.ErrorContext(ctx, "people.repoPgxpool.Count err: ", err.Error())
		err = errQueryCount
		return
	}
	return
}

func (repo *repoPgxpool) IsErrNotFound(err error) bool {
	return errors.Is(err, errNotFound)
}

func (repo *repoPgxpool) IsErrDuplicate(err error) bool {
	return errors.Is(err, errNickNameQueryCreate)
}

func (repo *repoPgxpool) validDuplicate(err error) bool {
	return err.Error() == "ERROR: duplicate key value violates unique constraint \"people_nick_name_key\" (SQLSTATE 23505)"
}

func (repo *repoPgxpool) scanPerson(ctx context.Context, row pgx.Row) (entities.Person, error) {
	person := entities.Person{}

	err := row.Scan(&person.ID, &person.Name, &person.NickName, &person.BirthDate, &person.StackDB)
	if err != nil {
		return person, fmt.Errorf("could not scan cid: %w", err)
	}

	person.ConvertStackDB()
	return person, nil
}

func scanRows[T any](ctx context.Context, rows pgx.Rows, scanRow func(ctx context.Context, row pgx.Row) (T, error)) ([]T, error) {
	defer rows.Close()
	destRows := make([]T, 0)
	for rows.Next() {
		destRow, err := scanRow(ctx, rows)
		if err != nil {
			return nil, err
		}
		destRows = append(destRows, destRow)
	}
	return destRows, nil
}
