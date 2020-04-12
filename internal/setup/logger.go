package setup

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wpwilson10/utility/setup"
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

// AppLog tracks the number of each warning level
type AppLog struct {
	AppName    string
	RunTime    string // use string to we can round that values before passing the html table
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

// SummarizeLog sends an email containing the number of log levels per application.
func SummarizeLog() {
	// get the log file
	r := readFile()
	defer r.Close()

	// read each line and count values
	values := parseLogs(r)

	// prepare html files
	t := setupTemplate()

	// buffer stores the output from the writer
	buf := new(bytes.Buffer)
	// this puts the struct values into the html template
	err := t.Execute(buf, values)
	if err != nil {
		setup.LogCommon(err).Error("Failed Execute")
	}

	// send email with the info
	Email("Log Summary", buf.String())
}

// open the three most recent log file for summarization
func readFile() *gzip.Reader {
	// file path
	fp := os.Getenv("LOG_FILEPATH") + os.Getenv("LOG_SUMMARY_FILE")
	// get our log file
	file, err := os.Open(fp)
	if err != nil {
		LogCommon(err).Fatal("Failed Open File")
	}
	// unzip file
	r, err := gzip.NewReader(file)
	if err != nil {
		LogCommon(err).Fatal("Failed gzip reader")
	}

	return r
}

// get and parse html template file
func setupTemplate() *template.Template {
	// get .env filepath
	absPath, err := filepath.Abs("./configs/")
	if err != nil {
		LogCommon(err).Error("Template filepath")
	}
	templatePath := absPath + "/" + os.Getenv("LOG_TEMPLATE_FILE")

	// setup html template
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		LogCommon(err).Error("Template parsing")
	}

	return t
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
				n := AppLog{AppName: app}
				n.Increment(level)
				appMap[app] = &n
			}
		}
	}
	if err := scanner.Err(); err != nil {
		LogCommon(err).Error("Scanner error")
	}

	// Convert map to slice of values since that it was the html template likes
	values := []AppLog{}
	for _, value := range appMap {
		values = append(values, *value)
	}

	return values
}
