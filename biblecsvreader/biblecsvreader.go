package biblecsvreader

import (
	"bibletool/basic"
	"encoding/csv"
	"os"
)

type Bibleindex struct {
	CSVData    [][]string
	Bibletrans []string
}

func ReadCSV(path *basic.OSPaths) *Bibleindex {
	var bi Bibleindex
	// open csv config file
	f, err := os.Open(path.CsvPath)
	basic.CheckErr(err, "Error open Bibleindex.csv file")
	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	bi.CSVData, err = csvReader.ReadAll()

	basic.CheckErr(err, "Error could not read csv file")
	// this function read all avaiable translation from csv config file
	bi.Bibletrans = bi.CSVData[0][1:]

	return &bi
}
