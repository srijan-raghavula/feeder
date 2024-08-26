package main

import (
	"encoding/json"
	"fmt"
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

func (apiCfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
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
	u, err := apiCfg.DB.CreateUser(r.Context(), userInfo)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	userToEncode := &user{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Name:      u.Name,
		ApiKey:    u.ApiKey,
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

func (apiCfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request, u database.User) {
	userToEncode := user{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Name:      u.Name,
		ApiKey:    u.ApiKey,
	}
	writeData, err := json.Marshal(userToEncode)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(writeData)
}

func (apiCfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	type body struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	type followNfeed struct {
		Feed       feed   `json:"feed"`
		FeedFollow follow `json:"feed_follow"`
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
		UserID:    u.ID,
	}

	feedFromDB, err := apiCfg.DB.CreateFeed(r.Context(), feedParams)
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
	followParams := database.FollowFeedParams{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
		FeedID:    feedToEncode.ID,
		UserID:    u.ID,
	}
	followFromDB, err := apiCfg.DB.FollowFeed(r.Context(), followParams)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	followToEncode := follow{
		ID:        followFromDB.ID,
		CreatedAt: followFromDB.CreatedAt,
		UpdatedAt: followFromDB.UpdatedAt,
		FeedID:    followFromDB.ID,
		UserID:    followParams.UserID,
	}
	resBody := followNfeed{
		Feed:       feedToEncode,
		FeedFollow: followToEncode,
	}
	writeData, err := json.Marshal(resBody)
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

func (apiCfg *apiConfig) followFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	type body struct {
		FeedID string `json:"feed_id"`
	}
	reqBody := body{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	followParams := database.FollowFeedParams{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
		FeedID:    reqBody.FeedID,
		UserID:    u.ID,
	}
	followFromDB, err := apiCfg.DB.FollowFeed(r.Context(), followParams)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	followToEncode := follow{
		ID:        followFromDB.ID,
		CreatedAt: followFromDB.CreatedAt,
		UpdatedAt: followFromDB.UpdatedAt,
		FeedID:    followFromDB.ID,
		UserID:    followParams.UserID,
	}
	writeData, err := json.Marshal(followToEncode)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(writeData)
}

func (apiCfg *apiConfig) unFollowFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	feedID := r.PathValue("feedFollowID")
	deleteParams := database.DeleteFollowParams{
		FeedID: feedID,
		UserID: u.ID,
	}
	err := apiCfg.DB.DeleteFollow(r.Context(), deleteParams)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Feed with ID %s deleted", feedID)))
}

func (apiCfg *apiConfig) getUserFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	followsFromDB, err := apiCfg.DB.GetFollows(r.Context(), u.ID)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	followsToEncode := make([]follow, 0)
	for _, followFromDB := range followsFromDB {
		followToEncode := follow{
			ID:        followFromDB.ID,
			CreatedAt: followFromDB.CreatedAt,
			UpdatedAt: followFromDB.UpdatedAt,
			UserID:    followFromDB.UserID,
			FeedID:    followFromDB.FeedID,
		}
		followsToEncode = append(followsToEncode, followToEncode)
	}
	writeData, err := json.Marshal(followsToEncode)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(writeData)
}
