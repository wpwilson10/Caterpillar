package news

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/rpc/v2/json2"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Newspaper handles asyncronous calls to the newspaper3k library and stores the returned articles.
type Newspaper struct {
	Articles []*RawArticle
	Sources  []*Source
	mutex    sync.Mutex
}

// NewNewspaper creates a Newspaper to handle calling and saving data from the Newspaper3k library.
func NewNewspaper() *Newspaper {
	// mutex will just exist because reasons
	return &Newspaper{Articles: []*RawArticle{}}
}

// Process extracts the article for the given source and adds it to the Articles list.
func (n *Newspaper) Process(source *Source, wg *sync.WaitGroup) {
	defer wg.Done()
	// perform the slow aync call
	article := NewRawArticle(source)

	// only add data if we got returned text
	if article != nil && len(article.Title) > 1 {
		// lock list then add return
		n.mutex.Lock()
		defer n.mutex.Unlock()
		n.Articles = append(n.Articles, article)
		n.Sources = append(n.Sources, source)
	}
}

// Args contains the request values send to the newspaper3k python library
type Args struct {
	Link string
}

// RawArticle contains the data from the newspaper3k python library
// PubDate will sometimes be an empty string
type RawArticle struct {
	Title     string
	Text      string
	Authors   []string
	Canonical string
	PubDate   string
}

// NewRawArticle calls the newspaper3k python library for the given source and parses the result.
// This will usually be a slow call; good to make async
// Can return nil if calls failed. Caller should check.
func NewRawArticle(source *Source) *RawArticle {
	// address to call for the newspaper3k application
	var url string = os.Getenv("NEWSPAPER_HOST")

	// Args object is what get sent out the RPC call
	args := Args{
		Link: source.Link,
	}

	// the return data
	var result RawArticle

	// Calls the extractNewspaper method on the reciever
	message, err := json2.EncodeClientRequest("extractNewspaper", args)
	if err != nil {
		setup.LogCommon(err).
			WithField("link", source.Link).
			Warn("Failed EncodeClientRequest")
		// stop to avoid null pointer issues
		return nil
	}

	// Setup an http request call
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		setup.LogCommon(err).
			WithField("link", source.Link).
			Warn("Failed http.NewRequest")
		// stop to avoid null pointer issues
		return nil
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the json rpc call using http
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// ignore common error where nothing is returned
		var nullErr = errors.New("result is null")
		if err == nullErr {
			return nil
		}

		setup.LogCommon(err).
			WithField("link", source.Link).
			Warn("Failed client.Do")
		// stop to avoid null pointer issues
		return nil
	}

	// Process the response into the ParsedNewspaper result type
	defer resp.Body.Close()
	err = json2.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		// ignore common error where nothing is returned
		var nullErr = errors.New("result is null")
		if err == nullErr {
			return nil
		}

		setup.LogCommon(err).
			WithField("link", source.Link).
			Warn("Failed DecodeClientResponse")
		// stop to avoid null pointer issues
		return nil
	}

	// everything is good
	return &result
}
