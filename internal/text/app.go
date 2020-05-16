package text

import (
	"fmt"

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

	body := target.Body.ValueOrZero()
	// Do initialize normization
	text := InitialClean(&body)
	// Divide into sentences
	targetSentences := Sentences(text)

	// Get articles published around the same time as the target article
	articles := AdjacentArticles(db, &target)
	// Only continue if we have a good number of articles to reference
	if len(articles) < 10 {
		return
	}

	// Iterate through the articles to collect sentences
	checkSentences := []string{}
	for _, each := range articles {
		// sanity check
		if !each.Body.IsZero() {
			// process the sentences
			body := each.Body.ValueOrZero()
			// Do initialize normization
			text := InitialClean(&body)
			// Divide into sentences
			newSentences := Sentences(text)
			// save for later
			checkSentences = append(checkSentences, newSentences...)
		}
	}

	// Return only unique sentences
	finalSentences := UniqueSentences(targetSentences, checkSentences)

	for i, each := range finalSentences {
		fmt.Println(i, " - ", each)
	}

	fmt.Println(len(finalSentences), len(targetSentences))
}
