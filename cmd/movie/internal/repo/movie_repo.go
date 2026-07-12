package repo

import (
	"context"
	"fmt"

	"github.com/abozorov/cinema/cmd/movie/internal/models"
	"github.com/abozorov/cinema/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		return fmt.Errorf("movie_repo.execAnalysis: %w", err)
	}
	if rows := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("movie_repo.execAnalysis: %w", pgx.ErrNoRows)
	}
	return nil
}

func (r *Repo) Create(ctx context.Context, m *models.Movie) (int, error) {
	query := `
		INSERT INTO movies (title, description, duration, age_limit, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	tx, err := r.pg.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("movie_repo.Create: %w", postgresToErrs(err))
	}
	defer tx.Commit(ctx)

	row := tx.QueryRow(ctx, query,
		m.Title,
		m.Description,
		m.Duration,
		m.AgeLimit,
		m.CreatedAt, 
		m.UpdatedAt,
	)
	var id int
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("movie_repo.Create: %w", postgresToErrs(err))
	}
	return id, nil
}

func (r *Repo) GetByID(ctx context.Context, id int) (*models.Movie, error) {
	query := `
		SELECT id, 
			title, 
			description, 
			duration, 
			age_limit, 
			created_at, 
			updated_at
		FROM movies
		WHERE id = $1;
		`
	m := &models.Movie{}
	err := r.pg.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.Title,
		&m.Description,
		&m.Duration,
		&m.AgeLimit,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("movie_repo.GetByID: %w", postgresToErrs(err))
	}
	return m, nil
}

func (r *Repo) List(ctx context.Context) ([]models.Movie, error) {
	query := `
		SELECT id, 
			title, 
			description, 
			duration, 
			age_limit, 
			created_at, 
			updated_at
		FROM movies
		ORDER BY id;
		`
	rows, err := r.pg.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("movie_repo.List: %w", postgresToErrs(err))
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.Duration,
			&m.AgeLimit,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("movie_repo.List: %w", postgresToErrs(err))
		}
		movies = append(movies, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("movie_repo.List: %w", postgresToErrs(err))
	}
	return movies, nil
}

func (r *Repo) Update(ctx context.Context, m *models.Movie) error {
	query := `
		UPDATE movies
		SET title = $2, 
		description = $3, 
		duration = $4, 
		age_limit = $5, 
		updated_at = NOW()
		WHERE id = $1;
		`
	err := execAnalysis(r.pg.Exec(ctx, query,
		m.ID,
		m.Title,
		m.Description,
		m.Duration,
		m.AgeLimit,
	))
	if err != nil {
		return fmt.Errorf("movie_repo.Update: %w", postgresToErrs(err))
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, id int) error {
	return nil
}
