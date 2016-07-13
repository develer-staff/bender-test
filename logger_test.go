package main

import (
    "path/filepath"
    "fmt"
    "testing"
    "time"
    "os"
)

var test_path string
var file_name string
var job Job

func TestMain(m *testing.M) {
    now := time.Now()
    job := Job{Script:"script", Uuid: "uuid"}
    test_path, _ = filepath.Abs(filepath.Join("log", "script"))
    file_name = fmt.Sprintf("%d.%d.%d-%d.%d.%d-%s.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), job.Uuid)
    os.RemoveAll(test_path)
}

func TestWriteLog(t *testing.T){
    go WriteLog()
    job := Job{Script:"script", Uuid: "uuid"}
    jobDone <- job

    file_path:= filepath.Join(test_path, file_name)

    if _, err := os.Stat(file_path); err != nil{
        t.Error("Test failed")
    }
}

func TestReadLogDir(t *testing.T){
    expected := "{\"Script\":\"script\",\"Path\":\"\",\"Args\":null,\"Uuid\":\"uuid\"," +
                "\"Output\":\"\",\"Exit\":\"\",\"Request\":\"0001-01-01T00:00:00Z\"," +
                "\"Start\":\"0001-01-01T00:00:00Z\",\"Finish\":\"0001-01-01T00:00:00Z\"," +
                "\"Status\":\"\"}\n\n*******************\n\n"

    actual := ReadLogDir(test_path)

    if actual != expected{
        t.Error("Test failed")
    }
}

func TestReadLog(t *testing.T){
    expected := "{\"Script\":\"script\",\"Path\":\"\",\"Args\":null,\"Uuid\":\"uuid\"," +
                "\"Output\":\"\",\"Exit\":\"\",\"Request\":\"0001-01-01T00:00:00Z\"," +
                "\"Start\":\"0001-01-01T00:00:00Z\",\"Finish\":\"0001-01-01T00:00:00Z\"," +
                "\"Status\":\"\"}"

    file_path := filepath.Join(test_path, file_name)
    actual := ReadLog(file_path)

    if actual != expected{
        t.Error("Test failed")
    }
}

func TestFindLog(t *testing.T){
    expected := filepath.Join(test_path, file_name)
    actual := FindLog("uuid")

    if actual != expected{
        t.Error("Test failed")
    }

    os.RemoveAll(test_path)
}

