package repo

import "github.com/abozorov/cinema/pkg/postgres"

type Repo struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		pg: pg,
	}
}
