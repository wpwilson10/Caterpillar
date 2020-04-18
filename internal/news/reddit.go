package news

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// RedditNews represents the database table for linking reddit submissions to news articles.
type RedditNews struct {
	LinkID       int64     `db:"link_id"`
	ArticleID    int64     `db:"article_id"`
	SubmissionID int64     `db:"submission_id"`
	DataTime     time.Time `db:"data_entry_time"`
}

// NewRedditNews creates an entry for the RedditNews database table
func NewRedditNews(articleID int64, submissionID int64) *RedditNews {
	var link = RedditNews{
		ArticleID:    articleID,
		SubmissionID: submissionID,
	}

	return &link
}

// Insert adds this redditnews relationship to the database table.
// Performs no validation.
func (link *RedditNews) Insert(db *sqlx.DB) {
	// Setup
	var insertStmt string = `INSERT INTO RedditNews (
								data_entry_time, 	-- Now()
								article_id, 		-- $1
								submission_id		-- $2
								)`

	var valueStmt string = `VALUES (
								Now(),
								:article_id,
								:submission_id
								)`

	var fullStmt string = insertStmt + " " + valueStmt

	// Insert article
	_, err := db.NamedExec(fullStmt, link)

	if err != nil {
		setup.LogCommon(err).
			WithField("ArticleID", link.ArticleID).
			WithField("SubmissionID", link.SubmissionID).
			Error("db.NamedExec")
	}

	setup.LogCommon(nil).
		WithField("ArticleID", link.ArticleID).
		WithField("SubmissionID", link.SubmissionID).
		Info("Inserting article")

}

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
			source := NewSource(FromRedditArticle(r))
			// add to list if we got something
			if len(source.Link) > 1 {
				sources = append(sources, source)
			}
		}
	}

	return sources
}

// For future backfills below

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
