package reddit

import (
	"fmt"
	"strings"

	"github.com/Workiva/go-datastructures/queue"
	"github.com/turnage/graw/reddit"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// wrapper that implements Compare so we can use the priority queue
type moreItem struct {
	more *reddit.More
}

func (m moreItem) Compare(other queue.Item) int {
	// make sure this is moreItem
	mI, ok := other.(moreItem)
	// can't compare unlike types
	if !ok {
		return 0
	} else if len(m.more.Children) > len(mI.more.Children) {
		return -1
	} else if len(m.more.Children) < len(mI.more.Children) {
		return 1
	}
	return 0
}

// MoreQueue handles retrieving more comment links from reddit submissions
type MoreQueue struct {
	pq           *queue.PriorityQueue // use priority queue to help focus on largest more queries
	maxDepth     int                  // limit depth of comment tree
	maxCalls     int                  // limit number of calls to gather more comments
	numCalls     int                  // how many calls have actually been done
	cutoff       int                  // number of comments a more must have to consider for further parsing
	reaperParams map[string]string    // common parameters for queries
	comments     []*reddit.Comment    // save all comments for output
}

// get the next more item from the queue
func (mQ *MoreQueue) pop() *reddit.More {
	if mQ.pq.Empty() {
		return nil
	}
	// ignore err because we just checked there are values
	item, _ := mQ.pq.Get(1)
	m := item[0]
	// assert these are in fact moreItems
	mI, ok := m.(moreItem)
	if ok {
		return mI.more
	}
	return nil
}

// checkConditions returns true if we should query reddit for the given more
func (mQ *MoreQueue) checkConditions(more *reddit.More) bool {
	return more.Depth <= mQ.maxDepth &&
		mQ.cutoff <= len(more.Children) &&
		mQ.numCalls <= mQ.maxCalls
}

// NewMoreQueue creates a queue for harvesting mores
func NewMoreQueue(harvest *reddit.Harvest, maxDepth int, cutoff int, maxCalls int) *MoreQueue {
	mQ := MoreQueue{
		pq:       queue.NewPriorityQueue(1, false),
		maxDepth: maxDepth,
		cutoff:   cutoff,
		maxCalls: maxCalls,
		numCalls: 0,
	}

	var link string // we will want to remember the link name
	// if this is a thread call, add the thread's objects
	if len(harvest.Posts) == 1 {
		mQ.addToQueue(harvest.Posts[0].Replies, []*reddit.More{harvest.Posts[0].More})
		// Posts have a discrete name
		link = harvest.Posts[0].Name
	}

	// then handle any mores or comments that are in the harvest
	mQ.addToQueue(harvest.Comments, harvest.Mores)

	// now figure out the link name if we didn't get it from the post
	if link == "" {
		// there really should be some comments to get this far
		if len(harvest.Comments) > 0 {
			link = harvest.Comments[0].ParentID
			// but maybe I will use this in the future
		} else if len(harvest.Mores) > 0 {
			link = harvest.Mores[0].ParentID
		}
	}

	// setup common parameters with that link
	mQ.reaperParams = map[string]string{
		"api_type": "json",
		"link_id":  link,
	}

	return &mQ
}

// MoreChildren repeated queries /api/morechildren using this queue to get comments from a reddit thread.
func (mQ *MoreQueue) MoreChildren() []*reddit.Comment {
	// create client that performs queries
	bot := MoreBotClient()

	// while we still have mores in the queue
	for !mQ.pq.Empty() {
		// get the next more with most number of children
		more := mQ.pop()
		// check our parameters
		if mQ.checkConditions(more) {
			// good so process
			mQ.numCalls = mQ.numCalls + 1
			// put children in comment delimited list
			mQ.reaperParams["children"] = strings.Join(more.Children, ",")
			// query reddit directly
			harvest, err := (*bot).ListingWithParams("/api/morechildren", mQ.reaperParams)
			if err != nil {
				setup.LogCommon(err).
					WithField("link", mQ.reaperParams["link_id"]).
					Error("Failed to query morechildren")
			}
			mQ.addToQueue(harvest.Comments, harvest.Mores)
		}
	}
	fmt.Println("MoreChildren - Link:", mQ.reaperParams["link_id"], "Number of Calls:", mQ.numCalls, "Comments gathered:", len(mQ.comments))

	return mQ.comments
}

// addToQueue takes outputs from a harvest and adds them to the moreQueue struct
func (mQ *MoreQueue) addToQueue(comments []*reddit.Comment, mores []*reddit.More) {
	// travel comment tree to get all comments
	allComments := ParseComments(comments)
	// save out comments
	mQ.comments = append(mQ.comments, allComments...)
	// check for mores associated with each comment
	for _, c := range allComments {
		if c.More != nil {
			mQ.pq.Put(moreItem{more: c.More})
		}
	}
	// add any directly found mores
	for _, m := range mores {
		mQ.pq.Put(moreItem{more: m})
	}
}
