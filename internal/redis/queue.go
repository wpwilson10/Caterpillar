package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Queue implements a Redis list as a queue where name is the list key
type Queue struct {
	client *redis.Client
	name   string
}

// NewQueue creates a Queue
func NewQueue(client *redis.Client, name string) *Queue {
	return &Queue{client, name}
}

// Push adds the input to the end of the queue
func (q *Queue) Push(input string) {
	err := q.client.RPush(q.name, input).Err()
	if err != nil {
		setup.LogCommon(err).WithField("input", input).Error("Failed RPush")
	}
}

// Pop returns the first value and removes it from the queue
func (q *Queue) Pop() string {
	out, err := q.client.LPop(q.name).Result()
	if err != nil {
		setup.LogCommon(err).Error("Failed LPop")
	}

	return out
}

// Peek returns the first value but does not remove it from the queue
func (q *Queue) Peek() string {
	out, err := q.client.LRange(q.name, 0, 0).Result()
	if err != nil {
		setup.LogCommon(err).Error("Failed LRange")
	} else if len(out) != 1 {
		setup.LogCommon(err).Error("Unexpected LRange return size")

		return ""
	}

	return out[0]
}
