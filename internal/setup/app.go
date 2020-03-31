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
