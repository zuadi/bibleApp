package bibletool

import (
	"bibletool/bibletool/env"
	"bibletool/bibletool/models"
	"encoding/csv"
	"os"
)

func (bt *Bibletool) ReadBibleIndexes() (models.BibleIndex, error) {
	path := env.BibleIndexFile.GetValue()
	bt.DebugLog("ReadBibleIndexes", "read bibleindex file "+path)

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
