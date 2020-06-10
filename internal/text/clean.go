package text

import (
	"html"
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"

	"github.com/agnivade/levenshtein"
	"github.com/microcosm-cc/bluemonday"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

// UniqueSentences returns sentences from the target list that are not in the check list.
func UniqueSentences(targetSentences []string, checkSentences []string) []string {
	out := []string{}
	match := false
	// compare string distance across all strings
	for _, s1 := range targetSentences {
		match = false
		for _, s2 := range checkSentences {
			// remove sentences that are not unique
			distance := levenshtein.ComputeDistance(s1, s2)
			// check that sentences are not mostly the same, arbitrary cutoff
			if distance < 10 {
				match = true
				break
			}
		}
		// if no match sentences found, it is unique so return
		if !match {
			out = append(out, s1)
		}
	}

	return out
}

// NormalizeString performs basic string unicode normalization and HTML escaping.
func NormalizeString(input *string) *string {
	// strip invalid characters
	text := strings.ToValidUTF8(*input, "")

	// StrictPolicy strips all HTML elements (and their attributes)
	text = bluemonday.StrictPolicy().Sanitize(*input)
	if len(text) == 0 {
		return nil
	}

	// Unescape remaining HTML
	text = html.UnescapeString(text)

	// Normalize unicode, see https://blog.golang.org/normalization
	// This may look like it does nothing,
	// but it becomes obvious when you do a string len compare
	text = norm.NFKC.String(text)

	return &text
}

// RemoveEmptySentences removes strings that contain only puncutation, whitespace, or control characters.
func RemoveEmptySentences(sentences []string) []string {
	// finds strings that are not just puncuation, whitespace, and control characters
	// Unicode operators, so should work for any language.
	re1, err := regexp.Compile(`\P{P}\P{Z}\P{C}`)
	if err != nil {
		setup.LogCommon(err).Error("Failed Regex compile")
	}

	// save sentences that are not "empty"
	out := []string{}
	for _, each := range sentences {
		// returns nil if no match, i.e. only whitespace and punctuation found
		if re1.FindStringIndex(each) != nil {
			out = append(out, each)
		}
	}

	return out
}
