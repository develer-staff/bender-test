package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"io/ioutil"
)

var logFileName string = "bender-test.log"

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

//WriteLog take a Job struct and save it in log/
func WriteLog() {
	for {
		scr := <-jobDone
		log_path, _ := filepath.Abs(filepath.Join("log", scr.Script))
		if _, err := os.Stat(log_path); os.IsNotExist(err) {
			os.MkdirAll(log_path, 0774)
		}

		now := time.Now()
		file_name := fmt.Sprintf("%d.%d.%d-%d.%d.%d-%s.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), scr.Uuid)
		file_path := filepath.Join(log_path, file_name)

		joutput, err := scr.ToJson()

		ioutil.WriteFile(file_path, joutput, 0664)
		if err != nil {
			LogErrors(err)
		} else {
			LogAppendLine(fmt.Sprintf("LOGGER log wrote succesfully"))
		}
	}
}
