package users

import (
	"database/sql"
	"fmt"
	"net/http"
	"reedsal/api"

	"github.com/jmoiron/sqlx"
)

var ErrNoUser = api.NewAPIError(http.StatusBadRequest, "User does not exist", nil)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r UserRepository) CreateUser(data UserCreatePayload) (*User, error) {
	query := `
		INSERT INTO users (email, password)
		VALUES (:email, :password)
		RETURNING id, email;
	`
	query, args, err := r.DB.BindNamed(query, data)
	if err != nil {
		return nil, err
	}

	var user User
	err = r.DB.QueryRowx(query, args...).StructScan(&user)
	return &user, err
}

func (r UserRepository) GetUserExistence(identifier string, value string) (bool, error) {
	query := `
		SELECT 1
		FROM users
		WHERE %s=$1
	`
	query = fmt.Sprintf(query, identifier)

	var exists bool
	err := r.DB.Get(&exists, query, value)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return exists, err
}

func (r UserRepository) GetUserWithPassword(email string) (*UserWithPassword, error) {
	query := `
		SELECT id, email, password
		FROM users
		WHERE email=$1
	`
	var user UserWithPassword
	err := r.DB.Get(&user, query, email)
	if err == sql.ErrNoRows {
		return nil, ErrNoUser
	}
	return &user, err
}
