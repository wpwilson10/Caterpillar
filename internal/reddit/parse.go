package reddit

import (
	"fmt"

	"github.com/turnage/graw/reddit"
)

func Test() {
	bot := BotClient()
	GetSubmission2(bot, "/r/politics/comments/hbo4a7/fauci_warns_of_antiscience_bias_being_a_problem/")

	sub := GetSubmission(bot, "/r/politics/comments/hbo4a7/fauci_warns_of_antiscience_bias_being_a_problem/")
	coms := ProcessComments(sub.Replies)
	fmt.Println(len(coms))
}

func GetSubmission2(bot *reddit.Bot, permalink string) *reddit.Post {
	opts := map[string]string{
		"raw_json": "1",
		"limit":    "1000",
		"depth":    "1000",
		// "sort":     "top",
	}

	harvest, err := (*bot).ListingWithParams(permalink, opts)
	if err != nil {
		fmt.Println("GetSub2 error", err)
		return nil
	}

	comments := parseComments(harvest.Posts[0].Replies)

	fmt.Println(len(harvest.Posts[0].Replies), len(comments))

	mQ := newMoreQueue(harvest, 5, 3, 15)
	mQ.MoreChildren(bot)

	return nil
}

// parseComments takes branching comment trees and returns a list of the comments.
// Each comment from geddit contains a tree of comments that are the comments children.
// This travels the trees and adds them to a simple list.
func parseComments(replies []*reddit.Comment) []*reddit.Comment {
	// save the current comments in a list
	out := replies
	// get child comments
	for _, comment := range replies {
		if len(comment.Replies) > 0 {
			// recursively explore child tree
			out = append(out, parseComments(comment.Replies)...)
		}
	}
	return out
}
