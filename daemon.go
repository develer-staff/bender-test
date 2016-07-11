// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocraft/web"
	"github.com/satori/go.uuid"
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

// SetDefaults initializes Context variables
func (c *Context) SetDefaults(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.ScriptsDir = GetScriptsDir()
	next(w, r)
}

// RunHandler handles /run requests
func (c *Context) RunHandler(w web.ResponseWriter, r *web.Request) {
	LogAppendLine(fmt.Sprintf("Requested execution of script '%s'", r.PathParams["script"]))
	r.ParseForm()

	job := Job{
		Script:  r.PathParams["script"],
		Args:    r.Form["args"],
		Uuid:    uuid.NewV4().String(),
		Path:    c.ScriptsDir,
		Request: time.Now(),
		Status:  "queued"}

	// encode job into a json
	js, err := json.Marshal(job)
	if err != nil {
		LogErrors(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

	// start http server
	LogFatal(http.ListenAndServe(":8080", router))
}
