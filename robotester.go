// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

// NewHash returns a sha1 hash of a randomly generated string of the specified
// size
func NewHash(size int) string {
	chars := make([]byte, size)
	for i := 0; i < size; i++ {
		chars[i] = byte(rand.Intn(255))
	}
	h := sha1.New()
	h.Write([]byte(string(chars)))
	return hex.EncodeToString(h.Sum(nil))
}

// RunScript parses the API request and extracts the script name and arguments
// with the appropriate checks.
// The request must have a 'script' field with the name of the script to call and
// optional 'arg' fields.
// Example: script.py foo -bar 42
//			=> /run/script.py?args=foo+-bar+42
// If the request is valid, the script dispatcher module is called.
func RunScript(w http.ResponseWriter, r *http.Request) {
	log.Printf("Responding to /run request from %s\n", r.Host)

	vars := mux.Vars(r)
	r.ParseForm()
	w.WriteHeader(http.StatusOK)
	script := vars["script"]

	if hasScript(script) {
		args := r.Form["args"]
		uuid := NewHash(32)
		wd, _ := os.Getwd()
		path := filepath.Join(wd, DIR_SCRIPTS, script)
		out := Runner(path, args)

		w.Write([]byte(fmt.Sprintf("Script: %s\n", script)))
		w.Write([]byte(fmt.Sprintf("Args: %s\n", args)))
		w.Write([]byte(fmt.Sprintf("Full Path: %s\n", path)))
		w.Write([]byte(fmt.Sprintf("UUID: %s\n", uuid)))

		fmt.Fprintln(w, "\n==== BEGIN SCRIPT OUTPUT ====")
		w.Write([]byte(out))
		fmt.Fprintln(w, "==== END SCRIPT OUTPUT ====")

	} else {
		fmt.Fprintf(w, "ERROR  no script named %s found in %s", script, DIR_SCRIPTS)
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
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
	router.HandleFunc("/run/{script}", RunScript).Methods("GET")

	// start http server
	log.Fatal(http.ListenAndServe(":1337", router))
}
