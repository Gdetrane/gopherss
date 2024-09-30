package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Gdetrane/gopherss/internal/database"
	"github.com/google/uuid"
)

func (cfg apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request)  {
  type parameters struct {
    Title string `json:"title"`
    Url string `json:"url"`
    Description string `json:"description"`
  }

  params := parameters{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
    return
  }

  post, err := cfg.DB.CreatePost(r.Context(), database.CreatePostParams{
    ID: uuid.New(),
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    Title: params.Title,
    Url: params.Url,
  })
  if err != nil {
    respondWithError(w, 201, fmt.Sprintf("Could not create post: %s", err))
  }

  respondWithJSON(w, 200, remapDatabasePost(post))
}
