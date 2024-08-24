package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/srijan-raghavula/feeder/internal/database"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("Welcome to feeder"))
}

func readiness(w http.ResponseWriter, r *http.Request) {
	type status struct {
		Status string `json:"status"`
	}
	statusOK := status{
		Status: "ok",
	}
	writeData, err := json.Marshal(statusOK)
	if err != nil {
		resWithJSON(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(writeData)
}

func errHandler(w http.ResponseWriter, r *http.Request) {
	type errJSON struct {
		Err string `json:"error"`
	}
	errBody := errJSON{
		Err: "Internal Server Error",
	}
	writeData, err := json.Marshal(errBody)
	if err != nil {
		resWithError(w, 500, err.Error())
	}
	w.WriteHeader(500)
	w.Write(writeData)
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Name string `json:"name"`
	}
	reqBody := body{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	newUUID := uuid.New()
	userInfo := database.CreateUserParams{
		ID:        newUUID.String(),
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
		Name:      reqBody.Name,
	}
	userFromDB, err := cfg.DB.CreateUser(r.Context(), userInfo)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	userToEncode := &user{
		ID:        userFromDB.ID,
		CreatedAt: userFromDB.CreatedAt,
		UpdatedAt: userFromDB.UpdatedAt,
		Name:      userFromDB.Name,
		ApiKey:    userFromDB.ApiKey,
	}
	writeData, err := json.Marshal(userToEncode)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(writeData)
}

func getApiKeyFromReq(r *http.Request) string {
	header := r.Header.Get("Authorization")
	apiKey := strings.TrimPrefix(header, "ApiKey ")
	return apiKey
}

func (cfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request, userFromDB database.User) {
	userToEncode := user{
		ID:        userFromDB.ID,
		CreatedAt: userFromDB.CreatedAt,
		UpdatedAt: userFromDB.UpdatedAt,
		Name:      userFromDB.Name,
		ApiKey:    userFromDB.ApiKey,
	}
	writeData, err := json.Marshal(userToEncode)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(writeData)
}

func (cfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, userFromDB database.User) {
	type body struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	reqBody := body{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	newUUID := uuid.New()
	feedParams := database.CreateFeedParams{
		ID:        newUUID.String(),
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
		Name:      reqBody.Name,
		Url:       reqBody.Url,
		UserID:    userFromDB.ID,
	}

	feedFromDB, err := cfg.DB.CreateFeed(r.Context(), feedParams)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	feedToEncode := feed{
		ID:        feedFromDB.ID,
		CreatedAt: feedFromDB.CreatedAt,
		UpdatedAt: feedFromDB.UpdatedAt,
		Name:      feedFromDB.Name,
		Url:       feedFromDB.Url,
		UserID:    feedFromDB.UserID,
	}
	writeData, err := json.Marshal(feedToEncode)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(writeData)
}

func (apiCfg *apiConfig) getAllFeeds(w http.ResponseWriter, r *http.Request) {
	feedsFromDB, err := apiCfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	feedsToEncode := make([]feed, 0)
	for _, feedFromDB := range feedsFromDB {
		feedToEncode := feed{
			ID:        feedFromDB.ID,
			CreatedAt: feedFromDB.CreatedAt,
			UpdatedAt: feedFromDB.UpdatedAt,
			Name:      feedFromDB.Name,
			Url:       feedFromDB.Url,
			UserID:    feedFromDB.UserID,
		}
		feedsToEncode = append(feedsToEncode, feedToEncode)
	}
	writeData, err := json.Marshal(feedsToEncode)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(writeData)
}
