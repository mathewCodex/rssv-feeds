package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mathewCodex/rssagg/auth"
)
/////////creating feed handler-------
func (apiCfg *apiConfig)handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct {
		FeedID uuid.UUID `just:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error Parsing json: %v", err))
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
		UserID: user.ID,//authenticated user
		FeedID: params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldnt create feed follow: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}






func (apiCfg *apiConfig)handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){
	// type parameters struct {
	// 	FeedID uuid.UUID `just:"feed_id"`
	// }
	// decoder := json.NewDecoder(r.Body)
	// params := parameters{}
	// err := decoder.Decode(&params)
	// if err != nil{
	// 	respondWithError(w, 400, fmt.Sprintf("Error Parsing json: %v", err))
	// 	return
	// }
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldnt get feed follow: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollows))
}

//delte func handler feed follow returning ""
func (apiCfg *apiConfig)handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollowIDstr := chi.URLParam(r, "feedFollowID");//to  pass the req and a key and  matches it between the {}
	feedFollowID, err := uuid.Parse(feedFollowIDstr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldnt Parser feed followID: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID: feedFollowID,
		UserID: user.ID,

})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldnt delete feed follow: %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}


