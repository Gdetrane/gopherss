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

func (cfg apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User)  {
  type parameters struct {
    Name string `json:"name"`
    Url string `json:"url"`
  }
  decoder := json.NewDecoder(r.Body)

  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
    return
  }

  feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
    ID: uuid.New(),
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    Name: params.Name,
    Url: params.Url,
    UserID: user.ID,
  })
  if err != nil {
    respondWithError(w, 400, fmt.Sprintf("Error creating feed: %v", err))
    return
  }
  
  respondWithJSON(w, 200, remapDatabaseFeed(feed))

}

func (cfg apiConfig) handlerGetFeedsByUserId(w http.ResponseWriter, r *http.Request)  {
  apiKey, err := auth.GetAPIKey(r.Header)
  if err != nil {
    respondWithError(w, 403, fmt.Sprintf("Auth error: %s", err))
    return
  }

  usr, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
  if err != nil {
    respondWithError(w, 400, fmt.Sprintf("Error retrieving user: %s", err))
  }

  feeds, err := cfg.DB.GetFeedsByUserId(r.Context(), usr.ID)
  if err != nil {
    respondWithError(w, 400, fmt.Sprintf("Error retrieving feeds: %s", err))
    return
  }

  respondWithJSON(w, 200, feeds)
}

func (cfg apiConfig) handlerGetAllFeeds(w http.ResponseWriter, r *http.Request)  {
  feeds, err :=cfg.DB.GetAllFeeds(r.Context())
  if err != nil {
    respondWithError(w, 400, fmt.Sprintf("Error retrieving feeds: %v", err))
  }

  respondWithJSON(w, 200, remapAllFeeds(feeds))
}
