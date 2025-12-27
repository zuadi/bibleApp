package bibletool

import (
	"bibletool/bibletool/env"
	"bibletool/bibletool/models"
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gitea.tecamino.com/paadi/tecamino-logger/logging"
	"github.com/skratchdot/open-golang/open"
)

type Bibletool struct {
	AppName          string
	config           *models.Config
	Logger           *logging.Logger
	OutputDir        string
	TempDir          string
	AbsIconPath      string
	bibleIndex       models.BibleIndex
	Wg               sync.WaitGroup
	TotalProgress    func(process float64)
	DocumentProgress func(title string, process float64)
	PdfProgress      func(process float64)
}

func NewBibletool() (bt *Bibletool, err error) {
	bt = &Bibletool{}

	//load enviroment variables
	if err := env.Load(".env"); err != nil {
		bt.LogError("load enviroment variables", err)
	}

	bt.DebugLog("main", "read app name")
	bt.AppName = env.AppName.GetValue()
	bt.DebugLog("main", bt.AppName)

	logConfig := logging.DefaultConfig()
	logConfig.Debug = strings.ToLower(env.Debug.GetValue()) == "true" || env.Debug.GetValue() == "1"

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	bt.Logger, err = logging.NewLogger(filepath.Join(configDir, bt.AppName+".log"), logConfig)
	if err != nil {
		return nil, err
	}

	//get user dir for Bibletranslation folder
	bt.DebugLog("NewBibletool", "set output directory")
	bt.OutputDir = env.OutputDir.GetValue()
	if bt.OutputDir == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			bt.LogError("get user home dir", err)
			return nil, err
		}
		bt.OutputDir = filepath.Join(homedir, "Bibletranslation")
	}
	bt.DebugLog("NewBibletool", bt.OutputDir)

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	//set icon path for html template
	bt.AbsIconPath = filepath.Join(workingDir, env.IconPath.GetValue())

	//read config if file exists
	err = bt.LoadUserSettings()
	if err != nil {
		bt.LogError("load user config", err)
	}
	return
}

func (bt *Bibletool) GetAllTranslations() (list []string, err error) {
	bt.DebugLog("GetAllTranslations", "load bibleindex from file "+env.BibleIndexFile.GetValue())
	f, err := os.Open(env.BibleIndexFile.GetValue())
	if err != nil {
		bt.LogError("open csv file", err)
		return nil, err
	}
	defer f.Close()

	bt.DebugLog("GetAllTranslations", "read as csv")
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
	bt.DebugLog("Close", "open file explorer "+bt.OutputDir)
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

func (bt *Bibletool) PdfProgressAdd(add int) {
	if bt.PdfProgress != nil {
		bt.PdfProgress(float64(add))
	}
}
