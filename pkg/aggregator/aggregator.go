package aggregator

import (
	"NewsAggregator/pkg/cache"
	"context"
	"errors"
	"fmt"
	"github.com/ungerik/go-rss"
	"log"
	"net/url"
	"os"
	"time"
)

type Aggregator struct {
	logger       *log.Logger
	channelCache map[string]*cache.Cache
	feeds        Feeds
	cancelFunc   context.CancelFunc
}

type Feed struct {
	URL        string
	Categories []string
}
type Feeds []Feed

func New(f Feeds, tick time.Duration) Aggregator {
	ctx, cancel := context.WithCancel(context.Background())

	a := Aggregator{
		logger:       log.New(os.Stdout, "[Aggregator] ", log.LstdFlags),
		channelCache: make(map[string]*cache.Cache),
		feeds:        f,
		cancelFunc:   cancel,
	}

	go a.Start(ctx, tick)

	return a
}

func (a Aggregator) Start(ctx context.Context, tick time.Duration) {
	a.RefreshData()

	ticker := time.Tick(tick)
	select {
	case <-ticker:
		a.logger.Println("Refreshing feed data")
		a.RefreshData()
	case <-ctx.Done():
		a.logger.Println("Aggregator has been stopped")
		return
	}
}

func (a Aggregator) Stop() {
	a.cancelFunc()
}

type data struct {
	Channel string
	Art     []Article `json:"Articles"`
}

func (a Aggregator) GetLatestData() []data {
	allData := make([]data, 0, 0)

	for channel, articles := range a.channelCache {
		d := data{Channel: channel, Art: make([]Article, 0, 0)}

		for _, k := range articles.Keys() {
			ia, ok := articles.Get(k)
			if !ok {
				continue
			}

			d.Art = append(d.Art, ia.(Article))
		}

		allData = append(allData, d)
	}

	return allData
}

func (a Aggregator) RefreshData() {
	for _, f := range a.feeds {
		a.logger.Println("fetching data for feed: " + f.URL)

		url, err := url.Parse(f.URL)
		if err != nil {
			panic(err)
		}

		//Always pass false here as Reddit is not a supported RSS host
		resp, err := rss.Read(f.URL, false)
		if err != nil {
			panic(err)
		}

		channel, err := rss.Regular(resp)
		if err != nil {
			panic(err)
		}

		a.logger.Println("channel is: " + channel.Description)

		// channel is a single RSS feed, containing multiple items.
		channelCache, ok := a.channelCache[f.URL]
		if !ok {
			// This is the first time we have seen this channel,
			// create a cache.
			channelCache = cache.New()
			a.channelCache[f.URL] = channelCache
		}

		for _, item := range channel.Item {
			article, err := convertItemToArticle(url, item, f.Categories)
			if err != nil {
				a.logger.Println(err)
				break
			}

			// Add or update item in the cache.
			updated, err := addArticleToCache(channelCache, article)
			if err != nil {
				a.logger.Println("unable to add value to the cache: " + err.Error() + " SKIPPING")
				continue
			}

			pubdate, _ := article.PublishDate()

			// each item is an individual news article
			a.logger.Println(fmt.Sprintf("\t%s %s %s UPDATE: %v", pubdate, item.Title, item.Link, updated))
		}

	}
}

func addArticleToCache(c *cache.Cache, a Article) (bool, error) {
	oldArticle, ok := c.Get(a.ID())
	if !ok {
		// This article doesn't currently exist in the cache
		// add it
		c.Put(a.ID(), a)

		return false, nil
	}

	// Article already exists, lets check if the article we received
	// has a later timestamp
	newPublishTime, err := a.PublishDate()
	if err != nil {
		return false, fmt.Errorf("unable to parse new article time: %w", err)
	}

	oldPublishTime, err := oldArticle.(Article).PublishDate()
	if err != nil {
		return false, fmt.Errorf("unable to parse old article time: %w", err)
	}

	if newPublishTime.After(oldPublishTime) {
		c.Put(a.ID(), a)
		return true, nil
	}

	// Article we have received has a publish date prior to
	// the article we already have.
	return false, nil
}

func convertItemToArticle(feed *url.URL, item rss.Item, additionalCategories []string) (Article, error) {
	switch feed.Host {
	case "feeds.bbci.co.uk":
		return BBCArticle{item, additionalCategories}, nil
	case "feeds.skynews.com":
		return SkyArticle{item, additionalCategories}, nil
	}

	return nil, errors.New("Unsupported feed type: " + feed.Host)
}
