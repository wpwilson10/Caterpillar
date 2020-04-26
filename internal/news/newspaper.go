package news

import (
	"bytes"
	"net/http"
	"os"

	"github.com/gorilla/rpc/v2/json2"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Args contains the request values send to the newspaper3k python library
type Args struct {
	Link string
}

// Newspaper contains the data from the newspaper3k python library
// PubDate will sometimes be an empty string
type Newspaper struct {
	Title     string
	Text      string
	Authors   []string
	Canonical string
	PubDate   string
}

// NewNewspaper calls the newspaper3k python library for the given source and parses the result.
// This will usually be a slow call; good to make async.
// Can return nil if calls failed. Caller should check.
func NewNewspaper(source *Source) *Newspaper {
	setup.LogCommon(nil).
		WithField("Link", source.Link).
		Info("Processing article")

	// address to call for the newspaper3k application
	var url string = os.Getenv("NEWSPAPER_HOST")

	// Args object is what get sent out the RPC call
	args := Args{
		Link: source.Link,
	}

	// the return data
	var result Newspaper

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
		if err == json2.ErrNullResult {
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
		if err == json2.ErrNullResult {
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
