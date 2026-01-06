package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func ImageToBase64(path string) (string, error) {
	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(path))
	var mime string
	switch ext {
	case ".png":
		mime = "image/png"
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".gif":
		mime = "image/gif"
	default:
		return "", fmt.Errorf("unsupported image type: %s", ext)
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(imgBytes), nil
}

func GetDistOsPath(path string) string {
	if runtime.GOOS == "darwin" {
		exe, _ := os.Executable()
		fmt.Println(10, exe)
		return filepath.Join(exe, "..", "..", "Resources", path)
	}
	return path
}
