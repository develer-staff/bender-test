package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/satori/go.uuid"
)

// NewJob builds a job struct and returns its pointer.
func NewJob(a *appContext, script string, args []string, requested time.Time) *Job {
	exists, scriptName := HasScript(a, script)
	path, err := filepath.Abs(a.ScriptsDir)

	if err != nil {
		LogErrors(err)
		panic(err.Error())
	}

	job := &Job{
		Script:  script,
		Path:    filepath.Join(path, scriptName),
		Args:    args,
		Uuid:    uuid.NewV4().String(),
		Request: requested}

	if !exists {
		job.Status = JOB_NOT_FOUND
	}

	return job
}

// RunWorker
func RunWorker(a *appContext) {
	LogAppendLine("WORKER  ready")
	for job := range a.JobQueue {
		LogAppendLine(fmt.Sprintf("WORKER  running job %s", job.Uuid))
		job.Status = JOB_WORKING
		job.Start = time.Now()

		cmdOutput := &bytes.Buffer{}
		cmd := exec.Command(job.Path, job.Args...)
		cmd.Stdout = cmdOutput
		err := cmd.Run()

		if err != nil {
			LogAppendLine(fmt.Sprintf("WORKER  failed job %s", job.Uuid))
			LogErrors(err)
			job.Exit = err.Error()
			job.Status = JOB_FAILED
		} else {
			LogAppendLine(fmt.Sprintf("WORKER  completed job %s", job.Uuid))
			job.Status = JOB_COMPLETED
		}

		job.Finish = time.Now()
		job.Output = string(cmdOutput.Bytes())

		// TODO send job to jobDone channel
	}
}

// Runner executes the specified script with the given parameters and returns
// the output
func Runner(a *appContext, name string, param []string) string {
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

// hasScript looks for a script in the default script dir and returns a bool
// and the first matching filename for the script (in alphabetical order).
func HasScript(a *appContext, script string) (bool, string) {
	files, err := ioutil.ReadDir(a.ScriptsDir)
	if err != nil {
		LogErrors(err)
		return false, ""
	}
	for _, file := range files {
		if file.Name() == script {
			return true, file.Name()
		}
		filename := file.Name()[0 : len(file.Name())-len(filepath.Ext(file.Name()))]
		if filename == script {
			return true, file.Name()
		}
	}
	return false, ""
}

// listScripts returns a list of scripts in the default script folder
func ListScripts(a *appContext) []string {
	files, err := ioutil.ReadDir(a.ScriptsDir)
	if err != nil {
		LogErrors(err)
		panic(err.Error())
	}
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names
}
