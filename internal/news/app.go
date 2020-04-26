package news

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/turnage/graw/reddit"
	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// App queries rss news sources for articles and adds new ones to the database
func App() {
	// connect to redis cache
	articleSet := redis.NewSet(setup.Redis(), os.Getenv("NEWSPAPER_SET"))
	// connect to database
	db := setup.SQL()
	// setup blacklist of article hosts to avoid
	blacklist := NewBlackList()
	// randomize time between calls
	rand.Seed(time.Now().UnixNano())
	// prep for async calls
	var wg sync.WaitGroup
	var numArticles uint64

	// check if python script is running
	if !setup.CheckOnce(setup.EnvToInt("PY_NEWSPAPER_PORT")) {
		setup.LogCommon(nil).Fatal("Python app not running")
	}

	// get data from rss feeds
	sources := SourceListFromRSS(articleSet)
	// process each source
	for _, source := range sources {
		// async parts - hands off a source for processing
		go func(source *Source) {
			defer wg.Done()
			// do the work of getting data and saving it
			a := Driver(source, db, articleSet, blacklist)
			// count number of articles successfully returned
			if a != nil {
				atomic.AddUint64(&numArticles, 1)
			}

		}(source)

		// random sleep time to be nice
		n := rand.Int63n(5000) // n will be between 0 and 5000
		// base time of one second + [0 - 5] seconds
		// average = 3.5 seconds = 1 second base + 2.5 second expected
		t := time.Second + (time.Millisecond * time.Duration(n))
		time.Sleep(t)
	}

	// wait to finish
	wg.Wait()

	// run summary
	setup.LogCommon(nil).
		WithField("NumSources", len(sources)).
		WithField("NumArticles", numArticles).
		WithField("RunTime", setup.RunTime().String()).
		Info("RunSummary")
}

// Driver uses a source to retrieve article data and save it into the database.
// Article will have be inserted into database and cached on successful calls.
// Returns nil if we have seen article before or failing to get or process article.
func Driver(source *Source, db *sqlx.DB, articleSet *redis.Set, blacklast *BlackList) *Article {
	// check if we have seen this source before
	if len(source.Link) > 1 && articleSet.IsMember(source.Link) {
		// skip
		return nil
	}
	// check if host is on blacklist
	if blacklast.IsBlackListed(source.Host) {
		// host may be updated to canonical, check if this is blacklisted
		fmt.Println("Black listed", source.Host, source.Source)
		// skip
		return nil
	}

	// call newspaper3k to get data
	newspaper := NewNewspaper(source)
	// check we got something
	if newspaper == nil {
		return nil
	}

	// check if we have seen the article before using the canonical link
	if len(newspaper.Canonical) > 1 && articleSet.IsMember(newspaper.Canonical) {
		// skip
		return nil
	}

	// create standard article
	article := NewArticle(newspaper, source)

	// check we got data worth inserting
	if article.Body.IsZero() {
		// no body text
		return nil
	} else if article.SourceTitle.IsZero() && article.Title.IsZero() {
		// no titles
		return nil
	} else if blacklast.IsBlackListed(article.Host) {
		// host may be updated to canonical, check if this is blacklisted
		fmt.Println("Black listed", article.Link, source.Host, source.Source)
		return nil
	}

	// Put article in database
	article.Insert(db)
	// cache known links so we don't duplicate articles
	articleSet.Add(source.Link)
	articleSet.Add(newspaper.Canonical)

	return article
}

// RedditNewsDriver adds news articles from reddit posts to the NewsArticle database
// and adds a RedditNews relationship entry to the RedditNews table.
func RedditNewsDriver(db *sqlx.DB, articleSet *redis.Set, blacklist *BlackList, submission *reddit.Post, sID int64) {
	// quick initial check that submissions have a link
	if len(submission.URL) <= 2 {
		// dont log error because it is normal for submissions to not have external link
		return
	}

	// submission that has been seen before
	if articleSet.IsMember(submission.URL) {
		// get the previous submission
		article := FindArticle(db, submission.URL)
		// sanity check that article exists
		if article == nil {
			setup.LogCommon(nil).
				WithField("redditID", submission.ID).
				WithField("redditURL", submission.URL).
				WithField("SubmissionID", sID).
				Error("Failed to find article in articleSet")
		} else {
			fmt.Println("Found existing article", article.ArticleID, article.Link, submission.Permalink)
			// create a table entry and insert it
			rn := NewRedditNews(article.ArticleID, sID)
			rn.Insert(db)
		}
	} else {
		// submission that has not been seen before
		// put into form ArticleDriver expects
		source := NewSource(FromReddit(submission))
		// get article and add to database
		article := Driver(source, db, articleSet, blacklist)
		// check that article exists
		if article != nil {
			// create a table entry and insert it
			rn := NewRedditNews(article.ArticleID, sID)
			rn.Insert(db)
		}
	}
}
