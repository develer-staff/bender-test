package main

import (
	"bytes"
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
var jobDone chan Job

func GetScriptsDir() string {
	return scriptsDir
}

func SetScriptsDir(dir string) {
	scriptsDir = dir
}

// hasScript checks for the script existance
func HasScript(search string) (bool, string) {
	files, err := ioutil.ReadDir(scriptsDir)
	if err != nil {
		return false, ""
	}
	for _, file := range files {
		namefile := file.Name()[0 : len(file.Name())-len(filepath.Ext(file.Name()))]
		if namefile == search {
			return true, file.Name()
		}
	}
	return false, ""
}

// listScripts returns a list of scripts in the default script folder
func ListScripts() []string {
	files, err := ioutil.ReadDir(scriptsDir)
	if err != nil {
		return nil
	}
	names := []string{}
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
}

// RunWorker listens on jobQueue for new jobs and executes them
func RunWorker() {
	for {
		job := <-jobQueue
		job.Status = "Running"
		job.Start = time.Now()
		LogAppendLine(fmt.Sprintf("WORKER  running job %s", job.Uuid))

		cmd := exec.Command(job.Path, job.Args...)
		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput
		err := cmd.Start()

		if err != nil {
			job.Output = err.Error()
			job.Status = "Runtime error"
		} else {
			cmd.Wait()
			job.Output = string(cmdOutput.Bytes())
			job.Status = "Completed"
		}

		job.Finish = time.Now()
		LogAppendLine(fmt.Sprintf("WORKER  finished job %s", job.Uuid))
		LogAppendLine(fmt.Sprintf("OUTPUT  %s", job.Output))
		jobDone <- job
	}
}

// NewJob determines if a scripts can be executed and returns a job
// struct
func NewJob(script string, args []string, path string, requested time.Time) (Job, error) {
	job := Job{
		Script:  script,
		Args:    args,
		Uuid:    uuid.NewV4().String(),
		Request: requested}

	check, name := HasScript(script)
	if !check {
		err := errors.New(fmt.Sprintf("No script '%s' found in dir '%s'", script, path))
		return job, err
	}

	job.Path = filepath.Join(path, name)

	return job, nil
}

func init() {
	scriptsDir, _ = filepath.Abs("scripts")
	jobQueue = make(chan Job, 2)
	jobDone = make(chan Job)
}
