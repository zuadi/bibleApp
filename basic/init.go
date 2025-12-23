package basic

import (
	"os"
	"path/filepath"
)

type OSPaths struct {
	Currentdirectory string
	HtmlTemplatePath string
	Outputpath       string
	Tempdir          string
}

func GetOSPaths() (*OSPaths, error) {

	//get executable dir
	curdir, err := Getexedir()
	if err != nil {
		return nil, err
	}

	//get user dir for Bibletranslation folder
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}

	return &OSPaths{
		Currentdirectory: curdir,
		Outputpath:       filepath.Join(homedir, "Bibletranslation"),
		Tempdir:          tempDir}, nil
}

// delete given path/folder
func Deltemp(path string) error {
	return os.RemoveAll(path)
}

// remove slice element out of slice
func Rmsliceelement(slice []string, index int) []string {
	copy(slice[index:], slice[index+1:]) // shift valuesafter the indexwith a factor of 1
	slice[len(slice)-1] = ""             // remove element
	slice = slice[:len(slice)-1]         // truncateslice
	return slice
}
