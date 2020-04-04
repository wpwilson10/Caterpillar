package stocks

import (
	"context"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	iex "github.com/wpwilson10/iexcloud/v2"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// IEXSetup prepares the iexcloud library.
func IEXSetup() *iex.Client {
	var token string = os.Getenv("IEX_PUBLIC_TOKEN")
	var baseURL string = os.Getenv("IEX_BASE_URL")
	var version string = os.Getenv("IEX_VERSION")

	var fullURL string = baseURL + "/" + version

	return iex.NewClient(token, iex.WithBaseURL(fullURL))
}

// IEXIntraday returns yesterday's intraday data for a given listing sorted from oldest to newest.
func IEXIntraday(client *iex.Client, listing Listing) []Intraday {
	// date of data to retreive
	var dateString string = os.Getenv("IEX_INTRADAY_DATE")
	// convert to int
	numDays, err := strconv.Atoi(dateString)
	if err != nil {
		setup.LogCommon(err).
			WithField("dateString", dateString).
			Error("Parse IEX_INTRADAY_DATE")
	}

	// yesterday is negative days ago
	yesterday := time.Now().AddDate(0, 0, -1*numDays)

	// From IEX API - string. Formatted as YYYYMMDD
	options := iex.IntradayHistoricalOptions{
		ExactDate: yesterday.Format("20060102"),
	}

	// IEX API call
	data, err := client.IntradayHistoricalPrices(context.Background(), listing.Symbol, &options)

	if err != nil {
		setup.LogCommon(err).
			WithField("listingID", listing.ListingID).
			WithField("symbol", listing.Symbol).
			Warn("Getting intraday prices")
	}

	// Add data to formatted struct
	intradayData := make([]Intraday, len(data))

	for i, s := range data {

		// Convert date and time strings to time struct
		layout := "2006-01-02 15:04"
		datetime, err := time.Parse(layout, s.Date+" "+s.Minute)

		if err != nil {
			setup.LogCommon(err).
				WithField("listingID", listing.ListingID).
				WithField("symbol", listing.Symbol).
				WithField("datetime", s.Date+" "+s.Minute).
				Error("Intraday datetime conversion")
		}

		// Convert iex library data to standard struct
		intradayData[i] = Intraday{
			IntradayID: 0,
			ListingID:  listing.ListingID,
			DataTime:   datetime,
			SourceID:   IEXSource,
			UpdateTime: time.Now(),
			Open:       decimal.NewFromFloat(s.Open),
			Close:      decimal.NewFromFloat(s.Close),
			High:       decimal.NewFromFloat(s.High),
			Low:        decimal.NewFromFloat(s.Low),
			Volume:     decimal.NewFromInt(int64(s.Volume)),
			Notional:   decimal.NewFromFloat(s.Notional),
			NumTrades:  decimal.NewFromInt(int64(s.NumberOfTrades)),
		}
	}

	// sort data from newest to oldest
	sort.Slice(intradayData, func(i, j int) bool {
		return intradayData[i].DataTime.Before(intradayData[j].DataTime)
	})

	return intradayData
}

// IEXSymbols returns all listings that are currently supported on IEX.
// This does not consider what is in our database or filter by any criteria.
func IEXSymbols(client *iex.Client) []Listing {
	// IEX API Call
	symbols, err := client.Symbols(context.Background())

	if err != nil {
		setup.LogCommon(err).Fatal("Getting IEX symbols")
	}

	listings := make([]Listing, len(symbols))

	for i, s := range symbols {
		// Convert iex library data to standard struct
		listings[i] = Listing{
			ListingID:  0,
			UpdateTime: time.Now(),
			IsEnabled:  s.IsEnabled,
			Symbol:     s.Symbol,
			Name:       s.Name,
			IexID:      s.IEXID,
			Type:       s.Type,
			Region:     s.Region,
			Currency:   s.Currency,
			Exchange:   s.Exchange,
		}
	}

	return listings
}

// ChangedOnIEX compares listings matching on IexID to find which have changed IEX fields.
// Returns the original listings that should go to audit, and the updated listings.
func ChangedOnIEX(db []Listing, fresh []Listing) (original []Listing, updated []Listing) {

	// hashmap to allow O(1) lookup instead of looping
	m := make(map[string]int)
	// put db listings in map
	for i, s := range db {
		m[s.IexID] = i
	}

	// check whether the fresh listings are in the database
	for _, f := range fresh {
		// Find a matching listing
		i, keyExists := m[f.IexID]

		var flag bool // true if listings do not match

		if keyExists {
			flag = false

			// get the listing
			s := db[i]

			// copy original listing so that fields we don't care about in this function don't change
			temp := s

			// Check whether values from IEX have changed. Dont need ListingID, UpdateTime. Iex_ID must be the same
			// If they have changed, put the fresh value in temp
			if s.IsEnabled != f.IsEnabled {
				flag = true
				temp.IsEnabled = f.IsEnabled
			}
			if strings.Compare(s.Symbol, f.Symbol) != 0 {
				flag = true
				temp.Symbol = f.Symbol
			}
			if strings.Compare(s.Name, f.Name) != 0 {
				flag = true
				temp.Name = f.Name
			}
			if strings.Compare(s.Type, f.Type) != 0 {
				flag = true
				temp.Type = f.Type
			}
			if strings.Compare(s.Region, f.Region) != 0 {
				flag = true
				temp.Region = f.Region
			}
			if strings.Compare(s.Currency, f.Currency) != 0 {
				flag = true
				temp.Currency = f.Currency
			}
			if strings.Compare(s.Exchange, f.Exchange) != 0 {
				flag = true
				temp.Exchange = f.Exchange
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
					WithField("iexIDFresh", f.IexID).
					Info("Changed Listing IEX")
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
