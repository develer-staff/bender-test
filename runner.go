package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

// Runner executes the specified script with the given parameters and returns
// the output
func (a *appContext) Runner(name string, param []string) string {
	cmd := exec.Command(name, param...)
	var output string
	out, err := cmd.Output()
	if err != nil {
		output = fmt.Sprintf("Error occurred\n%s", err)
	} else {
		output = fmt.Sprintf("%s", out)
	}
	return output
}

// hasScript checks for the script existance
func (a *appContext) HasScript(search string) bool {
	files, err := ioutil.ReadDir(a.ScriptsDir)
	if err != nil {
		return false
	}
	k := len(files)
	for i := 0; i < k; i++ {
		namefile := files[i].Name()[0 : len(files[i].Name())-len(filepath.Ext(files[i].Name()))]
		if namefile == search {
			return true
		}
	}
	return false
}

// listScripts returns a list of scripts in the default script folder
func (a *appContext) ListScripts() []string {
	files, err := ioutil.ReadDir(a.ScriptsDir)
	if err != nil {
		return nil
	}
	names := []string{}
	k := len(files)
	for i := 0; i < k; i++ {
		names = append(names, files[i].Name())
	}
	fmt.Println(names)
	return names
}
