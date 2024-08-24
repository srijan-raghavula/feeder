package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/srijan-raghavula/feeder/internal/database"
	"log"
	"net/http"
	"os"
)

func main() {
	godotenv.Load()

	// db shit
	dbURL := os.Getenv("CONN_STR")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("error opening db:", err.Error())
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		DB: dbQueries,
	}

	// actual server shit
	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	// idle handler
	mux.HandleFunc("/", welcome)
	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: mux,
	}

	// server health related endpoints
	mux.HandleFunc("GET /v1/healthz", readiness)
	mux.HandleFunc("GET /v1/err", errHandler)

	// user creation and getting
	mux.HandleFunc("POST /v1/users", apiCfg.createUser)
	mux.HandleFunc("GET /v1/users", apiCfg.midAuth(apiCfg.getUser))

	// feed creation and getting
	mux.HandleFunc("POST /v1/feeds", apiCfg.midAuth(apiCfg.createFeed))
	mux.HandleFunc("GET /v1/feeds", apiCfg.getAllFeeds)

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
