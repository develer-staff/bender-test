// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

// NewHash returns a unique sha1 hash.
func NewHash() string {
	h := sha1.New()
	u := fmt.Sprintf("%s", uuid.NewV4())

	h.Write([]byte(u))
	return hex.EncodeToString(h.Sum(nil))
}

// RunScript parses the API request and extracts the script name and arguments
// with the appropriate checks.
// The request must have a 'script' field with the name of the script to call and
// optional 'arg' fields.
// Example: => script.py foo -bar 42
// 			=> run?script=script.py&arg=foo&arg=-bar&arg=42
// If the request is valid, the script dispatcher module is called.
func RunScript(w http.ResponseWriter, r *http.Request) {
	log.Printf("Responding to /run request from %s\n", r.Host)

	r.ParseForm()
	w.WriteHeader(http.StatusOK)
	script := r.FormValue("script")

	if script != "" {
		options := fmt.Sprintf(strings.Join(r.Form["arg"], " "))
		uuid := NewHash()

		w.Write([]byte(fmt.Sprintf("Script: %s\n", script)))
		w.Write([]byte(fmt.Sprintf("Options: %s\n", options)))
		w.Write([]byte(fmt.Sprintf("UUID: %s\n", uuid)))
	} else {
		fmt.Fprintln(w, "ERROR  no script specified")
	}
}

func main() {
	// create and open log file 'robotester.log'
	logfile, err := os.OpenFile("robotester.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", logfile, ":", err)
	}
	multilog := io.MultiWriter(logfile, os.Stdout)
	log.SetOutput(multilog)

	// init http handlers
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/run", RunScript).
		Methods("GET")

	// start http server
	log.Fatal(http.ListenAndServe(":1337", router))
}
