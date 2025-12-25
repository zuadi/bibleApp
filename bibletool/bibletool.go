package bibletool

import (
	"bibletool/basic"
	"bibletool/bibletool/consts"
	"bibletool/bibletool/models"
	"encoding/csv"
	"os"
	"sync"

	"gitea.tecamino.com/paadi/tecamino-logger/logging"
	"github.com/skratchdot/open-golang/open"
)

type Bibletool struct {
	config           *models.Config
	Logger           *logging.Logger
	OsPaths          *basic.OSPaths
	bibleIndex       models.BibleIndex
	Wg               sync.WaitGroup
	TotalProgress    func(process float64)
	DocumentProgress func(title string, process float64)
	PdfProgress      func(process float64)
}

func NewBibletool() (bt *Bibletool, err error) {
	bt = &Bibletool{}
	bt.Logger, err = logging.NewLogger("", logging.DefaultConfig())
	if err != nil {
		return nil, err
	}

	//call basic function
	bt.OsPaths, err = basic.GetOSPaths()
	if err != nil {
		bt.LogError("get OS paths", err)
	}

	//read config if file exists
	err = bt.LoadUserSettings()
	if err != nil {
		bt.LogError("load user config", err)
	}
	return
}

func (bt *Bibletool) GetAllTranslations() (list []string, err error) {
	f, err := os.Open(consts.BibleIndexFile)
	if err != nil {
		bt.LogError("open csv file", err)
		return nil, err
	}
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	data, err := csvReader.Read()
	if err != nil {
		bt.LogError("read csv file", err)
		return nil, err
	}
	return data[1:], nil
}

func (bt *Bibletool) Close() error {
	// open file browser of translation
	if err := open.Run(bt.OsPaths.Outputpath); err != nil {
		bt.LogError("open file browser", err)
		return err
	}

	// save checkbox settings for next start
	err := bt.SaveUserSettings()
	if err != nil {
		return err
	}
	//remove temporary folder
	return basic.Deltemp(bt.OsPaths.Tempdir)
}

func (bt *Bibletool) TotalProgressAdd(add int) {
	if bt.TotalProgress != nil {
		bt.TotalProgress(float64(add))
	}
}

func (bt *Bibletool) DocumentProgressAdd(title string, add int) {
	if bt.DocumentProgress != nil {
		bt.DocumentProgress(title, float64(add))
	}
}

var processing float64
var CancelProgress bool

func (bt *Bibletool) PdfProgressAdd(add int) {
	///TODO:
	// if processing == 0.0 {
	// 	go func() {
	// 		for {
	// 			if CancelProgress {
	// 				break
	// 			}
	// 			bt.PdfProgress(0.01)
	// 			time.Sleep(500 * time.Millisecond)
	// 		}
	// 	}()
	// }
	if bt.PdfProgress != nil {
		bt.PdfProgress(float64(add))
	}
}
