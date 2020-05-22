package logsummary

import (
	"bytes"
	"compress/gzip"
	"html/template"
	"os"
	"path/filepath"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// SummarizeLog sends an email containing the number of log levels per application.
func SummarizeLog() {
	// get the latest log file
	logFile := readLastLog()
	defer logFile.Close()

	// read each line and count values from lastest log
	appLogs := parseLogs(logFile)

	// get the log summary file
	summaryFile := Load()
	// put logs into summary file
	for _, each := range appLogs {
		summaryFile.Append(&each)
	}
	// save the new data
	summaryFile.Save()

	// convert to output format with averages
	values := []Average{}
	for _, each := range appLogs {
		temp := NewAverage(&each, summaryFile)
		values = append(values, *temp)
	}

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
	setup.SendEmail("Log Summary", buf.String())
}

// Open the most recent log file for summarization
// Caller must close file.
func readLastLog() *gzip.Reader {
	// file path
	fp := os.Getenv("LOG_FILEPATH") + os.Getenv("LOG_SUMMARY_FILE")
	// get our log file
	file, err := os.Open(fp)
	if err != nil {
		setup.LogCommon(err).Fatal("Failed Open File")
	}
	// unzip file
	r, err := gzip.NewReader(file)
	if err != nil {
		setup.LogCommon(err).Fatal("Failed gzip reader")
	}

	return r
}

// get and parse html template file
func setupTemplate() *template.Template {
	// get .env filepath
	absPath, err := filepath.Abs("./configs/")
	if err != nil {
		setup.LogCommon(err).Error("Template filepath")
	}
	templatePath := absPath + "/" + os.Getenv("LOG_TEMPLATE_FILE")

	// setup html template
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		setup.LogCommon(err).Error("Template parsing")
	}

	return t
}
