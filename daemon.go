// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type Job struct {
	script  string
	path    string
	args    []string
	uuid    string
	output  string
	exit    string
	request time.Time
	start   time.Time
	finish  time.Time
}

var jobQueue chan Job

func RunHandler(w http.ResponseWriter, r *http.Request) {
	LogAppendLine(fmt.Sprintf("Responding to /run request from %s", r.Host))
	vars := mux.Vars(r)
	script := vars["script"]

	if hasScript(script) {
		r.ParseForm()
		wd, _ := os.Getwd()

		job := Job{
			script: script,
			args:   r.Form["args"],
			uuid:   uuid.NewV4().String(),
			path:   filepath.Join(wd, DIR_SCRIPTS, script)}

		jobQueue <- job
		fmt.Fprintf(w, "QUEUED with uuid: %s", job.uuid)

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "ERROR  no script named %s found in %s", script, DIR_SCRIPTS)
	}
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I handle /log requests!\n"))
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I handle /status requests!\n"))
}

func init() {
	jobQueue = make(chan Job, 10)
}

func main() {
	LogAppendLine(fmt.Sprintf("START  %s", time.Now()))

	// init http handlers
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/run/{script}", RunHandler).Methods("GET")
	router.HandleFunc("/log/script/{script}", LogHandler).Methods("GET")
	router.HandleFunc("/log/uuid/{uuid}", LogHandler).Methods("GET")
	router.HandleFunc("/status", StatusHandler).Methods("GET")

	// start http server
	http.ListenAndServe(":1337", router)
	os.Exit(1)
}
