package text

import (
	"context"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wpwilson10/caterpillar/internal/setup"
	"github.com/wpwilson10/caterpillar/protobuf"
)

// Sentences parses the given text into individual sentences and returns all non-empty strings.
// External call so it can be slow.
func Sentences(text *string) []string {
	// address to call for the text application
	var host string = os.Getenv("TEXT_HOST")
	// connect to server
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		setup.LogCommon(err).Error("Failed gRPC Dial")

		return nil
	}
	defer conn.Close()

	// create our client
	client := protobuf.NewTextClient(conn)

	// Make request
	response, err := client.Sentences(context.Background(),
		&protobuf.TextRequest{Text: *text})

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
			Warn("Failed gRPC request")

		return nil
	}

	// do gRPC call
	sentences := response.GetSentences()
	// return only non-empty strings
	return RemoveEmptySentences(sentences)
}
