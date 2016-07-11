// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gocraft/web"
)

type Context struct {
	ScriptsDir string
}

// SetDefaults initializes Context variables
func (c *Context) SetDefaults(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.ScriptsDir = GetScriptsDir()
	next(w, r)
}

// RunHandler handles /run requests
func (c *Context) RunHandler(w web.ResponseWriter, r *web.Request) {
	fmt.Fprintf(w, "Requested execution of script '%s'\n", r.PathParams["script"])
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
