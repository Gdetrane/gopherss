package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Gdetrane/gopherss/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (cfg apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating follow for feed: %v", err))
		return
	}

	respondWithJSON(w, 200, remapDatabaseFeedFollow(feedFollow))
}

func (cfg apiConfig) handlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollow, err := cfg.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error finding feed follows: %v", err))
		return
	}

	respondWithJSON(w, 200, remapAllFeedFollows(feedFollow))
}

func (cfg apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	if feedFollowIDStr == "" {
		respondWithError(w, 500, "Error deleting feed, no feed follow found")
	}

	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not parse feed follow id %s: %v", feedFollowIDStr, err))
		return
	}

  deletionErr := cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
    ID: feedFollowID,
    UserID: user.ID,
  })
  if deletionErr != nil {
    respondWithError(w, 400, fmt.Sprintf("Error deleting feed follow: %v", deletionErr))
    return
  }

  deletionSuccess := map[string]string{
    "msg": "item deleted successfully",
  }

  respondWithJSON(w, 200, deletionSuccess)
}
