package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// default scripts directory
var DIR_SCRIPTS string = "data"

// Runner executes the specified script with the given parameters and returns
// the output
func Runner(name string, param []string) string {
	cmd := exec.Command("bash", name)
	var output string
	for _, item := range param {
		cmd.Args = append(cmd.Args, item)
	}
	out, err := cmd.Output()
	if err != nil {
		output = fmt.Sprintf("Error occurred\n%s", err)
	} else {
		output = fmt.Sprintf("%s", out)
	}
	return output
}

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
