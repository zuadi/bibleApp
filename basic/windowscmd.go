// windowscmd.go
//go:build windows
// +build windows

package basic

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

var WaitingQue int

func StartcommandwithArgs(command string, Args []string) (err error) {
	cmd := exec.Command(command, Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}

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
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}

	err = cmd.Start()
	if err != nil {
		return
	}

	cmd.Wait()

	WaitingQue -= 1
	return
}

func Getexedir() (string, error) {
	//current working directory
	curdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/", curdir), nil
}
