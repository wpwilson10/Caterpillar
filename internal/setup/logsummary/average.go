package logsummary

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// Average holds the current and averages values for an App Log for the HTML template.
type Average struct {
	AppName           string
	RunTime           string // use string to we can round the values before passing to the html table
	TraceCount        int
	DebugCount        int
	InfoCount         int
	WarnCount         int
	ErrorCount        int
	FatalCount        int
	PanicCount        int
	RunTimeAverage    string // use string to we can round the values before passing to the html table
	TraceCountAverage int
	DebugCountAverage int
	InfoCountAverage  int
	WarnCountAverage  int
	ErrorCountAverage int
	FatalCountAverage int
	PanicCountAverage int
}

// NewAverage uses the given most recent appLog and a SummaryFile to make
// a struct containing values to put in the log summary email.
func NewAverage(currentLog *AppLog, summaryFile *SummaryFile) *Average {
	// get the related app
	appLogs := WrapperAppLog{}
	for _, each := range summaryFile.AppLogs {
		// check app names match
		if strings.Compare(each.AppName, currentLog.AppName) == 0 {
			appLogs = each
		}
	}

	// put in current values
	average := Average{
		AppName:    currentLog.AppName,
		RunTime:    currentLog.RunTime,
		TraceCount: currentLog.TraceCount,
		DebugCount: currentLog.DebugCount,
		InfoCount:  currentLog.InfoCount,
		WarnCount:  currentLog.WarnCount,
		ErrorCount: currentLog.ErrorCount,
		FatalCount: currentLog.FatalCount,
		PanicCount: currentLog.PanicCount,
	}

	// get average counts over last week
	var count int
	var tempDuration time.Duration
	for _, each := range appLogs.InnerAppLogs {
		// handle time
		if each.RunTime != "" {
			thisDuration, err := time.ParseDuration(each.RunTime)
			if err != nil {
				setup.LogCommon(err).
					WithField("runTime", each.RunTime).
					WithField("appName", each.AppName).
					Warn("Parse Duration")
			}
			tempDuration = tempDuration + thisDuration
		} else {
			tempDuration = 0
		}

		// add each count to the running total
		average.TraceCountAverage = average.TraceCountAverage + each.TraceCount
		average.DebugCountAverage = average.DebugCountAverage + each.DebugCount
		average.InfoCountAverage = average.InfoCountAverage + each.InfoCount
		average.WarnCountAverage = average.WarnCountAverage + each.WarnCount
		average.ErrorCountAverage = average.ErrorCountAverage + each.ErrorCount
		average.FatalCountAverage = average.FatalCountAverage + each.FatalCount
		average.PanicCountAverage = average.PanicCountAverage + each.PanicCount
		count = count + 1
	}

	// now get averages
	fmt.Println(average.InfoCountAverage, average.WarnCountAverage)
	average.TraceCountAverage = int(float64(average.TraceCountAverage) / float64(count))
	average.DebugCountAverage = int(float64(average.DebugCountAverage) / float64(count))
	average.InfoCountAverage = int(float64(average.InfoCountAverage) / float64(count))
	average.WarnCountAverage = int(float64(average.WarnCountAverage) / float64(count))
	average.ErrorCountAverage = int(float64(average.ErrorCountAverage) / float64(count))
	average.FatalCountAverage = int(float64(average.FatalCountAverage) / float64(count))
	average.PanicCountAverage = int(float64(average.PanicCountAverage) / float64(count))
	fmt.Println(average.InfoCountAverage, average.WarnCountAverage)

	// handle time
	tempNS := tempDuration.Nanoseconds() / int64(count)
	tempDuration, err := time.ParseDuration(strconv.FormatInt(tempNS, 10) + "ns")
	if err != nil {
		setup.LogCommon(err).
			WithField("tempNS", tempNS).
			WithField("appName", currentLog.AppName).
			Warn("Parse Duration NS")
	}
	average.RunTimeAverage = tempDuration.String()

	return &average
}
