# LogXM
Stands for
LOGrus eXtension Management

This is a golang logging library.Created based upon the famous Logrus library.

# features
Default JSON formatted output for Golang servers.
Configurable Log FileName and Directory.
Log Rotation for web servers.
## Able to Set Host Name In JSON
Since most current web services are running in multiple servers,
writing the host name into each log makes trouble shooting easier.

## Log Rotation
Log rotation is an required function for long running web servers.
Log rotation prevent log files occupy server's storage.

# Configuration
If you create a LogXM instance with nil argument. The instance has default settings.
If you would like to configure log, you may create an instance of LoggerConfiguration as an argument of "New" method.

```
type LoggerConfiguration struct {
	Loglevel        logrus.Level // loglevel
	DirName         string       // directory to put log files in.
	WriteToFile     bool         // if true writes to file, false writes to stdout.
	FileName        string       // logfile name
	DateFormat      string       // ex."2006-01-02T15:04:05.999Z07:00"
	LogRotation     int          // max date to hold daily log files. if 0 is set, logfile does not rotate.
	IncludeHostName bool         // always write host name to log.(use for multiple server logging)
}

```


## How to use with default settings
```	
logger := New(nil)

logger.Info(`Write info level log`)

logger.Warn(`Write warning level log.`)

```

all methods are compatible to logrus library.