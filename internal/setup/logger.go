package setup

import (
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// Logger configures the logrus package used by whole program.
func Logger(file *os.File) {
	// Default format
	log.SetFormatter(&log.TextFormatter{})
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	// use file if given, otherwise print to std out
	if file != nil {
		log.SetOutput(file)
	} else {
		// Output logs to stdout
		log.SetOutput(os.Stdout)
	}
}

// LogCommon returns a logger containing the optional error, application, and function name of the caller.
// By default, Logrus is protected by a mutex for concurrent writes.
func LogCommon(err error) *log.Entry {
	// this looks like FuncName(), but it needs to be internal here to return the correct function
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc).Name()

	// got an error, use it
	if err != nil {
		return log.WithField("application", ApplicationName).WithField("function", f).WithError(err)
	}
	// no error given
	return log.WithField("application", ApplicationName).WithField("function", f)
}

// LogFile returns a file to use for logging with name like 2006-01-02.txt.
// Caller must clost file.
func LogFile() *os.File {
	// file path
	prefix := os.Getenv("LOG_FILEPATH")
	filepath := prefix + "caterpillar.log"

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
