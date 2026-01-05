package bibletool

import (
	"bibletool/bibletool/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (bt *Bibletool) LoadUserSettings() error {
	bt.config = &models.Config{}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(configDir, bt.AppName)
	os.MkdirAll(appDir, 0755)

	configFile := filepath.Join(appDir, "config.json")
	bt.DebugLog("NewBibletool", "load user setting from file "+configFile)

	if _, err := os.Stat(configFile); err == nil {
		raw, err := os.ReadFile(configFile)
		if err != nil {
			bt.LogError("load user settings", err)
			return err
		}
		err = json.Unmarshal(raw, &bt.config)
		if err != nil {
			bt.LogError("json unmarshal user config", err)
			return err
		}
	}
	return nil
}

func (bt *Bibletool) SaveUserSettings() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(configDir, bt.AppName)
	os.MkdirAll(appDir, 0755)

	configFile := filepath.Join(appDir, "config.json")

	b, err := json.Marshal(bt.config)
	if err != nil {
		bt.LogError("json marshal user config", err)
		return err
	}
	err = os.WriteFile(configFile, b, 0644)
	if err != nil {
		bt.LogError("save user config", err)
		return err
	}
	return nil
}

func (bt *Bibletool) SetMaintranslation(t string) {
	bt.DebugLog("SetMaintranslation", t)
	bt.config.Maintransation = t
}

func (bt *Bibletool) GetMaintranslation() string {
	bt.DebugLog("GetMaintranslation", bt.config.Maintransation)
	return bt.config.Maintransation
}

func (bt *Bibletool) SetTranslations(t []string) {
	bt.DebugLog("SetTranslations", fmt.Sprint(t))
	bt.config.Translations = t
}

func (bt *Bibletool) GetSelectedTranslations() []string {
	bt.DebugLog("GetSelectedTranslations", fmt.Sprint(bt.config.Translations))
	return bt.config.Translations
}

func (bt *Bibletool) FilteredTranslations() (list []string) {
	bt.DebugLog("FilteredTranslations", "filter out main translation from list")
	for _, translation := range bt.config.Translations {
		if translation == bt.config.Maintransation {
			continue
		}
		list = append(list, translation)
	}
	return
}

func (bt *Bibletool) SetSameDocument(b bool) {
	bt.DebugLog("SetSameDocument", fmt.Sprint(b))
	bt.config.SameDocument = b
}

func (bt *Bibletool) GetSameDocument() bool {
	bt.DebugLog("GetSameDocument", fmt.Sprint(bt.config.SameDocument))
	return bt.config.SameDocument
}

func (bt *Bibletool) SetPastor(s string) {
	bt.DebugLog("SetPastor", s)
	bt.config.Pastor = s
}

func (bt *Bibletool) GetPastor() string {
	bt.DebugLog("SetPastor", bt.config.Pastor)
	return strings.TrimSpace(bt.config.Pastor)
}

func (bt *Bibletool) SetSermonTitle(s string) {
	bt.DebugLog("SetPastor", s)
	bt.config.SermonTitle = s
}

func (bt *Bibletool) GetSermonTitle() string {
	bt.DebugLog("GetSermonTitle", bt.config.SermonTitle)
	return strings.TrimSpace(bt.config.SermonTitle)
}

func (bt *Bibletool) SetOutputFile(s string) {
	bt.DebugLog("SetOutputFile", s)
	bt.config.OutputFile = s
}

func (bt *Bibletool) GetOutputFile() string {
	bt.DebugLog("GetOutputFile", bt.config.OutputFile)
	return strings.TrimSpace(bt.config.OutputFile)
}

func (bt *Bibletool) SetVerses(s string) {
	bt.DebugLog("SetVerses", s)
	bt.config.Verses = s
}

func (bt *Bibletool) GetVerses() string {
	bt.DebugLog("GetVerses", bt.config.Verses)
	return bt.config.Verses
}
