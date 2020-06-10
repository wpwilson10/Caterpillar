package text

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/wpwilson10/caterpillar/internal/news"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

func App() {
	// connect to database
	db := setup.SQL()

	// test article
	target := news.Article{}
	err := db.Get(&target, "SELECT * FROM NewsArticle WHERE article_id=$1", 2033)

	if err != nil {
		setup.LogCommon(err).Error("Get one article")
	}

	// we need something in this article to process
	if target.Body.IsZero() {
		setup.LogCommon(err).
			WithField("articleID", target.ArticleID).
			Warn("Article Body is empty")
		return
	}

	text := CleanArticle(db, &target)

	if text != nil {
		fmt.Println(len(*text), *text)
	}

	Summary(text)
}

// CleanArticle returns the article text after normaization and removing sentences
// common to multiple articles of the same source (i.e. ads, promotions, boilerplate).
// May return nil.
func CleanArticle(db *sqlx.DB, target *news.Article) *string {
	body := target.Body.ValueOrZero()
	// Clean up string
	text := NormalizeString(&body)
	// Divide into sentences
	targetSentences := Sentences(text)

	// Get articles published around the same time as the target article
	articles := AdjacentArticles(db, target)
	// Only continue if we have a good number of articles to reference
	if len(articles) < setup.EnvToInt("TEXT_ARTICLE_CUTOFF") {
		return nil
	}

	// Iterate through the articles to collect sentences
	checkSentences := []string{}
	for _, each := range articles {
		// sanity check
		if !each.Body.IsZero() {
			// process the sentences
			body := each.Body.ValueOrZero()
			// Clean up string
			text := NormalizeString(&body)
			// Divide into sentences
			newSentences := Sentences(text)
			// save for later
			checkSentences = append(checkSentences, newSentences...)
		}
	}

	// Sanity check we got a reasonable number of sentences
	if len(checkSentences) < setup.EnvToInt("TEXT_ARTICLE_CUTOFF") {
		return nil
	}

	// Find unique sentences
	finalSentences := UniqueSentences(targetSentences, checkSentences)
	// Sanity check we got a reasonable number of sentences
	if len(finalSentences) < 2 {
		return nil
	}

	out := strings.Join(finalSentences, " ")
	return &out
}
