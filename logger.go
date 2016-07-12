package main

import (
	"io"
	"log"
	"os"
)

var logFileName string = "robotester.log"

// init opens or creates (if non existent) a logfile.
// global string 'logFileName' defines the name of the logfile
func init() {
	logfile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", logfile, ":", err)
	}
	multilog := io.MultiWriter(logfile, os.Stdout)
	log.SetOutput(multilog)
}

// LogAppenLine appends a string to the logfile
func LogAppendLine(line string) {
	log.Println(line)
}

// LogFatal writes to logfile and terminates the program when the called
// interface ends
func LogFatal(v ...interface{}) {
	log.Fatal(v)
}

// LogErrors appends an error to the logfile
func LogErrors(err error) {
	log.Println(err.Error())
}
