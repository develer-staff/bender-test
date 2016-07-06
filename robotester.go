// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

// GetArgs parses a URL and returns a slice of strings maintaining the order of
// keys and values
func GetArgs(u *url.URL) []string {
	q := strings.Split(u.String(), "?")[1]
	order := strings.Split(q, "&")

	var opts []string
	for _, v := range order {
		if v != "" {
			if strings.ContainsAny(v, "=") {
				t := strings.SplitN(v, "=", 2)
				key := t[0]
				value := t[1]
				opts = append(opts, key)

				if value != "" {
					opts = append(opts, value)
				}
			} else {
				opts = append(opts, v)
			}
		}
	}
	return opts
}

// RunScript parses the API request and extracts the script name and arguments
// with the appropriate checks.
// The request must have a 'script' field with the name of the script to call and
// optional 'arg' fields.
// Example: script.py foo -bar 42
// 			=> /run/script.py?foo=-bar&42
//			=> /run/script.py?foo&-bar=42
// If the request is valid, the script dispatcher module is called.
func RunScript(w http.ResponseWriter, r *http.Request) {
	log.Printf("Responding to /run request from %s\n", r.Host)

	queries := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	script := queries["script"]

	if script != "" {
		if hasScript(script) {

			args := GetArgs(r.URL)
			uuid := NewHash()
			wd, _ := os.Getwd()
			path := filepath.Join(wd, DIR_SCRIPTS, script)
			out := Runner(path, args)

			w.Write([]byte(fmt.Sprintf("Script: %s\n", script)))
			w.Write([]byte(fmt.Sprintf("Full Path: %s\n", path)))
			w.Write([]byte(fmt.Sprintf("UUID: %s\n", uuid)))

			fmt.Fprintln(w, "\n==== BEGIN OUTPUT ====")
			w.Write([]byte(out))
			fmt.Fprintln(w, "==== END OUTPUT ====")
		} else {
			fmt.Fprintf(w, "ERROR  no script named %s found in %s", script, DIR_SCRIPTS)
		}
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
	router.HandleFunc("/run/{script}", RunScript).
		Methods("GET")

	// start http server
	log.Fatal(http.ListenAndServe(":1337", router))
}
