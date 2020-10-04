package aggregator

import (
	"github.com/ungerik/go-rss"
	"time"
)

type Article interface {
	ID() string
	PublishDate() (time.Time, error)
}

type BBCArticle struct {
	rss.Item
}

func (bbc BBCArticle) PublishDate() (time.Time, error) {
	return bbc.PubDate.ParseWithFormat(time.RFC1123)
}

func (bbc BBCArticle) ID() string {
	return bbc.GUID
}

type SkyArticle struct {
	rss.Item
}

func (sky SkyArticle) PublishDate() (time.Time, error) {
	return sky.PubDate.ParseWithFormat(time.RFC1123Z)
}

func (sky SkyArticle) ID() string {
	return sky.GUID
}
