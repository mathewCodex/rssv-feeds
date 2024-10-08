package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mathewCodex/rssagg/auth"
)

func (apiCfg *apiConfig)handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct {
		Name string `json:"name"`
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error Parsing json: %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
		Name: params.Name,
		url: params.URL,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldnt create feed: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}



//handler to get all feeds 
func (apiCfg *apiConfig)handlerGetFeeds(w http.ResponseWriter, r *http.Request){

	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldnt get feeds: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedToFeed(feeds))
}

