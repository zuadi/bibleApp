package basic

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type OSPaths struct {
	Pathseperator    string
	Currentdirectory string
	Outputpath       string
	Tempdir          string
	Configpath       string
	CsvPath          string
	IconPath         string
}

func Settings() *OSPaths {

	//get pathseparator from OS
	pathsep := string(os.PathSeparator)

	//get executable dir
	curdir := Getexedir()

	//get user dir for Bibletranslation folder
	homedir, err := os.UserHomeDir()
	CheckErr(err, "Error could not find user home dir")

	return &OSPaths{
		Pathseperator:    pathsep,
		Currentdirectory: curdir,
		Outputpath:       strings.Join([]string{homedir, pathsep, `Bibletranslation`}, ""),
		Tempdir:          CreateTempdir(),
		Configpath:       strings.Join([]string{curdir, "BibletoolConfig"}, ""),
		CsvPath:          strings.Join([]string{curdir, "cfg", pathsep, "Bibleindex.csv"}, ""),
		IconPath:         strings.Join([]string{curdir, `pics`, pathsep, `pottershouse.png`}, "")}
}

func CreateTempdir() string {

	//create temp folder
	tempdir, err := ioutil.TempDir("", "")
	if err != nil {
		CheckErr(err, "Error making temp dir")
	}
	return tempdir
}

// delete given path/folder
func Deltemp(path string) {
	err := os.RemoveAll(path)
	CheckErr(err, strings.Join([]string{"Error Removing: ", path}, ""))
}

// Error handling print out error message with a definded string
func CheckErr(err error, info string) {
	// error handling for sqlite file
	if err != nil {
		fmt.Println(info)

		f, err1 := os.OpenFile("Bibletool.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err1 != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
		fmt.Println(info)
		log.Println(info)
		log.Println(err)

	}

}

// remove slice element out of slice
func Rmsliceelement(slice []string, index int) []string {
	copy(slice[index:], slice[index+1:]) // shift valuesafter the indexwith a factor of 1
	slice[len(slice)-1] = ""             // remove element
	slice = slice[:len(slice)-1]         // truncateslice
	return slice
}
