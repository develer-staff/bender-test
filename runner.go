package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

var scriptsDir string

func GetScriptsDir() string {
	return scriptsDir
}

func SetScriptsDir(dir string) {
	scriptsDir = dir
}

func init() {
	scriptsDir = "scripts"
}

// default scripts directory
var DIR_SCRIPTS string = "scripts"

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
func HasScript(search string) bool {
	files, err := ioutil.ReadDir(DIR_SCRIPTS)
	if err != nil {
		return false
	} else {
		k := len(files)
		for i := 0; i < k; i++ {
			if strings.Contains(files[i].Name(), search) {
				fmt.Println("true")
				return true
			}
		}
		return false
	}
}

// listScripts returns a list of scripts in the default script folder
func ListScripts() []string {
	files, err := ioutil.ReadDir(DIR_SCRIPTS)
	if err != nil {
		return nil
	} else {
		names := []string{}
		k := len(files)
		for i := 0; i < k; i++ {
			names = append(names, files[i].Name())
		}
		fmt.Println(names)
		return names
	}
}
