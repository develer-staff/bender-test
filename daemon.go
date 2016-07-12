// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocraft/web"
)

type Context struct {
	ScriptsDir string
}

type Job struct {
	Script  string
	Path    string
	Args    []string
	Uuid    string
	Output  string
	Exit    string
	Request time.Time
	Start   time.Time
	Finish  time.Time
	Status  string
}

func (job Job) ToJson() ([]byte, error) {
	return json.Marshal(job)
}

// SetDefaults initializes Context variables
func (c *Context) SetDefaults(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.ScriptsDir = GetScriptsDir()
	next(w, r)
}

// RunHandler handles /run requests
func (c *Context) RunHandler(w web.ResponseWriter, r *web.Request) {
	LogAppendLine(fmt.Sprintf("Requested execution of script '%s'", r.PathParams["script"]))

	// parse http request query
	r.ParseForm()
	script := r.PathParams["script"]
	args := r.Form["args"]
	path := c.ScriptsDir
	requested := time.Now()

	// Validate request and build job
	w.Header().Set("Content-Type", "application/json")
	job, err := NewJob(script, args, path, requested)
	if err != nil {
		LogErrors(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Queue job (if there's space left)
	select {
	case jobQueue <- job:
	default:
		LogAppendLine(fmt.Sprintf("NOQUEUE  job queue is full :("))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Build and send response json
	job.Status = "queued"
	js, errj := job.ToJson()
	if errj != nil {
		LogErrors(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

// LogHandler handles /log requests
func (c *Context) LogHandler(w web.ResponseWriter, r *web.Request) {
	if r.PathParams["script"] != "" {
		fmt.Fprintf(w, "Requested log for script '%s'\n", r.PathParams["script"])
	} else if r.PathParams["uuid"] != "" {
		fmt.Fprintf(w, "Requested log for uuid '%s'\n", r.PathParams["uuid"])
	}
}

// StatusHandler handles /status requests
func (c *Context) StatusHandler(w web.ResponseWriter, r *web.Request) {
	if r.PathParams["script"] != "" {
		fmt.Fprintf(w, "Requested job status for script '%s\n'", r.PathParams["script"])
	} else if r.PathParams["uuid"] != "" {
		fmt.Fprintf(w, "Requested job status for uuid '%s'\n", r.PathParams["uuid"])
	} else {
		fmt.Fprintln(w, "Requested server status (general)")
		fmt.Fprintf(w, "  scripts dir: '%s'\n", c.ScriptsDir)
	}
}

func main() {
	LogAppendLine(fmt.Sprintf("START  %s", time.Now()))

	// init http handlers
	router := web.New(Context{})
	router.Middleware((*Context).SetDefaults)
	router.Get("/run/:script", (*Context).RunHandler)
	router.Get("/log/script/:script", (*Context).LogHandler)
	router.Get("/log/uuid/:uuid", (*Context).LogHandler)
	router.Get("/status", (*Context).StatusHandler)
	router.Get("/status/script/:script", (*Context).StatusHandler)
	router.Get("/status/uuid/:uuid", (*Context).StatusHandler)

	// worker
	go RunWorker()

	// start http server
	LogFatal(http.ListenAndServe(":8080", router))
}
