/*
	Package setup contains common tools for wpwilson10 applications.
*/

package setup

import "time"

// ApplicationName is the currently running program
var ApplicationName string

// startTime is used for timing application runs
var startTime time.Time

// Application sets up global variables
func Application(app string) {
	ApplicationName = app
	startTime = time.Now()
}

// RunTime returns the difference between setup.Application call and now.
func RunTime() time.Duration {
	now := time.Now()
	return now.Sub(startTime)
}

// CheckPythonServer blocks the function call until the python server is up.
func CheckPythonServer() {
	// check if python server is running
	var count int = 1
	for !CheckOnce(EnvToInt("PY_CATERPILLAR_PORT")) {
		LogCommon(nil).
			WithField("count", count).
			Warn("Caterpillar python server not running")
		// server should restart every 5 minutes
		time.Sleep(time.Minute * 5)
		count = count + 1
	}
}
