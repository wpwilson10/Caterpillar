package setup

import (
	"net/url"
	"runtime"
)

// FuncName returns the calling function's name
func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// IsValidURL tests a string to determine if it is a well-structured url or not.
func IsValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
