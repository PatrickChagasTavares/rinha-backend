package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type (
	Person struct {
		ID        string         `db:"id" json:"id"`
		Name      string         `db:"name" json:"name"`
		NickName  string         `db:"nick_name" json:"apelido"`
		BirthDate string         `db:"birth_date" json:"nascimento"`
		Stack     pq.StringArray `db:"stack" json:"stack"`
		CreatedAt time.Time      `db:"created_at" json:"created_at"`
	}

	PersonRequest struct {
		ID        string         `json:"-"`
		Name      string         `json:"name" validate:"required"`
		NickName  string         `json:"apelido" validate:"required"`
		BirthDate string         `json:"nascimento" validate:"required,dateformat"`
		Stack     pq.StringArray `json:"stack"`
		CreatedAt time.Time      `json:"-"`
	}
)

func (pr *PersonRequest) PreSave() {
	pr.ID = uuid.NewString()
	pr.CreatedAt = time.Now()
}

func (pr *PersonRequest) ToPerson() Person {
	return Person{
		ID:        pr.ID,
		Name:      pr.Name,
		NickName:  pr.NickName,
		BirthDate: pr.BirthDate,
		Stack:     pr.Stack,
		CreatedAt: pr.CreatedAt,
	}
}
