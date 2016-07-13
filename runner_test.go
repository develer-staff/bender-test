package main

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestListScripts(t *testing.T) {
	expected := []string{"fail.sh", "foo.sh", "sleep.sh"}
	actual := ListScripts()

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Test failed")
	}
}

func TestHasScriptFoo(t *testing.T) {
	expectedBool := true
	expectedString := "foo.sh"

	actualBool, actualString := HasScript("foo")

	if actualBool != expectedBool {
		t.Error("Test failed")
	}
	if actualString != expectedString {
		t.Error("Test failed")
	}
}

func TestHasScriptSleep(t *testing.T) {
	expectedBool := true
	expectedString := "sleep.sh"

	actualBool, actualString := HasScript("sleep")

	if actualBool != expectedBool {
		t.Error("Test failed")
	}
	if actualString != expectedString {
		t.Error("Test failed")
	}
}

func TestNewJob(t *testing.T) {
	now := time.Now()
	expected := Job{
		Script:  "foo",
		Args:    []string{"-b", "--ar"},
		Path:    filepath.Join(scriptsDir, "foo.sh"),
		Request: now}
	actual, err := NewJob("foo", []string{"-b", "--ar"}, scriptsDir, now)

	if err != nil {
		t.Error("Test failed")
	}
	if actual.Script != expected.Script {
		t.Error("Test failed")
	}
	if !reflect.DeepEqual(actual.Args, expected.Args) {
		t.Error("Test failed")
	}
	if actual.Path != expected.Path {
		t.Error("Test failed")
	}
	if actual.Request != expected.Request {
		t.Error("Test failed")
	}
}
