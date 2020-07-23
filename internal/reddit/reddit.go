package reddit

import (
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/turnage/graw/reddit"

	"github.com/wpwilson10/caterpillar/internal/news"
	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Driver contains the main application logic for adding submissions and comments to the database.
func Driver(db *sqlx.DB, bot *reddit.Bot, wg *sync.WaitGroup, q QueueSubmission, articleSet *redis.Set, blacklist *news.BlackList) {
	// async call
	defer wg.Done()

	// Get updated submission information
	harvest := GetSubmission(bot, q.Permalink)

	// sanity check that we got a single post
	if len(harvest.Posts) != 1 {
		setup.LogCommon(nil).
			WithField("permalink", q.Permalink).
			Error("More than one post returned")
		return
	}
	submission := harvest.Posts[0]

	// if we got a real not-deleted submission, and it was commented or scored enough
	if checkSubmission(submission) {
		// Put submission in database, returns the row ID
		sID := InsertSubmission(db, submission)
		// Transform comments from tree to list
		commentList := ParseComments(submission.Replies)

		// Get more comments if we can assume it will be worthwhile
		if checkGetMores(submission, harvest, commentList) {
			// Get more comments
			mQ := NewMoreQueue(harvest, 6, 2, 20)
			moreComments := mQ.MoreChildren()
			commentList = append(commentList, moreComments...)
		}

		// Add comments to database
		InsertComments(db, commentList, sID)

		// only process links that go externally
		if !(submission.IsRedditMediaDomain || submission.IsSelf) {
			// Handle getting and linking submission to a news article
			news.RedditNewsDriver(db, articleSet, blacklist, submission, sID)
		}
	}
}

// checkSubmission returns true if we got a real not-deleted submission,
// and it was commented + scored at least REDDIT_SCORE_CUTOFF times
func checkSubmission(submission *reddit.Post) bool {
	scoreCutoff := setup.EnvToInt("REDDIT_SCORE_CUTOFF")
	return submission != nil &&
		!submission.Deleted &&
		(submission.NumComments+submission.Score) >= int32(scoreCutoff)
}

// checkGetMores returns true if we should get more comments,
// based on arbitrary cutoff assumptions
func checkGetMores(submission *reddit.Post, harvest *reddit.Harvest, commentList []*reddit.Comment) bool {
	c1 := false
	if submission.More != nil {
		// if the more query will get several comments then do it
		c1 = submission.More != nil && len(submission.More.Children) >= 10
	}

	// if there are many comments missing then do it
	c2 := submission.NumComments > 20 && (float32(len(commentList))/float32(submission.NumComments) < 0.667)

	return c1 || c2
}

// GetSubmission returns a submission harvest based on it's permalink.
// May return nil in case the submission was not found (i.e. deleted)
func GetSubmission(bot *reddit.Bot, permalink string) *reddit.Harvest {
	// use graw to get submission content
	opts := map[string]string{
		"raw_json": "1",
		"limit":    "1000",
		"depth":    "1000",
		// "sort":     "top",
	}
	harvest, err := (*bot).ListingWithParams(permalink, opts)

	if err == reddit.BusyErr || err == reddit.RateLimitErr {
		// reddit is busy, wait and try again
		setup.LogCommon(err).
			WithField("permalink", permalink).
			Info("Recoverable error from reddit")

		time.Sleep(5 * time.Second)
		return GetSubmission(bot, permalink)
	} else if err == reddit.ThreadDoesNotExistErr {
		// don't log if nothing found
		return nil
	} else if err != nil {
		setup.LogCommon(err).
			WithField("permalink", permalink).
			Warn("Failed bot.Listing")
	}

	return &harvest
}

// InsertSubmission puts a submission into the RedditSubmission database table.
// Returns the ID given to the submission by the database.
func InsertSubmission(db *sqlx.DB, submission *reddit.Post) int64 {
	// Setup
	// Use positional bindvars to not have to recreate struct
	var insertStmt string = `INSERT INTO RedditSubmission (
								submission_id, -- DEFAULT
								data_entry_time, -- Now()
								reddit_id, -- $1
								title, -- $2
								url, -- $3
								permalink, -- $4
								created_time, -- $5
								user_name, -- $6
								subreddit_name, -- $7
								subreddit_id, -- $8
								selftext, -- $9
								selftext_html, -- $10
								num_comments, -- $11
								score, -- $12
								up_votes, -- $13
								down_votes, -- $14
								is_nsfw, -- $15
								is_self  -- $16
								)`

	var valueStmt string = `VALUES (DEFAULT, Now(), $1, $2, $3, $4, $5, $6, $7,
									$8, $9, $10, $11, $12, $13, $14, $15, $16)`

	var returnStmt string = "RETURNING submission_id;"
	var fullStmt string = insertStmt + " " + valueStmt + " " + returnStmt

	// convert time
	var y int64 = int64(submission.CreatedUTC)
	cTime := time.Unix(y, 0)

	// for the return
	var id int64

	// Use this hacky setup because libpq is stupid
	// See https://github.com/jmoiron/sqlx/issues/154
	err := db.QueryRow(fullStmt,
		submission.ID,           // $1
		submission.Title,        // $2
		submission.URL,          // $3
		submission.Permalink,    // $4
		cTime,                   // $5
		submission.Author,       // $6
		submission.Subreddit,    // $7
		submission.SubredditID,  // $8
		submission.SelfText,     // $9
		submission.SelfTextHTML, // $10
		submission.NumComments,  // $11
		submission.Score,        // $12
		submission.Ups,          // $13
		submission.Downs,        // $14
		submission.NSFW,         // $15
		submission.IsSelf).      // $16
		Scan(&id)

	if err != nil {
		setup.LogCommon(err).
			WithField("redditID", submission.ID).
			WithField("permalink", submission.Permalink).
			Error("Failed execute statement")
	}

	setup.LogCommon(nil).
		WithField("redditID", submission.ID).
		WithField("permalink", submission.Permalink).
		Info("Inserting reddit submission")

	return id
}

// InsertComments puts a list of comments into the RedditComment database table.
func InsertComments(db *sqlx.DB, comments []*reddit.Comment, sID int64) {
	// Setup
	// Use positional bindvars to not have to recreate struct
	var insertStmt string = `INSERT INTO RedditComment (
								comment_id, 		-- DEFAULT
								data_entry_time, 	-- Now()
								submission_id,		-- $1
								reddit_id,			-- $2
								parent_id,			-- $3
								created_time,		-- $4
								user_name, 			-- $5
								body,				-- $6
								body_html, 			-- $7
								up_votes, 			-- $8
								down_votes, 		-- $9
								is_deleted			-- $10
								)`

	var valueStmt string = "VALUES (DEFAULT, Now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	var fullStmt string = insertStmt + " " + valueStmt

	// start a transaction
	tx, err := db.Beginx()
	if err != nil {
		log.WithError(err).Warn("InsertComments start transaction")
	}

	// process each comment
	for _, comment := range comments {
		// convert time
		var y int64 = int64(comment.CreatedUTC)
		cTime := time.Unix(y, 0)

		// Insert into queue
		_, err = tx.Exec(fullStmt,
			sID,              // $1
			comment.ID,       // $2
			comment.ParentID, // $3
			cTime,            // $4
			comment.Author,   // $5
			comment.Body,     // $6
			comment.BodyHTML, // $7
			comment.Ups,      // $8
			comment.Downs,    // $9
			comment.Deleted)  // $10

		if err != nil {
			log.
				WithField("redditID", comment.ID).
				WithField("submissionID", sID).
				WithError(err).
				Error("InsertComments execute statement")
		}
	}

	// send it
	err = tx.Commit()
	if err != nil {
		setup.LogCommon(err).Warn("InsertComments commiting transaction")
	}
}

// ParseComments takes branching comment trees and returns a list of the comments.
// Each comment from geddit contains a tree of comments that are the comments children.
// This travels the trees and adds them to a simple list.
func ParseComments(replies []*reddit.Comment) []*reddit.Comment {
	// save the current comments in a list
	out := replies
	// get child comments
	for _, comment := range replies {
		if len(comment.Replies) > 0 {
			// recursively explore child tree
			out = append(out, ParseComments(comment.Replies)...)
		}
	}
	return out
}
