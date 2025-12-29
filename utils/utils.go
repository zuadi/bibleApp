package utils

import (
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
	switch runtime.GOOS {
	case "darwin":
		exe, _ := os.Executable()
		return filepath.Join(exe, "..", "..", "Resources", path)
	case "linux":
		if appDir := os.Getenv("APPDIR"); appDir != "" {
			return filepath.Join(appDir, "usr", "share", "bibletool")
		}
	}
	return path
}
