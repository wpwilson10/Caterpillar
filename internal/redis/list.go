package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// CappedList implements a Redis list with a max number of elements.
type CappedList struct {
	client *redis.Client
	name   string
	size   int64
}

// NewCappedList creates a CappedList
func NewCappedList(client *redis.Client, name string, size int) *CappedList {
	// sanity check then convert to int64
	var s int64
	if size < 1 {
		setup.LogCommon(nil).WithField("size", size).Error("Size must be positive")
		s = 1
	} else {
		s = int64(size)
	}

	return &CappedList{client, name, s}
}

// Add puts the input at the front of the list and removes extra elements from the end.
func (c *CappedList) Add(input string) {
	// add new element
	err := c.client.LPush(c.name, input).Err()
	if err != nil {
		setup.LogCommon(err).WithField("input", input).Error("Failed LPush")
	}
	// remove old elements
	err = c.client.LTrim(c.name, 0, c.size-1).Err()
	if err != nil {
		setup.LogCommon(err).WithField("input", input).Error("Failed LTrim")
	}
}

// List returns the non-null elements in the list
func (c *CappedList) List() []string {
	// get size of list
	size, err := c.client.LLen(c.name).Result()
	if err != nil {
		setup.LogCommon(err).Error("Failed LLen")
	} else if size > c.size {
		// sanity check that the list is not bigger than possible
		setup.LogCommon(nil).WithField("size", size).Error("Size is larger than expected")
	}

	// get the elements
	out, err := c.client.LRange(c.name, 0, size-1).Result()
	if err != nil {
		setup.LogCommon(err).Error("Failed LRange")
	}

	return out

}
