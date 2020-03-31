package news

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// App queries rss news sources for articles and adds new ones to the database
func App() {
	// connect to redis cache
	articleSet := redis.NewSet(setup.Redis(), os.Getenv("NEWSPAPER_SET"))

	// run the python script
	err := exec.Command(os.Getenv("NEWSPAPER_PYTHON_COMMAND")).Run()
	if err != nil {
		setup.LogCommon(err).Fatal("Failed executing python script")
	}

	// get data from rss feeds
	sources := sourceListFromRSS(articleSet)
	// common driver to turn sources into articles and insert into database
	articleDriver(sources, articleSet)
}

// Uses source objects to call newspaper3k then insert articles into database
func articleDriver(sources []*Source, articleSet *redis.Set) {
	// connect to database
	db := setup.SQL()

	// prep for async calls
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	// call newspaper to extract article data
	newspaper := NewNewspaper()
	for _, s := range sources {
		fmt.Println(s.Title, s.Link, s.Host)

		// call async function
		wg.Add(1)
		go newspaper.Process(s, &wg)
		// random sleep time to be nice
		n := rand.Int63n(5000) // n will be between 0 and 5000
		// base time of one second + [0 - 5] seconds
		// average = 3.5 seconds = 1 second base + 2.5 second expected
		t := time.Second + (time.Millisecond * time.Duration(n))
		time.Sleep(t)
	}
	// collect async data before continuing
	wg.Wait()

	// process each article
	count := 0
	for i, a := range newspaper.Articles {
		// check if we have seen the article before using the canonical link
		if len(a.Canonical) > 1 && articleSet.IsMember(a.Canonical) {
			// skip
		} else {
			// create standard article and put it in database
			article := NewArticle(a, newspaper.Sources[i])
			article.Insert(db)
			// cache known links so we don't duplicate articles
			articleSet.Add(newspaper.Sources[i].Link)
			articleSet.Add(a.Canonical)
			count = count + 1
		}
	}

	// run summary
	setup.LogCommon(nil).
		WithField("NumSources", len(sources)).
		WithField("NumArticles", len(newspaper.Articles)).
		WithField("NumInserted", count).
		WithField("RunTime", setup.RunTime().String()).
		Info("RunSummary")
}

// Returns source objects from rss feeds
func sourceListFromRSS(articleSet *redis.Set) []*Source {
	// get rss feeds
	rss := NewFeeds()

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
