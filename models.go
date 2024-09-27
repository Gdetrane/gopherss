package main

import (
	"time"

	"github.com/Gdetrane/gopherss/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollow struct {
  ID uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  UserID uuid.UUID `json:"user_id"`
  FeedID uuid.UUID `json:"feed_id"`
}

func remapDatabaseFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
  return FeedFollow{
    ID: dbFeedFollow.UserID,
    CreatedAt: dbFeedFollow.CreatedAt,
    UpdatedAt: dbFeedFollow.UpdatedAt,
    UserID: dbFeedFollow.UserID,
    FeedID: dbFeedFollow.FeedID,
  }
}

func remapAllFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
  remappedFeedFollows := make([]FeedFollow, 0)

  for _, dbFeedFollow := range dbFeedFollows {
    remappedFeedFollow := remapDatabaseFeedFollow(dbFeedFollow)
    remappedFeedFollows = append(remappedFeedFollows, remappedFeedFollow)
  }

  return remappedFeedFollows
}

func remapDatabaseUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

func remapDatabaseFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}

func remapAllFeeds(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, 0)

	for _, feed := range dbFeeds {
		remappedFeed := remapDatabaseFeed(feed)
		feeds = append(feeds, remappedFeed)
	}

	return feeds
}
