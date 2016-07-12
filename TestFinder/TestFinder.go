// TestFinder
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


var PATHSTRING string = filepath.Join("C:\\", "Stage", "robotester", "Scripts")

func check(e error) bool {
	if e != nil {
		panic(e)
		return false
	}
	return true
}


// list
// return names: a list of the names of the scripts in the directory of the 
// scripts, a void list otherwise
func list() []string {
	files, err := ioutil.ReadDir(PATHSTRING)
	if check(err) {
		names := []string{}
		k := len(files)
		for i := 0; i < k; i++ {
			names = append(names, files[i].Name())
		}
		return names
	}
	return []string
}

// findFile
// parameters search: a string that should be in the name of the script
// returns true: if there is at least one file that fulfills the requests false 
// otherwise
func findFile(search string) bool {
	files, err := ioutil.ReadDir(PATHSTRING)
	if check(err) {
		k := len(files)
		for i := 0; i < k; i++ {
			if strings.Contains(files[i].Name(), search) {
				return true
			}
		}
	}
	return false
}

func main() {
	scripts := list()
	fmt.Println(scripts)
	result := findFile(os.Args[1:])
	fmt.Println(result)
}
