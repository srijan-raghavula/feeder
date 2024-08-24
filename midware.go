package main

import (
	"context"
	"github.com/srijan-raghavula/feeder/internal/database"
	"net/http"
)

func (cfg *apiConfig) midAuth(handler authedUserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := getApiKeyFromReq(r)
		u, valid, err := cfg.authUser(r.Context(), apiKey)
		if err != nil || !valid {
			resWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		handler(w, r, u)
	}
}

func (cfg *apiConfig) authUser(context context.Context, apiKey string) (database.User, bool, error) {
	userFromDB, err := cfg.DB.GetUserByAPI(context, apiKey)
	if err != nil {
		return userFromDB, false, err
	}
	if apiKey != userFromDB.ApiKey {
		return userFromDB, false, err
	}
	return userFromDB, true, nil
}
