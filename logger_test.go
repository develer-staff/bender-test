package main

import (
    "testing"
    "fmt"
    "path/filepath"
    "time"
    "os"
)

var test_path string
var file_name string
var job Job

func TestMain(m *testing.M) {
    now := time.Now()
    job = Job{Script:"script", Uuid: "uuid"}
    test_path, _ = filepath.Abs(filepath.Join("log", "script"))
    file_name = fmt.Sprintf("%d.%d.%d-%d.%d.%d-%s.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), job.Uuid)
    code := m.Run()
    os.RemoveAll(test_path)
    os.Exit(code)
}

func TestWriteLog(t *testing.T){
    go WriteLog()
    jobDone <- job
    file_path:= filepath.Join(test_path, file_name)

    time.Sleep(20 * time.Millisecond)
    if _, err := os.Stat(file_path); os.IsNotExist(err){
        t.Error("Test failed")
    }
}