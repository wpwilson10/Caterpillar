package reddit

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/turnage/graw"
	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// App takes queued reddit submissions and gets most recent data to add to database
// Know issues - GetCommments cannot pull all comments for large threads. Limited by API
func App() {
	db := setup.SQL()
	queue := redis.NewQueue(setup.Redis(), os.Getenv("REDDIT_QUEUE"))
	bot := BotClient()

	// get submissions to process
	submissions := PopQueue(queue)

	// for tracking async calls
	var wg sync.WaitGroup

	// process each entry from the submission queue
	for _, s := range submissions {
		fmt.Println(s.Permalink)

		wg.Add(1)
		go Driver(db, bot, &wg, s)
		// reddit api has 60 calls/minute limit, and each run takes two calls
		time.Sleep(2 * time.Second)
	}

	// block until all done
	wg.Wait()
	// log summary
	setup.LogCommon(nil).
		WithField("NumQueued", len(submissions)).
		WithField("RunTime", setup.RunTime().String()).
		Info("RunSummary")
}

// BotApp creates and runs a bot that saves new submissions to our datebase queue.
// Running will block and run indefinitely.
func BotApp() {
	setup.LogCommon(nil).Info("Starting RedditBot")
	// Setup client
	bot := BotClient()
	// connect to queue
	queue := redis.NewQueue(setup.Redis(), os.Getenv("REDDIT_QUEUE"))

	// point bot to my struct with its handles
	handler := &redditBot{bot: *bot, queue: queue}

	// List of subreddits
	subreddits := readSubredditList()
	// Create a configuration specifying what event sources on Reddit graw
	// should connect to the bot.
	subredditCfg := graw.Config{Subreddits: subreddits}

	// Start up
	_, wait, err := graw.Run(handler, *bot, subredditCfg)
	if err != nil {
		setup.LogCommon(err).Fatal("SetupBot Run")
	}

	// block so the bot will announce (ideally) forever.
	err = wait()
	if err != nil {
		setup.LogCommon(err).Fatal("SetupBot Wait")
	}
}
