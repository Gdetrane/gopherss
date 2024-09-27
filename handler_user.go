package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Gdetrane/gopherss/internal/auth"
	"github.com/Gdetrane/gopherss/internal/database"
	"github.com/google/uuid"
)

func (cfg apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request)  {
  type parameters struct {
    Name string `json:"name"`
  }
  params := parameters{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s\n", err))
    return
  }

  user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
    ID: uuid.New(),
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    Name: params.Name,
  })
  if err != nil {
    respondWithError(w, 201, fmt.Sprintf("Could not create user: %s\n", err))
  }

  respondWithJSON(w, 200, remapDatabaseUser(user))
}

func (cfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User)  {
  apiKey, err := auth.GetAPIKey(r.Header)
  if err != nil {
    respondWithError(w, 403, fmt.Sprintf("Auth error: %s", err))
    return
  }
  
  usr, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
  if err != nil {
    respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
    return
  }

  respondWithJSON(w, 200, remapDatabaseUser(usr))

}
