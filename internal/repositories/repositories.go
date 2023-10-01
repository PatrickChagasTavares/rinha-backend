package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickchagastavares/rinha-backend/internal/repositories/database/people"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

type (
	// Container model to export instance repositories
	Container struct {
		Database SqlContainer
	}

	SqlContainer struct {
		People people.IRepository
	}

	// Options struct of options to create a new repositories
	Options struct {
		WriterSqlx *sqlx.DB
		ReaderSqlx *sqlx.DB
		Log        logger.Logger
	}
)

// New Create a new instance of repositories
func New(opts Options) *Container {
	return &Container{
		Database: SqlContainer{
			People: people.NewSqlx(opts.Log, opts.WriterSqlx, opts.ReaderSqlx),
		},
	}
}
