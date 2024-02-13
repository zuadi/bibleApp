// windowscmd.go
//go:build windows
// +build windows

package basic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

var WaitingQue int

func StartcommandwithArgs(command string, Args []string) {

	var debug bool
	var stderr, stdout io.ReadCloser
	var err error

	cmd := exec.Command(command, Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	// some command output will be input into stderr
	if debug {
		stderr, err = cmd.StderrPipe()
		if err != nil {
			fmt.Println(err)
		}

		stdout, err = cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	if debug {
		// print the output of the subprocess
		scanner := bufio.NewScanner(stdout)
		firstline := false
		for scanner.Scan() {
			m := scanner.Text()

			//skip first line with call of codesys information
			if firstline {
				fmt.Println(m)
			} else if m != "" {
				firstline = true
			}

		}

		// print the output of the subprocess
		scanner = bufio.NewScanner(stderr)
		firstline = false

		for scanner.Scan() {
			m := scanner.Text()
			if m != "" {
				//print empty line
				if !firstline {
					fmt.Println(" ")
					firstline = true
				}
				fmt.Println(m)
			}
		}
	}
	cmd.Wait()
}

func StartcommandwithArgsWG(command string, Args []string, wg *sync.WaitGroup) {
	defer wg.Done()
	var debug bool
	var stderr, stdout io.ReadCloser
	var err error

	WaitingQue += 1

	cmd := exec.Command(command, Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	// some command output will be input into stderr
	if debug {
		stderr, err = cmd.StderrPipe()
		if err != nil {
			fmt.Println(err)
		}

		stdout, err = cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	if debug {
		// print the output of the subprocess
		scanner := bufio.NewScanner(stdout)
		firstline := false
		for scanner.Scan() {
			m := scanner.Text()

			//skip first line with call of codesys information
			if firstline {
				fmt.Println(m)
			} else if m != "" {
				firstline = true
			}

		}

		// print the output of the subprocess
		scanner = bufio.NewScanner(stderr)
		firstline = false

		for scanner.Scan() {
			m := scanner.Text()
			if m != "" {
				//print empty line
				if !firstline {
					fmt.Println(" ")
					firstline = true
				}
				fmt.Println(m)
			}
		}
	}
	cmd.Wait()

	WaitingQue -= 1

}

func Getexedir() string {
	//current working directory
	curdir, err := os.Getwd()
	CheckErr(err, "Error findong current working directory")
	pathseparator := string(os.PathSeparator)

	return strings.Join([]string{curdir, pathseparator}, "")
}
