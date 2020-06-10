package text

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/wpwilson10/caterpillar/internal/news"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// AdjacentArticles returns articles that were published or entered immediatlely before and after
// the given target article.
func AdjacentArticles(db *sqlx.DB, target *news.Article) []news.Article {
	// First use published Time from article if it exists
	articles := byPublishedTime(db, target)
	if articles != nil {
		return articles
	}
	// Then try source published time
	articles = bySourcePublishedTime(db, target)
	if articles != nil {
		return articles
	}
	// Finally use the data entry time
	return byDataEntryTime(db, target)
}

// Returns the articles with a published_time immediately before and after the given article,
// or null if the given article has no published_time.
func byPublishedTime(db *sqlx.DB, target *news.Article) []news.Article {
	// check that we have a published time
	if target.PublishedTime.IsZero() {
		return nil
	}

	var selectStmtBefore string = `SELECT * FROM NewsArticle 
									WHERE article_id !=$1 AND host=$2
										AND published_time <= $3 
									ORDER BY published_time DESC
									LIMIT 10`

	var selectStmtAfter string = `SELECT * FROM NewsArticle 
									WHERE article_id !=$1 AND host=$2
										AND published_time >= $3 
									ORDER BY published_time ASC
									LIMIT 5`

	// we know this value exists from IsZero check above
	articleTime := target.PublishedTime.ValueOrZero()

	return selectBeforeAndAfter(db, target, &articleTime, selectStmtBefore, selectStmtAfter)
}

// Returns the articles with a source_published_time immediately before and after the given article,
// or null if the given article has no source_published_time.
func bySourcePublishedTime(db *sqlx.DB, target *news.Article) []news.Article {
	// check that we have a published time
	if target.SourcePublishedTime.IsZero() {
		return nil
	}

	var selectStmtBefore string = `SELECT * FROM NewsArticle 
									WHERE article_id !=$1 AND host=$2
										AND source_published_time <= $3 
									ORDER BY source_published_time DESC
									LIMIT 10`

	var selectStmtAfter string = `SELECT * FROM NewsArticle 
									WHERE article_id !=$1 AND host=$2
										AND source_published_time >= $3 
									ORDER BY source_published_time ASC
									LIMIT 5`

	// we know this value exists from IsZero check above
	articleTime := target.SourcePublishedTime.ValueOrZero()

	return selectBeforeAndAfter(db, target, &articleTime, selectStmtBefore, selectStmtAfter)
}

// Returns the articles with a data_entry_time immediately before and after the given article.
// Should always return values since data_entry_time is a default field.
func byDataEntryTime(db *sqlx.DB, target *news.Article) []news.Article {
	var selectStmtBefore string = `SELECT * FROM NewsArticle 
									WHERE article_id !=$1 AND host=$2
										AND data_entry_time <= $3 
									ORDER BY data_entry_time DESC
									LIMIT 10`

	var selectStmtAfter string = `SELECT * FROM NewsArticle 
									WHERE article_id !=$1 AND host=$2
										AND data_entry_time >= $3 
									ORDER BY data_entry_time ASC
									LIMIT 5`

	return selectBeforeAndAfter(db, target, &target.DataTime, selectStmtBefore, selectStmtAfter)
}

// Uses the given select statements and article time to perform the database search.
// Returns a combined list of articles corresponding to articles immediately before and after the articeTime.
func selectBeforeAndAfter(db *sqlx.DB, target *news.Article, articleTime *time.Time, selectStmtBefore string, selectStmtAfter string) []news.Article {
	articlesBefore := []news.Article{}
	articlesAfter := []news.Article{}

	err := db.Select(&articlesBefore, selectStmtBefore, target.ArticleID, target.Host, articleTime)
	if err != nil {
		setup.LogCommon(err).Error("Select articles before")
	}

	err = db.Select(&articlesAfter, selectStmtAfter, target.ArticleID, target.Host, articleTime)
	if err != nil {
		setup.LogCommon(err).Error("Select articles after")
	}

	return append(articlesAfter, articlesBefore...)
}
