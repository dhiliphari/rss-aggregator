package main

import (
	"context"
	"sync"
	"time"

	"github.com/dhiliphari/rss-aggregator/internal/models"
	"github.com/dhiliphari/rss-aggregator/internal/rss"
)

func (app *application) startScraping(concurrency int, interval time.Duration) {
	app.logger.Info("started scrapping", "concurrency", concurrency, "interval", interval)

	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		feeds, err := app.feeds.GetNextFeedToFetch(context.Background(), concurrency)
		if err != nil {
			app.logger.Error(err.Error())
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Go(func() {
				app.scrapeFeed(feed)
			})
		}
		wg.Wait()
	}
}

func (app *application) scrapeFeed(feed models.Feed) {
	err := app.feeds.MarkFeedFetched(context.Background(), feed.Id)
	if err != nil {
		app.logger.Error("failed marking feed as fetched", "error", err.Error())
		return
	}

	rss, err := rss.FetchFeed(feed.URL)
	if err != nil {
		app.logger.Error("failed fetching feed", "feed_id", feed.Id, "error", err.Error())
		return
	}

	_ = rss
}
