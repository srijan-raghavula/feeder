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

	"github.com/google/uuid"
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

func fetch(url string) (feeds.RssFeedXml, error) {
	var feed feeds.RssFeedXml
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
		go func() {
			apiCfg.fetchRoutine(feed.Url)
			defer wg.Done()
		}()
		log.Println("Fetcheing next")
	}
	wg.Wait()
	return true
}

func (apiCfg *apiConfig) fetchRoutine(url string) {
	ctx := context.Background()
	log.Println("Fetching", url)
	feed, err := fetch(url)
	if err != nil {
		log.Println(err.Error())
		log.Println("Fetch failed:", url)
		return
	}
	markFetchedParams := database.MarkFeedFetchedParams{
		Url: url,
		LastFetchedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}
	apiCfg.DB.MarkFeedFetched(ctx, markFetchedParams)
	feedFromDB, err := apiCfg.DB.GetFeedByUrl(ctx, url)
	if err != nil {
		log.Println(err.Error())
	}
	post, err := apiCfg.DB.GetPostByUrl(ctx, sql.NullString{
		Valid:  true,
		String: url,
	})
	if err == sql.ErrNoRows {
		apiCfg.DB.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New().String(),
			CreatedAt: time.Now().UTC().String(),
			UpdatedAt: time.Now().UTC().String(),
			Title:     feed.Channel.Title,
			Url: sql.NullString{
				Valid:  true,
				String: url,
			},
			Description: sql.NullString{
				Valid:  true,
				String: feed.Channel.Description,
			},
			PublishedAt: sql.NullString{
				Valid:  true,
				String: feed.Channel.PubDate,
			},
			FeedID: feedFromDB.ID,
		})
		return
	} else if err != nil {
		log.Println(err.Error())
		return
	}
	post.UpdatedAt = time.Now().UTC().String()
	apiCfg.DB.UpdatePostByUrl(ctx, database.UpdatePostByUrlParams{
		Url: sql.NullString{
			Valid:  true,
			String: url,
		},
		UpdatedAt: time.Now().UTC().String(),
	})
}
