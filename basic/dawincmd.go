// darwincmd.go
//go:build darwin
// +build darwin

package basic

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var WaitingQue int

func StartcommandwithArgs(command string, Args []string) (err error) {
	cmd := exec.Command(command, Args...)

	err = cmd.Start()
	if err != nil {
		return
	}

	cmd.Wait()
	return
}

func StartcommandwithArgsWG(command string, Args []string, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	WaitingQue += 1

	cmd := exec.Command(command, Args...)

	err = cmd.Start()
	if err != nil {
		return
	}

	cmd.Wait()

	WaitingQue -= 1
	return
}

func Getexedir() string {
	//current working directory
	curdir := filepath.Dir(os.Args[0])
	return fmt.Sprintf("%s/", curdir)
}
