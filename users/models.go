package users

import (
	"net/mail"
	"reedsal/api"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserCreatePayload struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (p UserCreatePayload) Validate() *api.ValidationError {
	details := api.ValidationDetails{}

	if _, err := mail.ParseAddress(p.Email); err != nil {
		details["email"] = err.Error()
	}
	if len(p.Password) < 8 {
		details["password"] = "Passwords is too short"
	} else if len(p.Password) > 32 {
		details["password"] = "Passwords is too long"
	} else if p.Password != p.PasswordConfirmation {
		details["password"] = "Passwords don't match"
	}

	if len(details) > 0 {
		return &api.ValidationError{Details: details}
	}
	return nil
}

func (p UserCreatePayload) GetHashedPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
}

type UserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p UserLoginPayload) Validate() *api.ValidationError {
	return nil
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type UserWithPassword struct {
	User
	Password string `json:"password"`
}
