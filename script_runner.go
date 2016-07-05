package main

import (
	"fmt"
	"os/exec"
)

func Runner(name string, param []string) string {
	/*Execute the specified script with is parameters
	  and return the output*/
	cmd := exec.Command("bash", name)
	var output string
	for _, item := range param {
		cmd.Args = append(cmd.Args, item)
	}
	out, err := cmd.Output()
	if err != nil {
		output = fmt.Sprintf("Error occurred\n%s", err)
	} else {
		output = fmt.Sprintf("%s", out)
	}
	return output
}
