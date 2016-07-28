package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type JobStatus string

const (
	JOB_QUEUED     = "queued"
	JOB_NOT_FOUND  = "not found"
	JOB_QUEUE_FULL = "queue full"
	JOB_WORKING    = "working"
	JOB_FAILED     = "failed"
	JOB_COMPLETED  = "completed"
)

type Job struct {
	Name    string
	Params  []string
	Uuid    string
	Created time.Time
	Status  JobStatus
	Timeout int
}

type JobInterface interface {
	Run(name, uuid string, args []string) int
	State()
}

var scriptsDir string
var run bool

func GetScriptsDir() string {
	return scriptsDir
}

func SetScriptsDir(dir string) {
	scriptsDir = dir
}

func init() {
	SetScriptsDir("scripts")
	logContextRunner = LoggerContext{
		name:  "RUNNER",
		level: 3}
}

func FakeRun(job *Job, script, uuid string, args []string) int {
	job.Name = script
	job.Uuid = uuid
	job.Params = args
	job.Status = JOB_WORKING

	var exit int

	if FakeHasScript(job.Name) {
		run = true
		go func() {
			time.Sleep(3 * time.Second)
			//execution...
			run = false
		}()
		exit = 0
	} else {
		exit = -1
	}

	return exit
}

//Check if a script exists
func FakeHasScript(script string) bool {
	return true
}

//Return the current stdout and stderr
func FakeLog(job *Job) string {
	buf := make([]byte, 100)
	//reading from stdout pipe
	return string(buf)
}

//Handle the status of script
func FakeState(job *Job) {
	if run {
		job.Status = JOB_WORKING
	} else {
		job.Status = JOB_COMPLETED
	}

}

var cmd = exec.Command("")
var outChan = make(chan string, 1)
var syncChan = make(chan bool)
var endReadStart = make(chan bool)
var logContextRunner LoggerContext

//Initialize the script command
func (job *Job) Run(script, uuid string, args []string) int {
	job.Name = script
	job.Uuid = uuid
	job.Params = args
	job.Status = JOB_WORKING

	var exit int

	if name, exist := HasScript(job.Name); exist {
		script_path := filepath.Join(GetScriptsDir(), name)

		cmd = exec.Command(script_path, job.Params...)
		go Start()
		exit = 0
	} else {
		LogErr(logContextRunner, "Script does not exist")
		exit = -1
	}

	return exit
}

//Run the command
func Start() {
	<-syncChan
	time.Sleep(100 * time.Millisecond)
	cmd.Start()
	LogInf(logContextRunner, "Execution started...")
	<-endReadStart
	err := cmd.Wait()
	LogInf(logContextRunner, "Execution finished")

	if err != nil {
		LogErr(logContextRunner, "Error occurred during execution")
	}
}

//Check if a script exists
func HasScript(script string) (string, bool) {
	files, err := ioutil.ReadDir(GetScriptsDir())
	var exist = false
	var name = ""

	if err != nil {
		LogErr(logContextRunner, "No scripts directory found")
	} else {
		for _, file := range files {
			if strings.Contains(file.Name(), script) {
				name = file.Name()
				exist = true
			}
		}
	}
	return name, exist
}

//Return the current stdout and stderr
func Log() *chan string {
	go func() {
		syncChan <- true
		stdout, err := cmd.StdoutPipe()
		stderr, err := cmd.StderrPipe()
		multi := io.MultiReader(stdout, stderr)
		scanner := bufio.NewScanner(multi)

		if err != nil {
			LogErr(logContextRunner, "Error occurred while reading stdout/stderr")
		}

		for scanner.Scan() {
			out := scanner.Text()
			outChan <- out
		}

		endReadStart <- true
		LogDeb(logContextRunner, "finished reading, sent sync to chan 'endReadStart'")
		endReadLog <- true
		LogDeb(logContextRunner, "finished reading, sent sync to chan 'endReadLog'")
	}()

	return &outChan
}

//Handle the status of script
func (job *Job) State() {
	if cmd.ProcessState == nil {
		job.Status = JOB_WORKING
	} else if cmd.ProcessState.Success() {
		job.Status = JOB_COMPLETED
	} else {
		job.Status = JOB_FAILED
	}
}

func List() []string {
	files, err := ioutil.ReadDir(GetScriptsDir())
	var scripts []string

	if err != nil {
		LogErr(logContextRunner, "No scripts directory found")
	} else {
		for _, file := range files {
			n := strings.LastIndexByte(file.Name(), '.')
			if n > 0 {
				scripts = append(scripts, file.Name()[:n])
			} else {
				scripts = append(scripts, file.Name())
			}
		}
	}
	return scripts
}
