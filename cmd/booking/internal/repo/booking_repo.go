package repo

import (
	"context"
	"fmt"

	"github.com/abozorov/cinema/cmd/booking/internal/models"
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

func (r *Repo) Create(ctx context.Context, b *models.Booking) (int, error) {
	query := `
		INSERT INTO bookings (user_id, movie_id, status)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	tx, err := r.pg.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("booking_repo.Create: %w", postgresToErrs(err))
	}
	defer tx.Commit(ctx)

	row := tx.QueryRow(ctx, query,
		b.UserID,
		b.MovieID,
		b.Status,
	)
	var id int
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("booking_repo.Create: %w", postgresToErrs(err))
	}
	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int) (*models.Booking, error) {
	query := `
		SELECT id, user_id, movie_id, status, created_at, updated_at
		FROM bookings
		WHERE id = $1;
	`

	b := &models.Booking{}
	err := r.pg.QueryRow(ctx, query, id).Scan(
		&b.ID,
		&b.UserID,
		&b.MovieID,
		&b.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("booking_repo.Get: %w", postgresToErrs(err))
	}
	return b, nil
}

func (r *Repo) GetUserBookings(ctx context.Context, userID int) ([]*models.Booking, error) {
	query := `
		SELECT id, user_id, movie_id, status, created_at, updated_at
		FROM bookings
		WHERE id = $1`

	rows, err := r.pg.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("booking_repo.GetUserBookings: %w", postgresToErrs(err))
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		b := &models.Booking{}
		err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.MovieID,
			&b.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("booking_repo.GetUserBookings: %w", postgresToErrs(err))
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *Repo) Cancel(ctx context.Context, id int) error {
	query := `
		UPDATE bookings
		SET status = $2
		WHERE id = $1;
	`

	err := execAnalysis(r.pg.Exec(ctx, query,
		id,
		models.StatusCanceled,
	))
	if err != nil {
		return fmt.Errorf("booking_repo.Cancel: %w", postgresToErrs(err))
	}
	return nil
}
