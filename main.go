package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/go-chi/cors"
	//_ "github.com/lib/pq"
	// "github.com/google/uuid"
)


type apiConfig struct{
	DB *database.Queries
}
func main() {
    // feed, err := urlToFeed("https://wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)
	// fmt.Println("Hello World")
   godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found in the env")
	}
   

	// import db connection
	
dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DBURL is not found in the env")
	}
	
	/** 
	// connecting to db with the driver as first arg 
	and return a connection string or an error if any
	*/
	conn, err := sql.Open("postgres", dbURL)
	if err != nil{
		log.Fatal("Can't connect to db", err)
	}
//converting conn to db query
db := database.New(conn)
// if err != nil {
// 	log.Fatal("cant create a db", err)
// }

///creating an api config that takes a db as one of its field
apiCfg := apiConfig{
	DB: db,
}
go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()
//creating a cors handler
   router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowedCredentials: false,
		MaxAge: 300,
	}))

//using the chi router to hookup our path
	v1Router := chi.NewRouter()
      
	v1Router.Get("/healthz", handlerReadiness)
    v1Router.Get("/err", handlerErr)
    v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostForUser))	
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))


//nesting the path e.g /v1/endpoint
	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}
	log.Println("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Port:", portString)

};
