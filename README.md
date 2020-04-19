# LogXM
Stands for
LOGrus Extension Management

This is a golang logging library.Created based upon the famous Logrus library.

## features
Default JSON formatted output for Golang servers.
Configurable Log FileName and Directory.
Log Rotation for web servers.

## How to use with default settings
```	
logger := New(nil)

logger.Info(`Write info level log`)

logger.Warn(`Write warning level log.`)

```

all methods are compatible to logrus library.