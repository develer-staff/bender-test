// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type JobStatus string

const (
	JOB_QUEUED     = "queued"
	JOB_NOT_FOUND  = "not found"
	JOB_QUEUE_FULL = "queue full"
	JOB_WORKING    = "working"
	JOB_FAILED     = "failed"
	JOB_COMPLETED  = "completed"
)

type Job struct {
	Script  string    `json:"script"`
	Path    string    `json:"path"`
	Args    []string  `json:"args"`
	Uuid    string    `json:"uuid"`
	Output  string    `json:"output"`
	Exit    string    `json:"exit"`
	Request time.Time `json:"request"`
	Start   time.Time `json:"start"`
	Finish  time.Time `json:"finish"`
	Status  JobStatus `json:"status"`
}

type appContext struct {
	ScriptsDir string
	JobQueue   chan Job
	JobDone    chan Job
}

// initAppContext initializes default scripts directory and channels for job
// handling
func initAppContext() *appContext {
	context := &appContext{
		ScriptsDir: GetScriptsDir(),
		JobQueue:   make(chan Job, 10),
		JobDone:    make(chan Job)}
	return context
}

// RunHandler handles /run requests
func (a *appContext) RunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Requested execution of script '%s'\n", vars["script"])
}

// LogHandler handles /log requests
func (a *appContext) LogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["script"] != "" {
		fmt.Fprintf(w, "Requested log for script '%s'\n", vars["script"])
	} else if vars["uuid"] != "" {
		fmt.Fprintf(w, "Requested log for uuid '%s'\n", vars["uuid"])
	}
}

// StatusHandler handles /status requests
func (a *appContext) StatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["script"] != "" {
		fmt.Fprintf(w, "Requested job status for script '%s\n'", vars["script"])
	} else if vars["uuid"] != "" {
		fmt.Fprintf(w, "Requested job status for uuid '%s'\n", vars["uuid"])
	} else {
		fmt.Fprintln(w, "Requested server status (general)")
	}
}

func main() {
	LogAppendLine("SERVER  listening...")

	// init context
	context := initAppContext()

	// init http handlers
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/run/{script}", context.RunHandler).Methods("GET")
	router.HandleFunc("/log/script/{script}", context.LogHandler).Methods("GET")
	router.HandleFunc("/log/uuid/{uuid}", context.LogHandler).Methods("GET")
	router.HandleFunc("/status", context.StatusHandler).Methods("GET")

	// start http server
	LogFatal(http.ListenAndServe(":8080", router))
}
