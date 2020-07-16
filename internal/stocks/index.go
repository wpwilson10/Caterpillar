package stocks

import (
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/jmoiron/sqlx"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// UpdateActive returns an array of listings with their active status values updated from .csv files.
// Does not change the input listings.
// The .csv filepaths are configured in the enivonment's .env file
func UpdateActive(db *sqlx.DB, listings []Listing) []Listing {
	// keep track of what changed
	updatedListings := []Listing{}

	// Get data from .csv file
	var filepath string = os.Getenv("ACTIVE_FILE")
	var active []*Membership = CSVtoMembership(filepath)

	// sanity check that we got something
	if len(active) > 0 {
		setup.LogCommon(nil).
			WithField("len(active)", len(active)).
			Error("Member length check")
	}

	// Standardize data and put in hashmap for easy lookup
	mapActive := make(map[string]int)
	for i, s := range active {
		var sym string = strings.TrimSpace(s.Symbol)
		sym = strings.ToUpper(sym)
		mapActive[sym] = i
	}

	// do updates
	for _, l := range listings {
		// find a matching listing
		var sym string = strings.ToUpper(l.Symbol)
		_, isActive := mapActive[sym]

		// copy listing
		temp := l
		// update listing values to match whether in map
		temp.IsActive = isActive
		updatedListings = append(updatedListings, temp)
	}

	return updatedListings
}

// ChangedActive compares listings matching on listingID to find which have changed active status.
// Returns the original listings that should go to audit, and the updated listings.
func ChangedActive(db []Listing, fresh []Listing) (original []Listing, updated []Listing) {

	// hashmap to allow O(1) lookup instead of looping
	m := make(map[int64]int)
	// put db listings in map
	for i, s := range db {
		m[s.ListingID] = i
	}

	// check whether the fresh listings are in the database
	for _, f := range fresh {
		// Find a matching listing
		i, keyExists := m[f.ListingID]

		var flag bool // true if listings do not match

		if keyExists {
			flag = false
			// get the listing
			s := db[i]
			// copy original listing so that fields we don't care about in this function don't change
			temp := s
			// Check whether index values we care about have changed.
			// If they have changed, put the fresh value in temp
			if s.IsActive != f.IsActive {
				flag = true
				temp.IsActive = f.IsActive
			}

			// save the differing listings
			if flag {
				original = append(original, s)
				updated = append(updated, temp)

				setup.LogCommon(nil).
					WithField("listingID", s.ListingID).
					WithField("iexID", s.IexID).
					WithField("symbol", s.Symbol).
					WithField("symbolFresh", f.Symbol).
					WithField("listingIDFresh", f.ListingID).
					Info("Changed Active Status")
			}
		}
	}

	// sanity check that same number of listings are returned
	if len(original) != len(updated) {
		setup.LogCommon(nil).
			WithField("len(changedDb)", len(original)).
			WithField("len(changedFresh)", len(updated)).
			Error("Return arrays size does not match")
	}

	return original, updated
}

// UpdateIndexTable returns an array of listings with their index values updated from .csv files.
// Does not change the input listings. Currently only checks SP500 and Russell3000.
// The .csv filepaths are configured in the enivonment's .env file
func UpdateIndexTable(db *sqlx.DB, listings []Listing) []Listing {
	// keep track of what changed
	updatedListings := []Listing{}

	// Get data from .csv file
	var filepathSP500 string = os.Getenv("SP500_FILE")
	var SP500 []*Membership = CSVtoMembership(filepathSP500)

	var filepathRussell3000 string = os.Getenv("RUSSELL3000_FILE")
	var Russell3000 []*Membership = CSVtoMembership(filepathRussell3000)

	// sanity check that sizes are what we expect and russell3000 will be larger
	if len(Russell3000) < len(SP500) {
		setup.LogCommon(nil).
			WithField("len(Russell3000)", len(Russell3000)).
			WithField("len(SP500)", len(SP500)).
			Error("UpdateIndexTable length check")
	}

	// Standardize data and put in hashmap for easy lookup
	mapSP500 := make(map[string]int)
	for i, s := range SP500 {
		var sym string = strings.TrimSpace(s.Symbol)
		sym = strings.ToUpper(sym)
		mapSP500[sym] = i
	}

	mapRussell3000 := make(map[string]int)
	for i, s := range Russell3000 {
		var sym string = strings.TrimSpace(s.Symbol)
		sym = strings.ToUpper(sym)
		mapRussell3000[sym] = i
	}

	// do updates
	for _, l := range listings {
		// find a matching listing
		var sym string = strings.ToUpper(l.Symbol)

		_, inSP500 := mapSP500[sym]
		_, inRussell3000 := mapRussell3000[sym]

		// copy listing
		temp := l

		// update listing values to match whether in map
		temp.IsSP500 = inSP500
		temp.IsRussell3000 = inRussell3000

		updatedListings = append(updatedListings, temp)
	}

	return updatedListings
}

// ChangedOnIndex compares listings matching on listingID to find which have changed index fields.
// Returns the original listings that should go to audit, and the updated listings.
func ChangedOnIndex(db []Listing, fresh []Listing) (original []Listing, updated []Listing) {

	// hashmap to allow O(1) lookup instead of looping
	m := make(map[int64]int)
	// put db listings in map
	for i, s := range db {
		m[s.ListingID] = i
	}

	// check whether the fresh listings are in the database
	for _, f := range fresh {
		// Find a matching listing
		i, keyExists := m[f.ListingID]

		var flag bool // true if listings do not match

		if keyExists {
			flag = false

			// get the listing
			s := db[i]

			// copy original listing so that fields we don't care about in this function don't change
			temp := s

			// Check whether index values we care about have changed.
			// If they have changed, put the fresh value in temp
			if s.IsSP500 != f.IsSP500 {
				flag = true
				temp.IsSP500 = f.IsSP500
			}
			if s.IsRussell3000 != f.IsRussell3000 {
				flag = true
				temp.IsRussell3000 = f.IsRussell3000
			}

			// save the differing listings
			if flag {
				original = append(original, s)
				updated = append(updated, temp)

				setup.LogCommon(nil).
					WithField("listingID", s.ListingID).
					WithField("iexID", s.IexID).
					WithField("symbol", s.Symbol).
					WithField("symbolFresh", f.Symbol).
					WithField("listingIDFresh", f.ListingID).
					Info("Changed Listing Index")
			}
		}
	}

	// sanity check that same number of listings are returned
	if len(original) != len(updated) {
		setup.LogCommon(nil).
			WithField("len(changedDb)", len(original)).
			WithField("len(changedFresh)", len(updated)).
			Error("Return arrays size does not match")
	}

	return original, updated
}

// Membership is used to create a list of listings for an index.
// For unmarshalling to work, it needs a column labeled 'Symbol'
type Membership struct {
	Symbol  string `csv:"Symbol"`
	NotUsed string `csv:"-"`
}

// CSVtoMembership creates an array of Membership structs from the file at filepath.
// See Membership struct for what fields the program is expectin on the CSV.
func CSVtoMembership(filepath string) []*Membership {
	// Get the file
	clientsFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		setup.LogCommon(err).
			WithField("filepath", filepath).
			Error("Index open file")
	}
	defer clientsFile.Close()

	// Read from file to array
	index := []*Membership{}
	// Load clients from file
	if err := gocsv.UnmarshalFile(clientsFile, &index); err != nil {
		setup.LogCommon(err).
			WithField("filepath", filepath).
			Error("Index unmarshal file")
	}

	// Check that we got something
	if len(index) <= 0 {
		setup.LogCommon(nil).
			WithField("filepath", filepath).
			Error("Index empty membership list")
	}

	return index
}
