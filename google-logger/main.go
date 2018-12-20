package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/google/logger"
)

const logPath = "./example.log"

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func main() {
	flag.Parse()

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()

	// Log to system log and a log file, Info logs don't write to stdout.
	loggerOne := logger.Init("LoggerExample", false, true, lf)
	defer loggerOne.Close()
	// Don't to system log or a log file, Info logs write to stdout..
	loggerTwo := logger.Init("LoggerExample", true, false, ioutil.Discard)
	defer loggerTwo.Close()

	loggerOne.Info("This will log to the log file and the system log")
	loggerTwo.Info("This will only log to stdout")
}