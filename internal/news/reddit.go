package news

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// RedditArticle represents external links from reddit submissions.
type RedditArticle struct {
	Title   string     `db:"title"`
	Link    string     `db:"url"`
	PubDate *time.Time `db:"created_time"`
}

func NewRedditArticles(db *sqlx.DB) []RedditArticle {
	// --- Get Submissions
	var selectStmt string = `SELECT title, url, created_time, data_entry_time FROM RedditSubmission
							 WHERE data_entry_time BETWEEN CURRENT_DATE - 1 AND CURRENT_DATE - 2'
							 	AND LENGTH(url) > 2
							 ORDER BY created_time ASC`

	submissions := []RedditArticle{}

	// pull submissions from database
	err := db.Select(&submissions, selectStmt)
	if err != nil {
		setup.LogCommon(err).Fatal("Failed Select Statement")
	}

	return submissions
}
