package main

import (
	"fmt"
	"os/exec"
)

func runner(name string, param []string) string {
	cmd := exec.Command("bash", name)
	output := ""

	for _, item := range param {
		cmd.Args = append(cmd.Args, item)
	}
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error occurred")
		output = fmt.Sprintf("%s", err)
	} else {
		output = fmt.Sprintf("%s", out)
	}
	return output
}
