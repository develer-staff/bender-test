// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func RunHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I handle /run requests!\n"))
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I handle /log requests!\n"))
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I handle /status requests!\n"))
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
