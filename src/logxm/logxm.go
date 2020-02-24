package logxm

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/robfig/cron"

	logrus "github.com/sirupsen/logrus"
)

/**
 * logxm is the Log library for Golang server.
 */

// LoggerConfiguration is a configuration for logging
type LoggerConfiguration struct {
	DirName     string // directory to put log files in.
	WriteToFile bool   // if true writes to file, false writes to stdout.
	FileName    string // logfile name
	DateFormat  string // ex."2006-01-02T15:04:05.999Z07:00"
	LogRotation int    // max date to hold daily log files. if 0 is set, logfile does not rotate.
}

/**
 * Using Logrus Library for Logging formatter.
 * See https://github.com/sirupsen/logrus for detail.
 * This library formats log to JSON format to make it easier to read from other log analyzer.
 * useFile true:uses logfile false: outputs to standard output
 */

// New : Create new Logger if the config was set to nil, default configuration is applyed.
func New(config *LoggerConfiguration) *logrus.Logger {
	return setupLog(config)
}

// SetupLog : Call this function first to start logging
func setupLog(config *LoggerConfiguration) *logrus.Logger {
	// if no configuration is set, use default
	if config == nil {
		config = StandardConfig()
	}
	// Configure Log Formats
	var log = logrus.New()
	// Make Log Directory
	createLogDir(config)
	logFile := getLogFile(config)
	f := new(logrus.JSONFormatter)
	f.TimestampFormat = config.DateFormat
	log.Formatter = f
	hostname, _ := os.Hostname()
	log.WithField("host", hostname) //always write log with hostname.

	if config.WriteToFile {
		writer := &fileWriter{w: logFile}
		log.SetOutput(writer)
		if config.LogRotation > 0 {
			go rotateLogging(config, writer)
		}
	} else {
		log.SetOutput(os.Stdout)
	}
	log.Info("LogXM is Setup for logging.")
	return log
}

func createLogDir(config *LoggerConfiguration) {
	mode := int32(0777)
	if config.WriteToFile { // create log directory
		logDir := `.` + string(filepath.Separator) + config.DirName
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			if err := os.Mkdir(logDir, os.FileMode(mode)); err != nil {
				fmt.Println(err)
			}
		}
	}
}

// StandardConfig is Standard Configuration for Logxm
func StandardConfig() *LoggerConfiguration {
	return &LoggerConfiguration{DirName: `log`, FileName: `application`, WriteToFile: true,
		DateFormat: "2006-01-02T15:04:05.999Z07:00", LogRotation: 3}
}

func getLogFile(config *LoggerConfiguration) *os.File {
	mode := int32(0777)
	file, err := os.OpenFile(getLogFileName(config), os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(mode))
	if err != nil {
		fmt.Println(err)
	}
	return file
}

func getLogFileName(config *LoggerConfiguration) string {
	if config.LogRotation == 0 {
		return `./` + config.DirName + `/` + config.FileName + `.log`
	}
	const dateFormat = `20060102030405`
	return `./` + config.DirName + `/` + config.FileName + `_` + time.Now().Format(dateFormat) + `.log`
}

// rotates log file every day.
func rotateLogging(config *LoggerConfiguration, w *fileWriter) {
	cron := cron.New()
	cron.AddFunc("0 0 * * * ", func() {
		fmt.Println(`Log rotate executed.`)
		deleteOutDatedLogFile(config, w)
		replaceFileByRotation(config, w)
	})
}

func deleteOutDatedLogFile(config *LoggerConfiguration, w *fileWriter) bool {
	var oldFile = w.w.(*os.File)
	w.history = append(w.history, oldFile)
	if len(w.history) >= config.LogRotation {
		// delete last appended file
		if err := os.Remove(w.history[0].Name()); err != nil {
			fmt.Printf(`Failed to delete old file %v`, err)
		}
		//delete file name from history
		w.history = w.history[1:]
	}
	return true
}

func replaceFileByRotation(config *LoggerConfiguration, w *fileWriter) {
	nextFile := getLogFile(config)
	w.exchange(nextFile)
}

type fileWriter struct {
	w       io.Writer
	mutex   sync.Mutex
	history []*os.File
}

func (f *fileWriter) Write(b []byte) (int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	fmt.Println(`Writing...`)
	return f.w.Write(b)
}

// 書き込み先を交換する
func (f *fileWriter) exchange(newFile *os.File) {
	f.mutex.Lock()
	if closer, ok := f.w.(io.Closer); ok {
		defer closer.Close()
	} else {
		fmt.Println(`Failed to close past file.`)
	}
	defer f.mutex.Unlock()
	fmt.Println(`File exchanged.`)
	f.w = newFile
}
