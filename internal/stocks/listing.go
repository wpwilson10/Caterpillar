package stocks

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Listing contains information about each stock listing.
// Most of the data comes from IEX
type Listing struct {
	ListingID     int64     `db:"listing_id"`
	UpdateTime    time.Time `db:"update_time"`
	IsEnabled     bool      `db:"is_enabled"`
	Symbol        string    `db:"symbol"`
	Name          string    `db:"name"`
	IexID         string    `db:"iex_id"`
	Type          string    `db:"type"`
	Region        string    `db:"region"`
	Currency      string    `db:"currency"`
	Exchange      string    `db:"exchange"`
	IsSP500       bool      `db:"is_sp500"`
	IsRussell3000 bool      `db:"is_russell3000"`
}

// InsertNewListings inserts all information from IEX as a new entry in the Listing table.
// It does not check if the listing already exists. Defaults index values to true.
func InsertNewListings(listings []Listing, db *sqlx.DB) {

	var insertStmt string = "INSERT INTO Listing (listing_id, update_time, is_enabled, symbol, name, iex_id, type, region, currency, exchange)"
	var valueStmt string = "VALUES (DEFAULT, Now(), :is_enabled, :symbol, :name, :iex_id, :type, :region, :currency, :exchange)"
	var wholeStmt string = insertStmt + " " + valueStmt

	// start a transaction
	tx, err := db.Beginx()
	if err != nil {
		setup.LogCommon(err).Warn("InsertNewListing start transaction")
	}

	for _, s := range listings {
		// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
		_, err := tx.NamedExec(wholeStmt, &s)
		if err != nil {
			setup.LogCommon(err).Warn("InsertNewListing inserting rows")
		}
	}

	err = tx.Commit()
	if err != nil {
		setup.LogCommon(err).Warn("InsertNewListing commiting transaction")
	}
}

// UpdateListings takes the freshListings output of ChangedOnIEX and makes listings updates.
// The values from the fresh array are used to update the listing in Listing.
func UpdateListings(fresh []Listing, db *sqlx.DB) {
	// setup statements
	var updateStmt string = "UPDATE Listing"
	var updateSetStmt1 string = " SET update_time = Now(), is_enabled = :is_enabled, symbol = :symbol, name = :name,"
	var updateSetStmt2 string = " type = :type, region = :region, currency = :currency, exchange = :exchange,"
	var updateSetIndex string = " is_sp500 = :is_sp500, is_russell3000 = :is_russell3000"
	var updateWhereStmt string = " WHERE listing_id = :listing_id"
	updateStmt = updateStmt + updateSetStmt1 + updateSetStmt2 + updateSetIndex + updateWhereStmt

	// start a transaction
	tx, err := db.Beginx()
	if err != nil {
		setup.LogCommon(err).Warn("UpdateListings start transaction")
	}

	// run on all listings
	for _, f := range fresh {
		_, err := tx.NamedExec(updateStmt, &f)
		if err != nil {
			setup.LogCommon(err).Warn("UpdateListings updating audit row")
		}
	}

	// send it
	err = tx.Commit()
	if err != nil {
		setup.LogCommon(err).Warn("UpdateListings commiting transaction")
	}

}

// AuditListings takes the output currentListing of ChangedOnIEX and inserts them as audit values.
// The values from the current array are put into the AuditListing table.
func AuditListings(current []Listing, db *sqlx.DB) {
	var insertStmt string = "INSERT INTO AuditListing (audit_id, listing_id, update_time,"
	var insertIEX string = " is_enabled, symbol, name, iex_id, type, region, currency, exchange,"
	var insertIndex string = " is_sp500, is_russell3000)"
	var valueStmt string = "VALUES (DEFAULT, :listing_id, :update_time,"
	var valueIEX string = " :is_enabled, :symbol, :name, :iex_id, :type, :region, :currency, :exchange,"
	var valueIndex string = " :is_sp500, :is_russell3000)"
	var auditStmt string = insertStmt + insertIEX + insertIndex + " " + valueStmt + valueIEX + valueIndex

	// start a transaction
	tx, err := db.Beginx()
	if err != nil {
		setup.LogCommon(err).Warn("AuditListings start transaction")
	}

	// run on all listings
	for _, c := range current {
		// insert audit
		_, err = tx.NamedExec(auditStmt, &c)
		if err != nil {
			setup.LogCommon(err).Warn("AuditListings inserting audit row")
		}
	}

	// send it
	err = tx.Commit()
	if err != nil {
		setup.LogCommon(err).Warn("AuditListings commiting transaction")
	}
}

// ActiveListings returns all entries from the Listing database table that are enabled.
// The returned array is sorted in ascending order by listing ID
func ActiveListings(db *sqlx.DB) []Listing {

	var selectStmt string = "SELECT * FROM Listing WHERE is_enabled='true' ORDER BY listing_id ASC"
	listings := []Listing{}

	// if you have null fields and use SELECT *, you must use sql.Null* in your struct
	err := db.Select(&listings, selectStmt)

	if err != nil {
		setup.LogCommon(err).Error("Select active listings")
	}

	return listings
}

// AllListings returns all entries from the Listing database table.
// The returned array is sorted in ascending order by listing ID
func AllListings(db *sqlx.DB) []Listing {

	var selectStmt string = "SELECT * FROM Listing ORDER BY listing_id ASC"
	listings := []Listing{}

	// if you have null fields and use SELECT *, you must use sql.Null* in your struct
	err := db.Select(&listings, selectStmt)

	if err != nil {
		setup.LogCommon(err).Error("Select all listings")
	}

	return listings
}

// Russell3000Listings returns all Listing entries that are enabled and on the Russell 3000 index.
// The returned array is sorted in ascending order by listing ID
func Russell3000Listings(db *sqlx.DB) []Listing {

	var selectStmt string = "SELECT * FROM Listing WHERE is_russell3000='true'AND is_enabled='true' ORDER BY listing_id ASC"
	listings := []Listing{}

	// if you have null fields and use SELECT *, you must use sql.Null* in your struct
	err := db.Select(&listings, selectStmt)

	if err != nil {
		setup.LogCommon(err).Error("Select all listings")
	}

	return listings
}

// NewListings compares two sets of listings by iex_id to see which are not already in the database.
// db are the listings currently in the database. fresh are the listings that may be new.
// Returns the new listings that are enabled on IEX.
func NewListings(db []Listing, fresh []Listing) []Listing {
	new := []Listing{}

	// hashmap to allow O(1) lookup instead of looping
	m := make(map[string]int)
	// put db listings in map
	for i, s := range db {
		m[s.IexID] = i
	}

	// check whether the fresh listings are in the database
	for _, s := range fresh {
		// only care about enabled listings
		if s.IsEnabled == true {
			// if not in the map, add to return array
			_, keyExists := m[s.IexID]

			if !keyExists {
				new = append(new, s)

				setup.LogCommon(nil).
					WithField("listingID", s.ListingID).
					WithField("iexID", s.IexID).
					Info("New Listing")

			}
		}
	}

	return new
}
