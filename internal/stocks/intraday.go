package stocks

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Intraday stores a single timepoint of data in a format matching the intraday table schema
type Intraday struct {
	IntradayID int64           `db:"intraday_id"`
	ListingID  int64           `db:"listing_id"`
	SourceID   int64           `db:"source_id"`
	DataTime   time.Time       `db:"data_time"`
	UpdateTime time.Time       `db:"update_time"`
	Open       decimal.Decimal `db:"open"`
	Close      decimal.Decimal `db:"close"`
	High       decimal.Decimal `db:"high"`
	Low        decimal.Decimal `db:"low"`
	Volume     decimal.Decimal `db:"volume"`
	Notional   decimal.Decimal `db:"notional"`
	NumTrades  decimal.Decimal `db:"num_trades"`
}

// InsertIntraday inserts all values for all given Intraday structs into the intraday database table.
// Currently performs no validation.
func InsertIntraday(db *sqlx.DB, data []Intraday) {
	// make sure we have new data, else panics
	if len(data) > 0 {
		// Setup
		var insertStmt string = "INSERT INTO Intraday (intraday_id, listing_id, source_id, data_time, update_time, open, close, high, low, volume, notional, num_trades)"
		var valueStmt string = "VALUES (DEFAULT, :listing_id, :source_id, :data_time, Now(), :open, :close, :high, :low, :volume, :notional, :num_trades)"
		var fullStmt string = insertStmt + " " + valueStmt

		// start a transaction
		tx, err := db.Beginx()

		if err != nil {
			setup.LogCommon(err).Warn("Intraday setup transaction")
		}

		for _, s := range data {
			// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
			_, err := tx.NamedExec(fullStmt, &s)

			if err != nil {
				setup.LogCommon(err).Warn("Intraday inserting rows")
			}
		}

		err = tx.Commit()

		if err != nil {
			setup.LogCommon(err).Warn("Intraday commit transaction")
		}

		setup.LogCommon(nil).
			WithField("listingID", data[0].ListingID).
			WithField("firstDataTime", data[0].DataTime).
			WithField("lastDataTime", data[len(data)-1].DataTime).
			Info("Inserting intraday data")
	}
}

// SanitizeIntraday removes data is that is old or all zeros values.
// Returns the data that should be input into the database.
func SanitizeIntraday(listing Listing, data []Intraday, latestTimes map[int64]time.Time) []Intraday {
	// check if we got any data
	if !(len(data) > 0) {
		setup.LogCommon(nil).
			WithField("listingID", listing.ListingID).
			WithField("symbol", listing.Symbol).
			Debug("Sanitize Intraday empty data")
	}

	lastTime := latestTimes[listing.ListingID]

	// check if we got a real time
	if lastTime.IsZero() {
		setup.LogCommon(nil).
			WithField("listingID", listing.ListingID).
			WithField("symbol", listing.Symbol).
			Warn("Sanitize Intraday no last date")
	}

	out := []Intraday{}

	for _, d := range data {
		// Add data that is newer than last data time and have non-zero data
		// run when there is no data
		if (d.DataTime.After(lastTime) || lastTime.IsZero()) && d.NumTrades.IsPositive() && d.Volume.IsPositive() && d.Open.IsPositive() && d.Close.IsPositive() {
			out = append(out, d)
		}
	}

	return out
}

// LatestIntraday returns the most recent data time for each listing in the database.
// Returns a map[ListingID] = DataTime
func LatestIntraday(db *sqlx.DB) map[int64]time.Time {

	// hashmap to allow O(1) lookup instead of looping
	m := make(map[int64]time.Time)

	// fetch latest intraday times from the db
	var selectStmt string = "SELECT listing_id, Max (data_time) FROM Intraday GROUP BY listing_id"
	rows, err := db.Query(selectStmt)

	if err != nil {
		setup.LogCommon(err).Error("Select latest intraday")
	}

	// iterate over each row
	for rows.Next() {
		var listingID int64
		var update time.Time

		err = rows.Scan(&listingID, &update)
		if err != nil {
			setup.LogCommon(err).Error("Scan latest intraday")
		}

		// insert into map
		m[listingID] = update
	}

	return m
}
