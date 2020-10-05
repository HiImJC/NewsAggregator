package aggregator

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/ungerik/go-rss"
	"time"
)

type Article interface {
	json.Marshaler
	ID() string
	PublishDate() (time.Time, error)
}

type BBCArticle struct {
	rss.Item
	additionalCategories []string // Not strictly needed, added to prove off the API
}

func (bbc BBCArticle) PublishDate() (time.Time, error) {
	return bbc.PubDate.ParseWithFormat(time.RFC1123)
}

func (bbc BBCArticle) ID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(bbc.GUID)))
}

func (bbc BBCArticle) MarshalJSON() ([]byte, error) {
	pd, err := bbc.PublishDate()
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		ID            string
		Title         string
		Snippet       string
		PublishDate   string
		Categories    []string
		ThumbnailLink string
		Provider      string
	}{
		ID:            bbc.ID(),
		Title:         bbc.Title,
		Snippet:       bbc.Description,
		PublishDate:   pd.Format(time.RFC1123),
		Categories:    append(bbc.Category, bbc.additionalCategories...),
		ThumbnailLink: "TODO",
		Provider:      "BBCNews",
	})
}

type SkyArticle struct {
	rss.Item
	additionalCategories []string // Not strictly needed, added to prove off the API
}

func (sky SkyArticle) PublishDate() (time.Time, error) {
	return sky.PubDate.ParseWithFormat(time.RFC1123Z)
}

func (sky SkyArticle) ID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(sky.GUID)))
}

func (sky SkyArticle) MarshalJSON() ([]byte, error) {
	pd, err := sky.PublishDate()
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		ID            string
		Title         string
		Snippet       string
		PublishDate   string
		Categories    []string
		ThumbnailLink string
		Provider      string
	}{
		ID:            sky.ID(),
		Title:         sky.Title,
		Snippet:       sky.Description,
		PublishDate:   pd.Format(time.RFC1123),
		Categories:    append(sky.Category, sky.additionalCategories...),
		ThumbnailLink: "TODO",
		Provider:      "SkyNews",
	})
}
