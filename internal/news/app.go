package news

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// App queries rss news sources for articles and adds new ones to the database
func App() {
	// connect to redis cache
	articleSet := redis.NewSet(setup.Redis(), os.Getenv("NEWSPAPER_SET"))
	// connect to database
	db := setup.SQL()

	// check if python script is running
	if !setup.CheckOnce(setup.EnvToInt("PY_NEWSPAPER_PORT")) {
		setup.LogCommon(nil).Fatal("Python app not running")
	}

	// get data from rss feeds
	sources := SourceListFromRSS(articleSet)
	// common driver to turn sources into articles and insert into database
	articleDriver(sources, articleSet, db)
}

// RedditLinksApp queries reddit submissions for articles and adds new ones to the database
func RedditLinksApp() {
	// connect to redis cache
	articleSet := redis.NewSet(setup.Redis(), os.Getenv("NEWSPAPER_SET"))
	// connect to database
	db := setup.SQL()

	// check if python script is running
	if !setup.CheckOnce(setup.EnvToInt("PY_NEWSPAPER_PORT")) {
		setup.LogCommon(nil).Fatal("Python app not running")
	}

	// get data from reddit submissions
	sources := SourceListFromReddit(articleSet, db)
	// common driver to turn sources into articles and insert into database
	articleDriver(sources, articleSet, db)
}

// Uses source objects to call newspaper3k then insert articles into database
func articleDriver(sources []*Source, articleSet *redis.Set, db *sqlx.DB) {
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
