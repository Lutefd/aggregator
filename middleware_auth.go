package main

import (
	"net/http"

	"github.com/Lutefd/aggregator/internal/auth"
	"github.com/Lutefd/aggregator/internal/database"
)

type authenticatedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) middlewareAuth(next authenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		next(w, r, user)
	}
}
