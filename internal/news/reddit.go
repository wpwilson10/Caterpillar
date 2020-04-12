package news

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// SourceListFromReddit returns source objects from reddit submissions.
func SourceListFromReddit(articleSet *redis.Set, db *sqlx.DB) []*Source {
	// get links from reddit submissions
	redditArticles := newRedditArticles(db)

	// iterate through each reddit link to create source structs
	// and filter out links we have already seen
	sources := []*Source{}
	for _, r := range redditArticles {
		// check if we have seen this link before
		if !articleSet.IsMember(r.Link) {
			// convert feed article to standard source
			source := NewSource(FromReddit(r))
			// add to list if we got something
			if len(source.Link) > 1 {
				sources = append(sources, source)
			}
		}
	}

	return sources
}

// RedditArticle represents external links from reddit submissions.
type RedditArticle struct {
	Title    string     `db:"title"`
	Link     string     `db:"url"`
	PubDate  *time.Time `db:"created_time"`
	DataTime *time.Time `db:"data_entry_time"`
}

// newRedditArticles returns a list of objects representing valid links from reddit submissions.
func newRedditArticles(db *sqlx.DB) []*RedditArticle {
	// --- Get Submissions
	/* var selectStmt string = `SELECT title, url, created_time, data_entry_time FROM RedditSubmission
	WHERE data_entry_time BETWEEN CURRENT_DATE - 2 AND CURRENT_DATE - 1
		AND LENGTH(url) > 2
	ORDER BY created_time ASC`
	*/

	var selectStmt string = `SELECT title, url, created_time, data_entry_time FROM RedditSubmission
							 WHERE LENGTH(url) > 2
							 ORDER BY created_time ASC`

	submissions := []RedditArticle{}

	// pull submissions from database
	err := db.Select(&submissions, selectStmt)
	if err != nil {
		setup.LogCommon(err).Fatal("Failed Select Statement")
	}

	// url cleanup and validation that URL is good
	out := []*RedditArticle{}
	for _, s := range submissions {
		s.Link = strings.TrimSpace(s.Link)
		if setup.IsValidURL(s.Link) {
			out = append(out, &s)
		}
	}

	return out
}
