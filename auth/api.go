package auth

import (
	"net/http"
	"reedsal/api"
	"reedsal/users"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	Srv *AuthService
}

func NewAuthHandler(db *sqlx.DB) AuthHandler {
	return AuthHandler{NewAuthService(users.NewUserRepository(db))}
}

func (h AuthHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	return r
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload users.UserCreatePayload

	err := api.ProcessPayload(w, r, &payload)
	if err != nil {
		return
	}

	err = api.ValidatePayload(w, payload)
	if err != nil {
		return
	}

	user, err := h.Srv.Register(payload)
	if err != nil {
		api.RespondWithError(w, err)
		return
	}

	api.Respond(w, http.StatusCreated, user)
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload users.UserLoginPayload

	err := api.ProcessPayload(w, r, &payload)
	if err != nil {
		return
	}

	err = api.ValidatePayload(w, payload)
	if err != nil {
		return
	}

	token, err := h.Srv.Login(payload)
	if err != nil {
		api.RespondWithError(w, err)
		return
	}

	api.Respond(w, http.StatusOK, TokenResponse{token})
}
