package bibletool

import (
	"bibletool/bibletool/env"
	"bibletool/bibletool/models"
	"encoding/csv"
	"os"
	"path/filepath"
	"sync"

	"gitea.tecamino.com/paadi/tecamino-logger/logging"
	"github.com/skratchdot/open-golang/open"
)

type Bibletool struct {
	config           *models.Config
	Logger           *logging.Logger
	OutputDir        string
	TempDir          string
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

	//get user dir for Bibletranslation folder
	bt.OutputDir = env.OutputDir.GetValue()
	if bt.OutputDir == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			bt.LogError("get user home dir", err)
			return nil, err
		}
		bt.OutputDir = filepath.Join(homedir, "Bibletranslation")
	}

	//load enviroment variables
	if err := env.Load(".env"); err != nil {
		bt.LogError("load enviroment variables", err)
	}

	//read config if file exists
	err = bt.LoadUserSettings()
	if err != nil {
		bt.LogError("load user config", err)
	}
	return
}

func (bt *Bibletool) GetAllTranslations() (list []string, err error) {
	f, err := os.Open(env.BibleIndexFile.GetValue())
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
	if err := open.Run(bt.OutputDir); err != nil {
		bt.LogError("open file browser", err)
		return err
	}

	// save checkbox settings for next start
	err := bt.SaveUserSettings()
	if err != nil {
		return err
	}
	return nil
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
