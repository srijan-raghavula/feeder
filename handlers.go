package main

import (
	"encoding/json"
	"net/http"
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
