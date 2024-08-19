package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	mux.HandleFunc("/", welcome)
	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: mux,
	}

	mux.HandleFunc("GET /v1/healthz", readiness)
	mux.HandleFunc("GET /v1/err", errHandler)

	log.Println("Listening and serving at port:", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func resWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	writeData, err := json.Marshal(payload)
	if err != nil {
		resWithError(w, 500, err.Error())
		return
	}
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	w.Write(writeData)
}

func resWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	writeData, err := json.Marshal(errorRes{
		Message: msg,
	})
	if err != nil {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(writeData)
}
