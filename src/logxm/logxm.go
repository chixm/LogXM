package logxm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	logrus "github.com/sirupsen/logrus"
)

/**
 * logxm is the Log library for Golang server.
 */

var logger *logrus.Entry

var logFile *os.File

// LoggerConfiguration is a configuration for logging
type LoggerConfiguration struct {
	DirName     string
	WriteToFile bool
	FileName    string
	DateFormat  string // ex."2006-01-02T15:04:05.999Z07:00"
	LogRotation int    // max date to hold daily log files. if 0 is set, logfile does not rotate.
}

/**
 * Using Logrus Library for Logging formatter.
 * See https://github.com/sirupsen/logrus for detail.
 * This library formats log to JSON format to make it easier to read from other log analyzer.
 * useFile true:uses logfile false: outputs to standard output
 */

// SetupLog : Call this function first to start logging
func SetupLog(config *LoggerConfiguration) {
	// if no configuration is set, use default
	if config == nil {
		config = StandardConfig()
	}
	// Configure Log Formats
	var lg = logrus.New()
	mode := int32(0777)
	if config.WriteToFile { // create log directory
		logDir := `.` + string(filepath.Separator) + config.DirName
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			err := os.Mkdir(logDir, os.FileMode(mode))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	logFile := getLogFile(config)
	f := new(logrus.JSONFormatter)
	f.TimestampFormat = config.DateFormat
	lg.Formatter = f
	hostname, _ := os.Hostname()
	logger = lg.WithField("host", hostname) //always write log with hostname.

	if config.WriteToFile {
		lg.SetOutput(logFile)
	} else {
		lg.SetOutput(os.Stdout)
	}
	if config.LogRotation > 0 {
		go rotateLogging(config)
	}
	logger.Info("LogXM is Setup for logging.")
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
		log.Fatal(err)
	}
	return file
}

func getLogFileName(config *LoggerConfiguration) string {
	return `./` + config.DirName + `/` + config.FileName
}

// TerminateLogging : use this when you want to finish writing log.
func TerminateLogging() {
	logFile.Close()
}

// rotates log file every day.
func rotateLogging(config *LoggerConfiguration) {
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			// rename current logging file and create new one
			fmt.Println(`Log Rotationt executed.`)
			current := getLogFileName(config)
			const dateFormat = `20200101`
			rotate := current + time.Now().Format(dateFormat)
			if err := os.Rename(current, rotate); err != nil {
				fmt.Println(err)
			}
		}
	}
}
