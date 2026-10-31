package main

import (
	"context"
	"strings"
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

	for _, item := range rss.Channel.Item {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			pubDate, err = time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				app.logger.Error("error parsing published date", "published_date", item.PubDate, "error", err.Error())
				continue
			}
		}
		_, err = app.posts.Create(context.Background(), item.Title, item.Description, item.Link, pubDate, feed.Id)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			app.logger.Error("failed to create post", "error", err.Error())
		}
	}
}
