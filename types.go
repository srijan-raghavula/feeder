package main

import (
	"github.com/srijan-raghavula/feeder/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

type errorRes struct {
	Message string `json:"error"`
}

type user struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	ApiKey    string `json:"api_key"`
}
