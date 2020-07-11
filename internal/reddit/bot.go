package reddit

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/turnage/graw/reddit"

	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Simple struct so we can describe the custom function handles
type redditBot struct {
	bot   reddit.Bot
	queue *redis.Queue
}

// Post handler for our custom reddit bot
func (r *redditBot) Post(p *reddit.Post) error {
	// turn into queue submission
	submission := NewQueueSubmission(p)
	// add to queue
	submission.Push(r.queue)

	return nil
}

// BotClient returns a graw client using my login config from a .env.
func BotClient() *reddit.Bot {
	// Bot account and login info
	botCfg := reddit.BotConfig{
		Agent: os.Getenv("REDDIT_USER_AGENT"),
		App: reddit.App{
			ID:       os.Getenv("REDDIT_CLIENT_ID"),
			Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
			Username: os.Getenv("REDDIT_USER"),
			Password: os.Getenv("REDDIT_PASSWORD"),
		},
		// reddit api has 60 calls/minute limit
		// https://github.com/reddit-archive/reddit/wiki/API#rules
		Rate: time.Second,
	}

	bot, err := reddit.NewBot(botCfg)
	if err != nil {
		setup.LogCommon(err).Fatal("SetupBot NewBot")
	}

	return &bot
}

func readSubredditList() []string {
	// get subreddit list filepath
	absPath, err := filepath.Abs(os.Getenv("REDDIT_LIST_FILEPATH"))
	if err != nil {
		setup.LogCommon(err).Fatal("Filepath")
	}

	// Open the file
	csvfile, err := os.Open(absPath)
	if err != nil {
		setup.LogCommon(err).Fatal("Open csv")
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	var subreddits []string

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			setup.LogCommon(err).Fatal("Read csv")
		}

		// add to list
		if len(record) > 0 {
			subreddits = append(subreddits, record[0])
		}
	}

	return subreddits
}
