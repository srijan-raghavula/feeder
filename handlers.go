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
		Id:        userFromDB.ID,
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

func (cfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	apiKey := strings.TrimPrefix(header, "ApiKey ")
	userFromDB, err := cfg.DB.GetUserByAPI(r.Context(), apiKey)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	userToEncode := user{
		Id:        userFromDB.ID,
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
