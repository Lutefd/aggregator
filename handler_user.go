package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Lutefd/aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing request")
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (cfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	query := r.URL.Query()
	limit := query.Get("limit")
	conv, err := strconv.Atoi(limit)
	parsedLimit := int32(conv)
	if err != nil {
		parsedLimit = int32(10)
	}
	posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  parsedLimit,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting posts")
		return
	}
	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
