package repo

import (
	"context"
	"fmt"

	"github.com/abozorov/cinema/cmd/user/internal/models"
	"github.com/abozorov/cinema/pkg/postgres"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repo struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		pg: pg,
	}
}

func execAnalysis(res pgconn.CommandTag, err error) error {
	if err != nil {
		return fmt.Errorf("user_repo.execAnalysis: %w", err)
	}
	if rows := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("user_repo.execAnalysis: %w", models.ErrUserNotFound)
	}
	return nil
}

func (r *Repo) Add(ctx context.Context, user models.User) error {
	const query = `INSERT INTO users (name, email, phone, password_hash, age) 
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err := r.pg.Exec(ctx, query,
		user.Name,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.Age,
	)
	if err != nil {
		return fmt.Errorf("user_repo.Add: %w", postgresToErrs(err))
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id int) (*models.User, error) {
	const query = `
		SELECT  id,
			name,
			email,
			phone,
			password_hash,
			age
		 FROM users
		 WHERE id = $1;
	`
	row := r.pg.QueryRow(ctx, query, id)

	var (
		user  models.User
		phone pgtype.Text
	)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&phone,
		&user.PasswordHash,
		&user.Age,
	)
	if err != nil {
		return &models.User{},
			fmt.Errorf("user_repo.GetByID: %w", postgresToErrs(err))
	}
	user.Phone = phone.String
	return &user, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `
		SELECT  id,
			name,
			email,
			phone,
			password_hash,
			age
		 FROM users
		 WHERE email = $1;
	`
	row := r.pg.QueryRow(ctx, query, email)

	var (
		user  models.User
		phone pgtype.Text
	)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&phone,
		&user.PasswordHash,
		&user.Age,
	)
	if err != nil {
		return &models.User{},
			fmt.Errorf("user_repo.GetByID: %w", postgresToErrs(err))
	}
	user.Phone = phone.String
	return &user, nil
}

func (r *Repo) Update(ctx context.Context, user models.User) error {
	const query = `
		UPDATE users
		SET name = $2,
		phone = $3
		WHERE id = $1;
	`

	err := execAnalysis(r.pg.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Phone,
	))

	if err != nil {
		return fmt.Errorf("user_repo.Update: %w", postgresToErrs(err))
	}
	return nil
}
