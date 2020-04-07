package news

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/mmcdole/gofeed"

	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// SourceListFromRSS returns source objects from rss feeds.
func SourceListFromRSS(articleSet *redis.Set) []*Source {
	// get rss feeds
	rss := newFeeds()

	// iterate through each news feed to create source structs
	// and filter out links we have already seen
	sources := []*Source{}
	for _, r := range rss {
		// iterate through each article in feed
		for _, a := range r.Items {
			// check if we have seen this link before
			if !articleSet.IsMember(a.Link) {
				// convert feed article to standard source
				source := NewSource(FromFeed(a))
				// add to list if we got something
				if len(source.Link) > 1 {
					sources = append(sources, source)
				}
			}
		}
	}

	return sources
}

// newFeeds returns an array of gofeeds using the rss list in the file at RSS_FILEPATH.
// This is potentially slow
func newFeeds() []*gofeed.Feed {
	// get rss feeds from file
	rss := rssFromFile()
	// setup parser
	fp := gofeed.NewParser()
	// set timeout
	fp.Client = &http.Client{Timeout: time.Second * 10}
	// save values
	out := []*gofeed.Feed{}

	// process each
	for _, r := range rss {
		// sanity check that there is an rss link, arbitrary length value used
		if len(r.RSS) < 5 {
			setup.LogCommon(nil).
				WithField("RSS", r.Name).
				Error("No RSS feed")
		}
		// parse the rss feed
		feed, err := fp.ParseURL(strings.TrimSpace(r.RSS))
		if err != nil {
			setup.LogCommon(err).
				WithField("RSS", r.Name).
				WithField("URL", r.RSS).
				Error("Failed ParseURL")
		}
		// check we got good data
		err = checkRSSFeed(feed)
		if err != nil {
			setup.LogCommon(err).
				WithField("RSS", r.Name).
				WithField("URL", r.RSS).
				Error("Failed checkRSSFeed")
		} else {
			// add to output
			out = append(out, feed)
		}
	}

	return out
}

// checkRSSFeed performs basic checks to see if feed is valid.
func checkRSSFeed(feed *gofeed.Feed) error {
	// check we got something
	if feed == nil {
		return errors.New("No feed returned")
	} else if len(feed.Items) < 1 {
		// check we got some articles
		return errors.New("Feed returned no articles")
	} else if len(feed.Title) < 1 {
		// check there is a title
		return errors.New("Feed has no title")
	}
	return nil
}

// temp struct to parse rss source file
type rss struct {
	Name     string `csv:"Name"`
	Link     string `csv:"Link"`
	RSS      string `csv:"RSS"`
	Category string `csv:"Category"`
}

// rssFromFile creates an array of rss structs from the file at filepath.
func rssFromFile() []*rss {
	filepath := os.Getenv("NEWSPAPER_RSS_FILEPATH")
	// Get the file
	sourceFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		setup.LogCommon(err).
			WithField("filepath", filepath).
			Error("OpenFile")
	}
	defer sourceFile.Close()

	// Read from file to array
	s := []*rss{}
	// Load clients from file
	if err := gocsv.UnmarshalFile(sourceFile, &s); err != nil {
		setup.LogCommon(err).
			WithField("filepath", filepath).
			Error("UnmarshalFile")
	}

	// Check that we got something
	if len(s) <= 0 {
		setup.LogCommon(nil).
			WithField("filepath", filepath).
			Error("Empty source list")
	}

	return s
}
