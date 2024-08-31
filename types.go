package main

import (
	"github.com/gorilla/feeds"
	"net/http"
	"time"

	"github.com/srijan-raghavula/feeder/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

type errorRes struct {
	Message string `json:"error"`
}

type user struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	ApiKey    string `json:"api_key"`
}

type feed struct {
	ID            string    `josn:"id"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	Name          string    `json:"name"`
	Url           string    `json:"url"`
	UserID        string    `json:"user_id"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
}

type follow struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	FeedID    string `json:"feed_id"`
	UserID    string `json:"user_id"`
}

type feedPage struct {
	Title         string       `xml:"title"`
	Link          string       `xml:"link"`
	Description   string       `xml:"description"`
	Generator     string       `xml:"generator"`
	Lang          string       `xml:"language"`
	LastBuildDate string       `xml:"lastBuildDat"`
	Items         []feeds.Feed `xml:"item"`
}

type authedUserHandler func(http.ResponseWriter, *http.Request, database.User)
