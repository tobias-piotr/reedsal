package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Responds with an api error, defaults to 'internal server error' if provided error
// is not an api error instance.
func RespondWithError(w http.ResponseWriter, err error) error {
	aerr, ok := err.(APIError)
	if !ok {
		slog.Error("Internal server error", "err", err)
		aerr = NewAPIError(http.StatusInternalServerError, "Internal server error", nil)
		// defer panic(err)
	}
	w.WriteHeader(aerr.StatusCode)
	res, err := json.Marshal(aerr)
	if err != nil {
		return err
	}
	w.Write(res)
	return nil
}

// Responds with marshalled data.
func Respond(w http.ResponseWriter, status int, data any) error {
	w.WriteHeader(status)
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Write(res)
	return nil
}
