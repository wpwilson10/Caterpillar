package news

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Article stores a news article in a format matching the NewsArticle table schema
type Article struct {
	ArticleID           int64       `db:"article_id"`
	DataTime            time.Time   `db:"data_entry_time"`
	Source              string      `db:"source"`                // source of the link (e.g. RSS, Reddit)
	Host                string      `db:"host"`                  // host or vendor of article (e.g. www.cnn.com, www.espn.com)
	Link                string      `db:"link"`                  // original url from the reference
	SourcePublishedTime null.Time   `db:"source_published_time"` // published time from source if given
	PublishedTime       null.Time   `db:"published_time"`        // published time from newspaper3k if given
	SourceTitle         null.String `db:"source_title"`
	Title               null.String `db:"title"`
	CanonicalLink       null.String `db:"canonical_link"`
	Body                null.String `db:"body"`
	Authors             null.String `db:"authors"`
}

// NewArticle parses data from newspaper3k and source into a standard article.
func NewArticle(raw *RawArticle, source *Source) *Article {
	var article = Article{
		DataTime:            time.Now(),
		Source:              source.Source,
		Host:                source.Host,
		Link:                source.Link,
		SourcePublishedTime: sourcePublishedTime(source.PubDate),
		PublishedTime:       publishedTime(raw.PubDate, source),
		SourceTitle:         sourceTitle(source.Title),
		Title:               title(raw.Title),
		CanonicalLink:       canonicalLink(raw.Canonical),
		Body:                body(raw.Text, source),
		Authors:             authors(raw.Authors, source),
	}

	return &article
}

// Attribute Setters

func authors(list []string, source *Source) null.String {
	// turn into json
	j, err := json.Marshal(list)

	if err != nil {
		setup.LogCommon(err).
			WithField("link", source.Link).
			WithField("authors", list).
			Warn("Failed json.Marshal")
	}

	// convert from bytes to string
	s := string(j)

	return null.StringFromPtr(&s)
}

func body(s string, source *Source) null.String {
	clean := strings.TrimSpace(s)
	// check that we got a reasonable amount of text
	// arbitrary lengths
	if len(clean) < 3 {
		return null.NewString("", false)
	} else if len(clean) < 30 {
		setup.LogCommon(nil).
			WithField("link", source.Link).
			WithField("length", len(clean)).
			Warn("Short article")
	}

	return null.StringFromPtr(&clean)
}

func canonicalLink(s string) null.String {
	// not any different than title currently
	return title(s)
}

func sourceTitle(s string) null.String {
	// not any different than title currently
	return title(s)
}

func title(s string) null.String {
	clean := strings.TrimSpace(s)
	// check that we got a reasonable title
	// 280 is twitter max tweet size, means this probably isn't a title
	if len(clean) < 3 || len(clean) > 280 {
		return null.NewString("", false)
	}

	return null.StringFromPtr(&clean)
}

func sourcePublishedTime(t *time.Time) null.Time {
	return null.TimeFromPtr(t)
}

func publishedTime(t string, s *Source) null.Time {
	// published time from newspaper is not reliable so pubDate may be null
	// case where we know the time parse will not work, so return a null time
	if len(t) != len("2006-01-02T15:04:05-07:00") {
		return null.NewTime(time.Time{}, false)
	}

	// use actual time instead of placeholder
	pubDate, err := time.Parse("2006-01-02T15:04:05-07:00", t)
	if err != nil {
		setup.LogCommon(err).
			WithField("Link", s.Link).
			WithField("pubDate", t).
			Warn("Failed time.Parse")
	}
	// this constructor handles cases where pubDate is nil
	return null.TimeFromPtr(&pubDate)
}

// Insert adds this article to the NewsArticle database table.
// Performs no validation.
func (article *Article) Insert(db *sqlx.DB) {
	// Setup
	var insertStmt string = `INSERT INTO NewsArticle (
								article_id, 		-- DEFAULT
								data_entry_time, 	-- Now()
								source,
								host,
								link,
								source_published_time,
								published_time,
								source_title,
								title,
								canonical_link,
								body,
								authors
								)`

	var valueStmt string = `VALUES (
								DEFAULT,
								Now(),
								:source,
								:host,
								:link,
								:source_published_time,
								:published_time,
								:source_title,
								:title,
								:canonical_link,
								:body,
								:authors
								)`

	var fullStmt string = insertStmt + " " + valueStmt

	// Insert article
	_, err := db.NamedExec(fullStmt, article)

	if err != nil {
		setup.LogCommon(err).
			WithField("Link", article.Link).
			Error("db.NamedExec")
	}

	setup.LogCommon(nil).
		WithField("Link", article.Link).
		Info("Inserting article")

}
