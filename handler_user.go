package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mathewCodex/rssagg/auth"
)
/** 
* create user handler function
considered as method apicfg is a pointer to the apiConfig

*/
func (apiCfg *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)//decoding into a pointer to params
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error Parsing json: %v", err))
		return
	}
	//create user resturns a user and an err
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		//this is a struct that should have all the required field...
		ID: uuid.New(),
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
		Name: params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldnt create user: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseUserToUser(user))
}
/** 
handler to get user
*/
func (apiCfg *apiConfig)handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User){
	// apiKey, err := auth.GetAPIKey(r.Header)
	// if err != nil{
	// 	respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
	// 	return 
	// }

	// user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)//the Context gives you way to track multiple content in the go routine
	// if err != nil {
	// 	respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
	// 	return 
	// }
	respondWithJSON(w, 200, databaseUserToUser(user))
}
func (apiCfg *apiConfig)handlerGetPostForUser(w http.ResponseWriter, r *http.Request, user database.User){
	post, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID:	 user.ID,
		Limit:	10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldnt get post: %v", err))
		return
	}
	respondWithJSON(w, 200, databasePostsToPosts(post))
}