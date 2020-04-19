package reddit

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/turnage/graw/reddit"
	"github.com/wpwilson10/caterpillar/internal/redis"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// QueueSubmission represents a submission sent to our queue for later processing.
type QueueSubmission struct {
	CreatedTime time.Time
	Permalink   string
}

// NewQueueSubmission creates a queue object for the given reddit post.
func NewQueueSubmission(p *reddit.Post) *QueueSubmission {
	// convert time
	var y int64 = int64(p.CreatedUTC)
	createdTime := time.Unix(y, 0)

	return &QueueSubmission{
		CreatedTime: createdTime,
		Permalink:   p.Permalink,
	}
}

// Push marshals this object into JSON and adds it to the queue.
func (s *QueueSubmission) Push(queue *redis.Queue) {
	jsonData, err := json.Marshal(s)

	if err != nil {
		setup.LogCommon(err).
			WithField("permaLink", s.Permalink).
			Error("Failed json.Marshal")
	}

	queue.Push(string(jsonData))
}

// PopQueue returns submissions older than 24 hours from the REDDIT_QUEUE redis queue.
func PopQueue(queue *redis.Queue) []QueueSubmission {
	out := []QueueSubmission{}

	// flag tracks whether we have found a submission older than cut off
	flag := true
	for flag {
		// look at most recent submission
		s := queue.Peek()

		if s == nil {
			// nothing returned so quit
			flag = false
		} else {
			// convert to struct
			q := QueueSubmission{}
			err := json.Unmarshal([]byte(*s), &q)
			if err != nil {
				setup.LogCommon(err).
					WithField("queueSubmission", q).
					Error("Failed json.Unmarshal")
			}

			// get lookback time
			lookback, err := strconv.ParseFloat(os.Getenv("REDDIT_LOOKBACK"), 64)
			if err != nil {
				setup.LogCommon(err).
					WithField("lookback", lookback).
					Fatal("Failed REDDIT_LOOKBACK float conversion")
			}

			// check age
			if time.Since(q.CreatedTime).Hours() >= lookback {
				// remove peeked value
				queue.Pop()
				// add to return
				out = append(out, q)
			} else {
				// got a submission newer than REDDIT_LOOKBACK
				flag = false
			}
		}
	}

	return out
}
