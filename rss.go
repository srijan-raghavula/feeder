package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/feeds"
	"github.com/srijan-raghavula/feeder/internal/database"
)

func (apiCfg *apiConfig) workerLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		ctx := context.Background()
		go apiCfg.worker(ctx)
	}
}

func fetch(url string) (feeds.Feed, error) {
	var feed feeds.Feed
	res, err := http.Get(url)
	if err != nil {
		return feed, err
	}
	defer res.Body.Close()
	if res.StatusCode > 399 {
		return feed, errors.New(res.Status)
	}

	err = xml.NewDecoder(res.Body).Decode(&feed)
	if err != nil {
		return feed, err
	}

	return feed, nil
}

func (apiCfg *apiConfig) worker(ctx context.Context) bool {
	wg := sync.WaitGroup{}
	limit := int32(10)
	feeds, err := apiCfg.DB.GetNextFeedsToFetch(ctx, limit)
	if err != nil {
		log.Println(err.Error())
	}
	for _, feed := range feeds {
		wg.Add(1)
		go func(url string) {
			log.Println("Fetching", url)
			defer wg.Done()
			feed, err := fetch(url)
			if err != nil {
				log.Println(err.Error())
				log.Println("Fetch failed:", url)
			} else {
				markFetchedParams := database.MarkFeedFetchedParams{
					Url: url,
					LastFetchedAt: sql.NullTime{
						Valid: true,
						Time:  time.Now(),
					},
				}
				apiCfg.DB.MarkFeedFetched(ctx, markFetchedParams)
				log.Println("Fetched", url, ":\n", feed)
			}
		}(feed.Url)
		log.Println("Fetcheing next")
	}
	wg.Wait()
	return true
}
