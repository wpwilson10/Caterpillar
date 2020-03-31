/*
	Package redis creates data structures that interact with 
	a redis database.
*/

package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Set implements a Redis set where name is the set key
type Set struct {
	client *redis.Client
	name   string
}

// NewSet creates a Set
func NewSet(client *redis.Client, name string) *Set {
	return &Set{client, name}
}

// Add puts input in the set if it does not already exist, otherwise does nothing.
func (s *Set) Add(input string) {
	err := s.client.SAdd(s.name, input).Err()
	if err != nil {
		setup.LogCommon(err).WithField("input", input).Error("Failed SAdd")
	}
}

// IsMember returns true if input is in the set, false otherwise.
func (s *Set) IsMember(input string) bool {
	out, err := s.client.SIsMember(s.name, input).Result()
	if err != nil {
		setup.LogCommon(err).WithField("input", input).Error("Failed SIsMember")
		return false
	}

	return out
}
