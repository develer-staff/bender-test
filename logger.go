package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"io/ioutil"
	"errors"
	"strings"
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
	for scr := range jobDone {
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

//ReadLog returns the content of a log file
func ReadLog(path string) (string, error) {
	output, err := ioutil.ReadFile(path)
	var log string

	if err != nil {
		err = errors.New("Log not found")
	}

	log = string(output)
	return log, err
}

//ReadLogDir returns the content of each file in the given dir
func ReadLogDir(path string) (string, error) {
	out_log := ""
	files, err := ioutil.ReadDir(path)
	for _, file := range files {
		file_path := filepath.Join(path, file.Name())
		tmp, _ := ReadLog(file_path)
		out_log += tmp
		out_log += "\n\n*******************\n\n"
	}

	if err != nil || len(files) == 0{
		err = errors.New("No logs found for the given script")
	}

	return out_log, err
}

//FindLog returns the path of the log file
//for the given id
func FindLog(id string) (string, error) {
	path := ""
	log_path, err_log := filepath.Abs("log")
	var err error
	var files []os.FileInfo

	dirs, err_log := ioutil.ReadDir(log_path)

	for _, dir_path := range dirs {
		dir, _ := os.Open(filepath.Join(log_path, dir_path.Name()))
		files, err = dir.Readdir(-1)
		for _, file := range files {
			if strings.Contains(file.Name(), id) {
				path = filepath.Join(dir.Name(), file.Name())
			}
		}
		if err != nil || path == ""{
			err = errors.New("No log found")
		}
	}

	if err_log != nil{
		err = errors.New("No log found")
	}
	return path, err
}