package logsummary

import (
	"bufio"
	"compress/gzip"
	"strings"
	"time"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// AppLog tracks the number of each warning level
type AppLog struct {
	AppName    string
	LogDate    time.Time // date of the log
	RunTime    string    // use string to we can round the values before passing to the html table
	TraceCount int
	DebugCount int
	InfoCount  int
	WarnCount  int
	ErrorCount int
	FatalCount int
	PanicCount int
}

// Increment the count for a given level
func (log *AppLog) Increment(level string) {
	switch level {
	case "trace":
		log.TraceCount = log.TraceCount + 1
	case "debug":
		log.DebugCount = log.DebugCount + 1
	case "info":
		log.InfoCount = log.InfoCount + 1
	case "warning":
		log.WarnCount = log.WarnCount + 1
	case "error":
		log.ErrorCount = log.ErrorCount + 1
	case "fatal":
		log.FatalCount = log.FatalCount + 1
	case "panic":
		log.PanicCount = log.PanicCount + 1
	}
}

// checks each line and counts number of occurences of log info
func parseLogs(r *gzip.Reader) []AppLog {
	// stores count of times we see app
	appMap := make(map[string]*AppLog)

	// read each line
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// get data we need to increment the app count
		var app string
		var level string
		var runtime string
		// get line
		line := scanner.Text()
		// split line into pieces
		pieces := strings.Split(line, " ")
		// check each piece for what we care about
		for _, each := range pieces {
			// application
			if strings.Contains(each, "application=") {
				app = each[len("application="):]
			}
			// level
			if strings.Contains(each, "level=") {
				level = each[len("level="):]
			}
			// run time
			if strings.Contains(each, "RunTime=") {
				runtime = each[len("RunTime="):]
			}
		}

		// check if app in map
		count, check := appMap[app]
		if check {
			// if so increment count and runtime for that appLog
			if len(app) > 2 && len(level) > 2 {
				count.Increment(level)
			}
			if len(runtime) > 2 {
				rt, err := time.ParseDuration(runtime)
				if err != nil {
					setup.LogCommon(err).Error("Failed parsing runtime")
				} else {
					count.RunTime = rt.Round(time.Second).String()
				}
			}
		} else {
			// if not in map and got data, create new app log
			if len(app) > 2 && len(level) > 2 {
				n := AppLog{
					AppName: app,
					LogDate: time.Now().AddDate(0, 0, -1), // should be yesterday's log
				}
				n.Increment(level)
				appMap[app] = &n
			}
		}
	}
	if err := scanner.Err(); err != nil {
		setup.LogCommon(err).Error("Scanner error")
	}

	// Convert map to slice of values since that it was the html template likes
	values := []AppLog{}
	for _, value := range appMap {
		values = append(values, *value)
	}

	return values
}
