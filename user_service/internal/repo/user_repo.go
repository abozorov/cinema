package repo

import "github.com/abozorov/cinema/user_service/pkg/postgres"

type Repo struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		pg: pg,
	}
}
