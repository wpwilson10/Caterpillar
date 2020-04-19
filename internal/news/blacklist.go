package news

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// BlackList is a list of hosts that should be excluded from the NewsArticle database.
type BlackList struct {
	blacklist map[string]int
}

// NewBlackList creates a list of hosts that should be excluded from the NewsArticle database.
func NewBlackList() *BlackList {
	// get filepath
	absPath, err := filepath.Abs(os.Getenv("NEWSPAPER_BLACKLIST_FILEPATH"))
	if err != nil {
		setup.LogCommon(err).Fatal("Filepath")
	}

	// Open the file
	csvfile, err := os.Open(absPath)
	if err != nil {
		setup.LogCommon(err).Fatal("Open csv")
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	m := make(map[string]int)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			setup.LogCommon(err).Fatal("Read csv")
		}

		// add to list
		if len(record) > 0 {
			m[record[0]] = 1
		}
	}

	// put map in blacklist
	blacklist := BlackList{blacklist: m}

	return &blacklist
}

// IsBlackListed returns true if the given host is in the blacklist.
// False otherwise.
func (b *BlackList) IsBlackListed(host string) bool {
	_, prs := b.blacklist[host]
	return prs
}
