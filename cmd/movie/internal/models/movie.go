package models

import "time"

type Movie struct {
	ID          int
	Title       string
	Description string
	Dutation    int
	AgeLimit    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MovieService interface {
}

type MovieRepository interface {
}

// CREATE TABLE movies
// (
//     id          SERIAL PRIMARY KEY,
//     title       VARCHAR(255) NOT NULL,
//     description TEXT         NOT NULL,
//     duration    INT          NOT NULL,
//     age_limit   INT          NOT NULL,
//     created_at  TIMESTAMP    NOT NULL,
//     updated_at  TIMESTAMP    NOT NULL
// );
