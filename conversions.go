package main

import (
	"github.com/srijan-raghavula/feeder/internal/database"
)

func feedFromDBtoEnc(feedFromDB database.Feed) feed {
	return feed{
		ID:            feedFromDB.ID,
		CreatedAt:     feedFromDB.CreatedAt,
		UpdatedAt:     feedFromDB.UpdatedAt,
		Name:          feedFromDB.Name,
		Url:           feedFromDB.Url,
		UserID:        feedFromDB.UserID,
		LastFetchedAt: feedFromDB.LastFetchedAt.Time,
	}
}

func userFromDBtoEnc(userFromDB database.User) user {
	return user{
		ID:        userFromDB.ID,
		CreatedAt: userFromDB.CreatedAt,
		UpdatedAt: userFromDB.UpdatedAt,
		Name:      userFromDB.Name,
		ApiKey:    userFromDB.ApiKey,
	}

}
