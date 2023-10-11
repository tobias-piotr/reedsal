package auth

import (
	"net/http"
	"reedsal/api"
	"reedsal/users"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *users.UserRepository
}

func NewAuthService(repo *users.UserRepository) *AuthService {
	return &AuthService{repo}
}

func (s AuthService) Register(data users.UserCreatePayload) (*users.User, error) {
	exists, err := s.Repo.GetUserExistence("email", data.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, api.NewAPIError(http.StatusBadRequest, "User already exists", nil)
	}

	password, err := data.GetHashedPassword()
	if err != nil {
		return nil, err
	}
	data.Password = string(password)

	return s.Repo.CreateUser(data)
}

func (s AuthService) Login(data users.UserLoginPayload) (string, error) {
	user, err := s.Repo.GetUserWithPassword(data.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return "", api.NewAPIError(http.StatusBadRequest, "Invalid password", nil)
	}

	return NewToken(user.ID, user.Email)
}
