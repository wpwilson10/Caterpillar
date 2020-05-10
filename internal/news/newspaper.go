package news

import (
	"context"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wpwilson10/caterpillar/internal/setup"
	"github.com/wpwilson10/caterpillar/protobuf"
)

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
	var host string = os.Getenv("NEWSPAPER_HOST")
	// connect to server
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		setup.LogCommon(err).Error("Failed gRPC Dial")

		return nil
	}
	defer conn.Close()

	// create our client
	client := protobuf.NewNewspaperClient(conn)

	// Make request
	response, err := client.Request(context.Background(),
		&protobuf.NewspaperRequest{Link: source.Link})

	// handle possible failure codes
	if err != nil {
		// get error code
		if e, ok := status.FromError(err); ok {
			// check if this is a known code, don't throw warning
			if e.Code() == codes.Internal {
				return nil
			}
		}

		// else handle as an error
		setup.LogCommon(err).
			WithField("Link", source.Link).
			Warn("Failed gRPC request")

		return nil
	}

	// perform link consistency checks
	if strings.Compare(source.Link, response.GetLink()) != 0 {
		setup.LogCommon(nil).
			WithField("Link", source.Link).
			WithField("response", response.GetLink).
			Error("Links do not match")

		return nil
	}
	// check that we got text to return
	if len(response.GetTitle()) < 3 || len(response.GetText()) < 3 {
		return nil
	}

	// Put into internal format
	return &Newspaper{
		Title:     response.GetTitle(),
		Text:      response.GetText(),
		Authors:   response.GetAuthors(),
		Canonical: response.GetCanonical(),
		PubDate:   response.GetPubdate(),
	}
}
