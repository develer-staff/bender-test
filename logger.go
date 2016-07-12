package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

//WriteLog take a Job struct and save it in log/
func WriteLog() {
	for {
		scr := <-jobDone
		log_path, _ := filepath.Abs(filepath.Join("log", scr.Script))
		if _, err := os.Stat(log_path); os.IsNotExist(err) {
			os.MkdirAll(log_path, 0777)
		}

		now := time.Now()
		file_name := fmt.Sprintf("%d.%d.%d-%d.%d.%d-%s.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), scr.Uuid)
		file_path := filepath.Join(log_path, file_name)
		outfile, _ := os.OpenFile(file_path, os.O_CREATE|os.O_WRONLY, 0666)

		joutput, err := scr.ToJson()

		if err != nil {
			fmt.Fprintln(outfile, "Failed to convert to json")
		} else {
			fmt.Fprintf(outfile, "%s", joutput)
		}
		LogAppendLine(fmt.Sprintf("LOGGER log succesfully wrote"))
	}
}

//ReadLogDir returns the content of each file in the given dir
func ReadLogDir(path string) string {
	out_log := ""
	dir, err := os.Open(path)

	if err != nil {
		out_log = "no logs found for the given script"
	} else {
		files, _ := dir.Readdir(-1)

		for _, file := range files {
			file_path := filepath.Join(path, file.Name())
			out_log += ReadLog(file_path)
			out_log += "\n\n*******************\n\n"
		}
	}

	return out_log
}

//ReadLog returns the content of a log file
func ReadLog(path string) string {
	output, err := ioutil.ReadFile(path)
	var log string

	if err != nil {
		log = "log not found"
	}

	log = string(output)
	return log
}

//FindLog returns the path of the log file
//for the given id
func FindLog(id string) string {
	path := "log not found"
	log_path, _ := filepath.Abs("log")

	log_dir, _ := os.Open(log_path)
	dirs, _ := log_dir.Readdir(-1)

	for _, dir_path := range dirs {
		dir, _ := os.Open(filepath.Join(log_path, dir_path.Name()))
		files, _ := dir.Readdir(-1)
		for _, file := range files {
			if strings.Contains(file.Name(), id) {
				path = filepath.Join(log_path, dir_path.Name(), file.Name())
			}
		}
	}
	return path
}
