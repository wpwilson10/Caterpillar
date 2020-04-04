package setup

import (
	"fmt"
	"net"
	"strings"
)

// RunOnce checks if the given app is already running, and if not runs it.
// Blocks indefinitely, calling code should consider running as go routine.
func RunOnce(port int, appFunc func()) {
	// check if already running
	if !CheckOnce(port) {
		// bind to app's port
		once(port)
		// main app logic, runs and blocks indefinitely
		appFunc()
	}
}

// once binds to the given port on localhost and does nothing forever.
// Useful for when we want a single instances of an app at a time.
func once(port int) {
	// point to port on localhost
	address := fmt.Sprintf("127.0.0.1:%d", port)
	// connect to port
	listener, err := net.Listen("tcp", address)
	// failed, return error
	if err != nil {
		LogCommon(err).
			WithField("port", port).
			Fatal("Failed once listen")
	}
	// send listener off to wait forever
	go silent(listener)
}

func silent(listener net.Listener) {
	// hang out and do nothing forever
	for {
		listener.Accept()
	}
}

// CheckOnce returns true if the given port is busy which represents a running app,
// false otherwise.
func CheckOnce(port int) bool {
	// point to port on localhost
	address := fmt.Sprintf("127.0.0.1:%d", port)
	// try to connect to port
	conn, err := net.Listen("tcp", address)
	// failed, check error
	if err != nil {
		// if we got the expected single connection error, given port is busy
		// windows
		if strings.Index(err.Error(), "Only one usage of each socket address") != -1 {
			return true
		} else if strings.Index(err.Error(), "address already in use") != -1 {
			// ubuntu
			return true
		}

		// got an unexpected error
		LogCommon(err).
			WithField("port", port).
			Fatal("Unexpected error")

		return true
	}
	// close connection else it stays connected
	err = conn.Close()
	if err != nil {
		LogCommon(err).
			WithField("port", port).
			Error("Closing error")
	}
	// nothing using the port
	return false

}
