package news

import (
	"net/url"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Source is a website news source.
type Source struct {
	Title   string
	Link    string // URL of content
	Source  string // where the article link came from (e.g. RSS Feed, Reddit)
	Host    string // hostname parsed from the link
	PubDate *time.Time
}

// SourceOption is the signature of our option functions
// see https://www.sohamkamani.com/blog/golang/options-pattern/
type SourceOption func(*Source)

// NewSource creates a standard Source from multiple possible reference structs.
func NewSource(opts ...SourceOption) *Source {
	// Need a default struct
	s := &Source{}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *House as the argument
		opt(s)
	}

	return s
}

// FromFeed uses gofeed.Items to make a source
func FromFeed(feed *gofeed.Item) SourceOption {
	// get host from link
	u, err := url.Parse(feed.Link)
	if err != nil {
		setup.LogCommon(err).
			WithField("link", feed.Link).
			Warn("Failed url.Parse")
	}

	return func(s *Source) {
		s.Title = feed.Title
		s.Link = feed.Link
		s.Source = "RSS Feed"
		s.Host = u.Hostname()
		s.PubDate = feed.PublishedParsed
	}
}

// FromReddit uses reddit submissions to make a source
func FromReddit(item *RedditArticle) SourceOption {
	// get host from link
	u, err := url.Parse(item.Link)
	if err != nil {
		setup.LogCommon(err).
			WithField("link", item.Link).
			Warn("Failed url.Parse")
	}

	return func(s *Source) {
		s.Title = item.Title
		s.Link = item.Link
		s.Source = "Reddit Submission"
		s.Host = u.Hostname()
		s.PubDate = item.PubDate
	}
}
