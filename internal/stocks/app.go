package stocks

import (
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
	listings := ActiveListings(db)
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

	// log summary
	setup.LogCommon(nil).
		WithField("RunTime", setup.RunTime().String()).
		Info("RunSummary")
}

// UpdateActiveDriver update the active status for listings in the IEX listing table.
func UpdateActiveDriver() {
	// Setup necessary clients
	db := setup.SQL()

	// get all listings from database
	dbListings := AllListings(db)
	// update index values
	freshListings := UpdateActive(db, dbListings)

	// update existing listings
	toAuditListings, updatedListings := ChangedActive(dbListings, freshListings)
	AuditListings(toAuditListings, db)
	UpdateListings(updatedListings, db)
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

	// log summary
	setup.LogCommon(nil).
		WithField("NumNewListings", len(newListings)).
		WithField("NumUpdatedListings", len(updatedListings)).
		WithField("RunTime", setup.RunTime().String()).
		Info("RunSummary")
}
