package logxm

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

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

// New : Create new Logger
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
		log.SetOutput(writer.w)
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
	return &LoggerConfiguration{DirName: `log`, FileName: `application.log`, WriteToFile: true,
		DateFormat: "2006-01-02T15:04:05.999Z07:00", LogRotation: 7}
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
		return `./` + config.DirName + `/` + config.FileName
	}
	const dateFormat = `20060102030405`
	return `./` + config.DirName + `/` + config.FileName + `_` + time.Now().Format(dateFormat)
}

// rotates log file every day.
func rotateLogging(config *LoggerConfiguration, w *fileWriter) {
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			// rename current logging file and create new one
			fmt.Println(`Log rotate executed.`)
			replaceFileByRotation(config, w)
		}
	}
}

func replaceFileByRotation(config *LoggerConfiguration, w *fileWriter) {
	nextFile := getLogFile(config)
	w.exchange(nextFile)
}

type fileWriter struct {
	w     io.Writer
	mutex sync.Mutex
}

func (f *fileWriter) Write(b []byte) (int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.w.Write(b)
}

// 書き込み先を交換する
func (f *fileWriter) exchange(newFile *os.File) {
	f.mutex.Lock()
	if closer, ok := f.w.(io.Closer); ok {
		defer closer.Close()
	}
	defer f.mutex.Unlock()
	f.w = newFile
}
