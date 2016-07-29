// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gocraft/web"
	"github.com/satori/go.uuid"
)

var logContextDaemon LoggerContext
var daemon_localStatus *StatusModule

type statusJobs struct {
	Jobs []Job `json:"jobs"`
}

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
	LogInf(logContextDaemon, "Receive RUN[%v] request from: %v", "Daemon", r.RemoteAddr)
	r.ParseForm()

	name := r.PathParams["script"]
	uuid := uuid.NewV4().String()
	timeout := 2000
	params := r.Form

	status, _ := daemon_localStatus.GetState()
	if status == SERVER_WORKING {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	js, err := json.Marshal(uuid)

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		LogErr(logContextDaemon, "json creation failed")
		return
	}

	w.Write(js)

	Submit(name, uuid, params, timeout)
}

// LogHandler handles /log requests
func (c *Context) LogHandler(w web.ResponseWriter, r *web.Request) {
	//	LogInf(logContextDaemon, "Receive LOG[%v] request from: %v", "Daemon", r.RemoteAddr)
	r.ParseForm()
	name := r.PathParams["script"]
	id := r.Form["uuid"][0]
	var buffer []byte

	if id != "" {
		buffer = Read(name, id, 0, 0)
	}
	js, err := json.Marshal(string(buffer))

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		LogErr(logContextDaemon, "json creation failed")
		return
	}

	w.Write(js)
}

// StatusHandler handles /state requests
func (c *Context) StatusHandler(w web.ResponseWriter, r *web.Request) {
	//general state requests

	if r.RequestURI == "/state" {
		//		LogInf(logContextDaemon, "Receive STATE[%v] request from: %v", "Daemon", r.RemoteAddr)
		js, err := json.Marshal(daemon_localStatus)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			LogErr(logContextDaemon, "json creation failed")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		// script-name specific requests
		r.ParseForm()

		LogInf(logContextDaemon, "Receive STATE[%v] request from: %v", r.PathParams["script"], r.RemoteAddr)

		response := statusJobs{
			Jobs: daemon_localStatus.GetJobs(r.PathParams["script"])}
		js, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			LogErr(logContextDaemon, "json creation failed")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *Context) HomeHandler(w web.ResponseWriter, r *web.Request) {
	LogInf(logContextDaemon, "Receive HOME[%v] request from: %v", "Daemon", r.RemoteAddr)
	job := Job{Name: "hellosleep",
		Uuid:    "acfjsif-909D",
		Created: time.Now(),
		Timeout: 54}

	t := template.New("New template")
	t, _ = template.ParseFiles("html/home.html")
	t.Execute(w, job)
}

func (c *Context) ListHandler(w web.ResponseWriter, r *web.Request) {
	LogInf(logContextDaemon, "Receive LIST[%v] request from: %v", "Daemon", r.RemoteAddr)
	scripts := List()
	js, err := json.Marshal(scripts)

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		LogErr(logContextDaemon, "json creation failed")
		return
	}

	w.Write(js)
}

func DaemonInit(sm *StatusModule, cm *ConfigModule) {

	daemon_localStatus = sm

	// init logger
	logContextDaemon = LoggerContext{
		level: cm.GetLogLevel("daemon", 3),
		name:  "DAEMON"}
	LogInf(logContextDaemon, "START")

	// init http handlers
	router := web.New(Context{})
	router.Middleware((*Context).SetDefaults)
	router.Get("/run/:script", (*Context).RunHandler)
	router.Get("/log/:script", (*Context).LogHandler)
	router.Get("/state", (*Context).StatusHandler)
	router.Get("/state/:script", (*Context).StatusHandler)
	router.Get("/", (*Context).HomeHandler)
	router.Get("/service/list", (*Context).ListHandler)

	// start http server
	address := cm.Get("daemon", "address", "0.0.0.0")
	port := cm.Get("daemon", "port", "8080")
	LogInf(logContextDaemon, "Listening on %s:%s", address, port)
	LogFatal(logContextDaemon, http.ListenAndServe(address+":"+port, router))
}
