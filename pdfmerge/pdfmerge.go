package pdfmerge

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
)

var size, margin string
var scaleH, scaleW, landscape, verbose, version bool
var JPEGQuality int

const (
	DefaultSize   = "IMG-SIZE"
	DefaultMargin = "0,0,0,0"
	VERSION       = "1.2.0"
)

func Pdfmerge(inputpath []string, outputpath string) {

	if verbose {
		unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
	}

	c := creator.New()

	for _, arg := range inputpath {
		err := NewSource(arg).MergeTo(c)
		if err != nil {
			fmt.Printf("Error: %s (%s) \n", err.Error(), arg)
			os.Exit(1)
		}
	}

	err := c.WriteToFile(outputpath)
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
	}

	debugInfo(fmt.Sprintf("Complete, see output file: %s", outputpath))
}
