package biblecsvreader

import (
	"encoding/csv"
	"os"
)

type BibleData [][]string

func ReadCSV(path string) (BibleData, error) {
	// open csv config file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'

	return csvReader.ReadAll()
}
