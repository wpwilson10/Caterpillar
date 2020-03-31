package setup

import (
	"runtime"
)

// FuncName returns the calling function's name
func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
