package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/satori/go.uuid"
)

var scriptsDir string
var jobQueue chan Job

func GetScriptsDir() string {
	return scriptsDir
}

func SetScriptsDir(dir string) {
	scriptsDir = dir
}

func init() {
	scriptsDir = "scripts"
}

// Runner executes the specified script with the given parameters and returns
// the output
func Runner(name string, param []string) string {
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
func HasScript(search string) bool {
	files, err := ioutil.ReadDir(scriptsDir)
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
func ListScripts() []string {
	files, err := ioutil.ReadDir(scriptsDir)
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

// RunWorker listens on jobQueue for new jobs and executes them
func RunWorker() {
	for {
		job := <-jobQueue
		job.Status = "Running"
		job.Start = time.Now()
		LogAppendLine(fmt.Sprintf("WORKER  running job %s", job.Uuid))

		cmd := exec.Command("/bin/bash", job.Script)
		cmd.Args = job.Args

		out, err := cmd.Output()
		if err != nil {
			job.Status = "runtime error"
			job.Output = fmt.Sprintf("Error occurred\n%s", err)
		} else {
			job.Status = "completed"
			job.Output = fmt.Sprintf("%s", out)
		}

		time.Sleep(time.Second * 10)

		job.Finish = time.Now()
		LogAppendLine(fmt.Sprintf("WORKER  finished job %s", job.Uuid))
		LogAppendLine(fmt.Sprintf("OUTPUT  %s", job.Output))
	}
}

// NewJob determines if a scripts can be executed and returns a job
// struct
func NewJob(script string, args []string, path string, requested time.Time) (Job, error) {
	job := Job{
		Script:  script,
		Args:    args,
		Uuid:    uuid.NewV4().String(),
		Path:    path,
		Request: time.Now()}

	if !HasScript(script) {
		err := errors.New(fmt.Sprintf("No script '%s' found in dir '%s'", script, path))
		return job, err
	}

	return job, nil
}

func init() {
	scriptsDir = "scripts"
	jobQueue = make(chan Job, 2)
}
