package stocks

import (
	"fmt"
	"time"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// App runs the IEX intraday data colletion application.
func App() {
	// Setup necessary clients
	client := IEXSetup()
	db := setup.SQL()

	// get active listings and latest intraday times from database
	// use russell3000 index to reduce number of calls
	listings := Russell3000Listings(db)
	latestTimes := LatestIntraday(db)

	// Update data for all active listings
	for _, l := range listings {
		// sanity check
		if l.IsEnabled == true {
			// do the work
			data := IEXIntraday(client, l)
			cleanData := SanitizeIntraday(l, data, latestTimes)
			InsertIntraday(db, cleanData)
			// be polite
			time.Sleep(1 * time.Second)
		}
	}
}

// UpdateIndex update the index for listings in the IEX listing table.
func UpdateIndex() {
	// Setup necessary clients
	db := setup.SQL()

	// get all listings from database
	dbListings := AllListings(db)
	// update index values
	freshListings := UpdateIndexTable(db, dbListings)

	// update existing listings
	toAuditListings, updatedListings := ChangedOnIndex(dbListings, freshListings)
	AuditListings(toAuditListings, db)
	UpdateListings(updatedListings, db)

	fmt.Println(len(dbListings), " ", len(freshListings), " ", len(toAuditListings), " ", len(updatedListings))
}

// UpdateListingsDriver checks IEX for new or changed listings and updates database.
func UpdateListingsDriver() {
	// Setup necessary clients
	client := IEXSetup()
	db := setup.SQL()

	// get all listings from database
	dbListings := AllListings(db)
	// new listings from IEX
	freshListings := IEXSymbols(client)

	// update existing listings
	toAuditListings, updatedListings := ChangedOnIEX(dbListings, freshListings)
	AuditListings(toAuditListings, db)
	UpdateListings(updatedListings, db)

	// add new listings
	newListings := NewListings(dbListings, freshListings)
	InsertNewListings(newListings, db)

}
