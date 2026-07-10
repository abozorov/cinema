package models

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// User представляет доменную модель пользователя
type User struct {
	ID           int
	Name         string
	Email        string
	Phone        string
	PasswordHash string
	Age          int
}

// UserService определяет бизнес-логику для работы с пользователями
type UserService interface {
	Add(ctx context.Context, user *User) (int, error)
	GetByID(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, user *User) error
}

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Add(ctx context.Context, user *User) (int, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
}

// Validation errors
var (
	ErrInvalidName   = errors.New("name must be between 2 and 100 characters")
	ErrInvalidPhone  = errors.New("phone must be 12 characters")
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrInvalidAge    = errors.New("invalid age")
	ErrInvalidUSerId = errors.New("invalid user id")
	ErrUserExists    = errors.New("user with this email already exists")
	ErrEmptyName     = errors.New("name cannot be empty")
	ErrEmptyPhone    = errors.New("phone cannot be empty")
	ErrEmptyEmail    = errors.New("email cannot be empty")
	ErrEmptyUserID   = errors.New("user id cannot be empty")
)

// emailRegex для валидации email
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewUser создает нового пользователя с валидацией
func NewUser(name, email, phone, password_hash string, age int) (*User, error) {
	if err := validateName(&name); err != nil {
		return nil, err
	}

	if err := validateEmail(&email); err != nil {
		return nil, err
	}

	if err := validatePhone(&phone); err != nil {
		return nil, err
	}

	if err := validateAge(&age); err != nil {
		return nil, err
	}

	return &User{
		Name:         name,
		Email:        email,
		Phone:        phone,
		PasswordHash: password_hash,
		Age:          age,
	}, nil
}

// validateAge проверяет корректность возраста
func validateAge(age *int) error {
	if *age < 0 || *age > 100 {
		return ErrInvalidAge
	}

	return nil
}

// validateName проверяет корректность номера телефона
func validatePhone(phone *string) error {

	*phone = strings.TrimSpace(*phone)
	if *phone == "" {
		return ErrEmptyPhone
	}

	length := utf8.RuneCountInString(*phone)
	if length != 12 {
		return ErrInvalidPhone
	}

	return nil
}

// validateName проверяет корректность имени
func validateName(name *string) error {

	*name = strings.TrimSpace(*name)
	if *name == "" {
		return ErrEmptyName
	}

	length := utf8.RuneCountInString(*name)
	if length < 2 || length > 100 {
		return ErrInvalidName
	}

	return nil
}

// validateEmail проверяет корректность email
func validateEmail(email *string) error {

	*email = strings.ToLower(strings.TrimSpace(*email))
	if *email == "" {
		return ErrEmptyEmail
	}

	// Проверяем длину
	if len(*email) < 5 || len(*email) > 254 {
		return ErrInvalidEmail
	}

	// Проверяем формат с помощью regex
	if !emailRegex.MatchString(*email) {
		return ErrInvalidEmail
	}

	// Дополнительные проверки
	parts := strings.Split(*email, "@")
	if len(parts) != 2 {
		return ErrInvalidEmail
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Проверяем локальную часть
	if len(localPart) == 0 || len(localPart) > 64 {
		return ErrInvalidEmail
	}

	// Проверяем доменную часть
	if len(domainPart) == 0 || len(domainPart) > 253 {
		return ErrInvalidEmail
	}

	// Проверяем, что домен содержит точку
	if !strings.Contains(domainPart, ".") {
		return ErrInvalidEmail
	}

	return nil
}

// Update обновляет данные пользователя
func (u *User) Update(name, email, phone string, age int) error {
	if err := validateName(&name); err != nil {
		return err
	}

	if err := validateEmail(&email); err != nil {
		return err
	}

	if err := validatePhone(&phone); err != nil {
		return err
	}

	if err := validateAge(&age); err != nil {
		return err
	}

	u.Name = name
	u.Email = email
	u.Phone = phone
	u.Age = age

	return nil
}

// IsValid проверяет, что пользователь содержит валидные данные
func (u *User) IsValid() error {

	if err := validateName(&u.Name); err != nil {
		return err
	}

	if err := validateEmail(&u.Email); err != nil {
		return err
	}

	if err := validatePhone(&u.Phone); err != nil {
		return err
	}

	if err := validateAge(&u.Age); err != nil {
		return err
	}

	return nil
}
