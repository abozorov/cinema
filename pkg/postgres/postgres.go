package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	*pgxpool.Pool
}

func NewConn(conn string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), conn)

	return &Postgres{
		Pool: pool,
	}, err
}
