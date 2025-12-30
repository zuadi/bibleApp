package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func MkDirs(root string, subfolder ...string) error {
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		if err := os.RemoveAll(root); err != nil {
			return err
		}
	}

	if err := os.Mkdir(root, 0755); err != nil {
		return err
	}

	//create subfolders
	for _, sub := range subfolder {
		if err := os.Mkdir(filepath.Join(root, sub), 0755); err != nil {
			return err
		}
	}

	return nil
}

func GetDistOsPath(path string) string {
	if runtime.GOOS == "darwin" {
		exe, _ := os.Executable()
		fmt.Println(10, exe)
		return filepath.Join(exe, "..", "..", "Resources", path)
	}
	return path
}
