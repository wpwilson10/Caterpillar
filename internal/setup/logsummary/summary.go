package logsummary

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/wpwilson10/caterpillar/internal/setup"
)

// SummaryFile stores log counts for each app over several days.
type SummaryFile struct {
	AppLogs []WrapperAppLog
}

// WrapperAppLog is a collection of similar AppLogs by app.
type WrapperAppLog struct {
	AppName      string
	InnerAppLogs []AppLog
}

// Load creates a SummaryFile from the values saved in the LOG_FILEPATH summary file.
func Load() *SummaryFile {
	// file path
	prefix := os.Getenv("LOG_FILEPATH")
	filepath := prefix + "summary.json"
	// get the file
	file, err := os.Open(filepath)
	if err != nil {
		setup.LogCommon(err).Error("Open File")
	}
	defer file.Close()

	// read summary file into our struct
	out := SummaryFile{}
	err = (json.NewDecoder(file)).Decode(&out)
	if err != nil {
		setup.LogCommon(err).Error("Decode File")
	}

	return &out
}

// Save puts the SummaryFile's current values in a file.
func (data *SummaryFile) Save() {
	// create json file from struct
	jsonFile, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		setup.LogCommon(err).Error("Json marshall indent")
	}

	// file path
	prefix := os.Getenv("LOG_FILEPATH")
	filepath := prefix + "summary.json"
	// get the file
	file, err := os.OpenFile(filepath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		setup.LogCommon(err).Error("Open file")
	}
	defer file.Close()
	// write to file
	_, err = file.Write(jsonFile)
	if err != nil {
		setup.LogCommon(err).Error("Write File")
	}
}

// Append adds a new log entry to the summary file while removing the oldest value
// for the same app.
func (data *SummaryFile) Append(newLog *AppLog) {
	// find matching WrapperAppLog
	found := false
	for i, each := range data.AppLogs {
		// check app names match
		if strings.Compare(each.AppName, newLog.AppName) == 0 {
			found = true
			// sort by date using Slice() function
			sort.Slice(each.InnerAppLogs, func(p, q int) bool {
				return each.InnerAppLogs[p].LogDate.Before(each.InnerAppLogs[q].LogDate)
			})
			// we will assume the new log is today
			// add to end because it is newest
			temp := append(each.InnerAppLogs, *newLog)
			// if less than 7 days, just add data
			if len(temp) <= 7 {
				data.AppLogs[i].InnerAppLogs = temp
			} else {
				// cut off most recent day
				data.AppLogs[i].InnerAppLogs = temp[1:]
			}
		}
	}
	if !found {
		// app has not been seen before, so add it
		newWrapper := WrapperAppLog{
			AppName:      newLog.AppName,
			InnerAppLogs: []AppLog{*newLog},
		}
		data.AppLogs = append(data.AppLogs, newWrapper)
	}
}
