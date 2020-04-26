package text

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"

	"github.com/microcosm-cc/bluemonday"
	"github.com/wpwilson10/caterpillar/internal/setup"
	"golang.org/x/text/unicode/norm"
	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"
)

func App() {
	// connect to database
	// db := setup.SQL()

	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("./test/foxnews1.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	str := string(content)

	text := (Clean(&str))
	Sentences(text)

	// run summary
	setup.LogCommon(nil).
		WithField("RunTime", setup.RunTime().String()).
		Info("RunSummary")
}

func Clean(input *string) *string {
	// StrictPolicy strips all HTML elements (and their attributes)
	text := bluemonday.StrictPolicy().Sanitize(*input)
	if len(text) == 0 {
		return nil
	}

	fmt.Println(text)
	fmt.Println("LENGTH", len(text))
	fmt.Println("**********")

	// Unescape remaining HTML
	text = html.UnescapeString(text)

	fmt.Println(text)
	fmt.Println("LENGTH", len(text))
	fmt.Println("-----------")

	// Normalize unicode, see https://blog.golang.org/normalization
	// This may look like it does nothing,
	// but it becomes obvious when you do a string len compare
	text = norm.NFKC.String(text)

	fmt.Println(text)
	fmt.Println("LENGTH", len(text))
	fmt.Println("===========")

	return &text
}

func Sentences(input *string) {

	// Compiling language specific data into a binary file can be accomplished
	// by using `make <lang>` and then loading the `json` data:
	b, _ := data.Asset("data/english.json")

	// load the training data
	training, _ := sentences.LoadTraining(b)

	// create the default sentence tokenizer
	tokenizer := sentences.NewSentenceTokenizer(training)

	sentences := tokenizer.Tokenize(*input)

	for i, s := range sentences {
		fmt.Println(i, s.Text)
	}
}

/*
func Sentences(input *string) []*string {
	// Parse text into pieces
	doc, err := prose.NewDocument(*input)
	if err != nil {
		setup.LogCommon(err).Warn("Failed Prose.NewDocument")
		return nil
	}

	out := []*string{}

	// Iterate over the doc's sentences:
	for i, sent := range doc.Sentences() {
		fmt.Println(i, sent.Text)
		out = append(out, &sent.Text)
	}

	return out
}
*/
