package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Gdetrane/gopherss/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
  DB *database.Queries
}
 
func main() {
	fmt.Println("Hello world!")

  /*DEBUG!*/
  // feed, err := urlToFeed("https://rss.draghetti.it/ilpost.xml")
  // if err != nil {
  //   log.Fatal(err)
  // }
  // fmt.Println(feed)

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("No PORT environment variable found in env.")
	}

  dbUrl := os.Getenv("DB_URL")
  if dbUrl == "" {
    log.Fatal("No DB_URL found in env.")
  }

  conn, err := sql.Open("postgres", dbUrl)
  if err != nil {
    log.Fatal("Can't connect to database:", err)
  }

  queries := database.New(conn)

  apiCfg := apiConfig{
    DB: queries,
  }

  go startScraping(queries, 20, time.Minute)

	router := chi.NewRouter()

	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}),
	)

  v1Router := chi.NewRouter()
  v1Router.Get("/healthz", handlerReadiness)

  v1Router.Get("/err", handlerErr)

  v1Router.Post("/users", apiCfg.handlerCreateUser)
  v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

  v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
  v1Router.Get("/feeds", apiCfg.handlerGetAllFeeds)

  v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
  v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))
  v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

  v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))
  
  router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %s", portString)

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v", err)
		os.Exit(1)
	}
}
