package main

import (
	"io/ioutil"
	"strings"
)

// default scripts directory
var DIR_SCRIPTS string = "scripts"

// hasScript checks for the script existance
func hasScript(search string) bool {
	files, err := ioutil.ReadDir(DIR_SCRIPTS)
	if err == nil {
		k := len(files)
		for i := 0; i < k; i++ {
			if strings.Contains(files[i].Name(), search) {
				return true
			}
		}
	}
	return false
}

// listScripts returns a list of scripts in the default script folder
func listScripts() []string {
	files, err := ioutil.ReadDir(DIR_SCRIPTS)
	if err == nil {
		names := []string{}
		k := len(files)
		for i := 0; i < k; i++ {
			names = append(names, files[i].Name())
		}
		return names
	}
	return nil
}
