package models

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

type Movie struct {
	ID          int
	Title       string
	Description string
	Duration    int
	AgeLimit    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var TimeFormat = time.RFC1123Z

var (
	ErrEmptyTitle         = errors.New("empty title")
	ErrEmptyDescrittion   = errors.New("empty description")
	ErrInvalidMovieId     = errors.New("invalid movieId")
	ErrInvalidTitle       = errors.New("invalid title")
	ErrInvalidDescrittion = errors.New("invalid description")
	ErrInvalidDuration    = errors.New("invalid duration")
	ErrInvalidAgeLimit    = errors.New("invalid age limit")
)

func NewMovie(title, description string, duration, ageLimit int) (*Movie, error) {
	err := validateTitle(&title)
	if err != nil {
		return &Movie{}, err
	}

	err = validateDescription(&description)
	if err != nil {
		return &Movie{}, err
	}

	err = validateDuration(duration)
	if err != nil {
		return &Movie{}, err
	}

	err = validateAgeLimit(ageLimit)
	if err != nil {
		return &Movie{}, err
	}

	return &Movie{
		Title:       title,
		Description: description,
		Duration:    duration,
		AgeLimit:    ageLimit,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil

}

// validateTitle проверяет название фильма.
// Если название не указано (nil) или пустое – возвращает ErrEmptyTitle.
// Если название слишком длинное (>255 символов) – возвращает ErrInvalidTitle.
// В остальных случаях – nil.
func validateTitle(t *string) error {
	if t == nil || strings.TrimSpace(*t) == "" {
		return ErrEmptyTitle
	}
	if utf8.RuneCountInString(*t) > 255 {
		return ErrInvalidTitle
	}
	return nil
}

// validateDescription проверяет описание фильма.
// Если описание не указано (nil) или пустое – возвращает ErrEmptyDescrittion.
// Если описание слишком длинное (>1000 символов) – возвращает ErrInvalidDescrittion.
// В остальных случаях – nil.
func validateDescription(d *string) error {
	if d == nil || strings.TrimSpace(*d) == "" {
		return ErrEmptyDescrittion
	}
	if utf8.RuneCountInString(*d) > 1000 {
		return ErrInvalidDescrittion
	}
	return nil
}

// validateDuration проверяет продолжительность фильма в минутах.
// Возвращает ErrInvalidDuration, если значение <= 0 или > 600.
func validateDuration(d int) error {
	if d <= 0 || d > 600 {
		return ErrInvalidDuration
	}
	return nil
}

// validateAgeLimit проверяет возрастной рейтинг.
// Возвращает ErrInvalidAgeLimit, если значение < 0 или > 21.
func validateAgeLimit(a int) error {
	if a < 0 || a > 21 {
		return ErrInvalidAgeLimit
	}
	return nil
}

func (m *Movie) IsValid() error {
	err := validateTitle(&m.Title)
	if err != nil {
		return err
	}

	err = validateDescription(&m.Description)
	if err != nil {
		return err
	}

	err = validateDuration(m.Duration)
	if err != nil {
		return err
	}

	err = validateAgeLimit(m.AgeLimit)
	if err != nil {
		return err
	}

	return nil
}

func (m *Movie) Update(title, description string, duration, ageLimit int) error {
	err := validateTitle(&title)
	if err != nil {
		return err
	}

	err = validateDescription(&description)
	if err != nil {
		return err
	}

	err = validateDuration(duration)
	if err != nil {
		return err
	}

	err = validateAgeLimit(ageLimit)
	if err != nil {
		return err
	}

	m.Title = title
	m.Description = description
	m.Duration = duration
	m.AgeLimit = ageLimit

	return nil
}
