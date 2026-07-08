package postgres

import (
	"context"
	"fmt"

	"github.com/abozorov/cinema/user_service/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	*pgxpool.Pool
}

func NewConn(cfg config.Config) (*Postgres, error) {

	pool, err := pgxpool.New(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.Name,
	))

	return &Postgres{
		Pool: pool,
	}, err
}
