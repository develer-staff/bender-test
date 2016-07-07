package main

import (
	"os"
	"path/filepath"
	"testing"
)

var exp_res_runner string = "Hello World\n"

func TestRunner(t *testing.T) {
	str := []string{"Hello", "World"}
	wd, _ := os.Getwd()
	fpath := filepath.Join(wd, "scripts", "hello.sh")
	out := Runner(fpath, str)
	if exp_res_runner != out {
		t.Error("Expected Hello World, got", out)
	}
}
