package logxm

import (
	"log"
	"os"
	"path/filepath"

	logrus "github.com/sirupsen/logrus"
)

/**
 * logxm is the Log library for Golang server.
 */

var logger *logrus.Entry

var logFile *os.File

// LoggerConfiguration ...
type LoggerConfiguration struct {
	DirName     string
	WriteToFile bool
	FileName    string
	DateFormat  string
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
	if config.WriteToFile {
		err := os.Mkdir(`.`+string(filepath.Separator)+config.DirName, os.FileMode(mode))
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.OpenFile(`./`+config.DirName+`/`+config.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(mode))
	if err != nil {
		log.Fatal(err)
	}
	f := new(logrus.JSONFormatter)
	f.TimestampFormat = config.DateFormat
	lg.Formatter = f

	hostname, _ := os.Hostname()
	logger = lg.WithField("host", hostname) //always write log with hostname.

	if config.WriteToFile {
		lg.SetOutput(file)
	} else {
		lg.SetOutput(os.Stdout)
	}
	logger.Info("Logrus is Setup for logging.")
}

// StandardConfig is Standard Configuration for Logxm
func StandardConfig() *LoggerConfiguration {
	return &LoggerConfiguration{DirName: `log`, FileName: `application.log`, WriteToFile: true, DateFormat: "2006-01-02T15:04:05.999Z07:00"}
}

// TerminateLogging : use this when you want to finish writing log.
func TerminateLogging() {
	logFile.Close()
}
