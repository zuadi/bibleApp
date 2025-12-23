package bibletool

import (
	"bibletool/bibletool/consts"
	"bibletool/bibletool/models"
	"bytes"
	"encoding/gob"
	"os"
	"strings"
)

func (bt *Bibletool) LoadUserSettings() error {
	bt.config = &models.Config{}

	if _, err := os.Stat(consts.ConfigPath); err == nil {
		raw, err := os.ReadFile(consts.ConfigPath)
		if err != nil {
			bt.LogError("load user settings", err)
			return err
		}
		buffer := bytes.NewBuffer(raw)
		dec := gob.NewDecoder(buffer)
		err = dec.Decode(&bt.config)
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			}
			bt.LogError("load user settings", err)
			return err
		}
	}
	return nil
}

func (bt *Bibletool) SaveUserSettings() error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(*bt.config)
	if err != nil {
		bt.LogError("save user config", err)
		return err
	}
	err = os.WriteFile(consts.ConfigPath, buffer.Bytes(), 0644)
	if err != nil {
		bt.LogError("save user config", err)
		return err
	}
	return nil
}

func (bt *Bibletool) SetMaintranslation(t string) {
	bt.config.Maintransation = t
}

func (bt *Bibletool) GetMaintranslation() string {
	return bt.config.Maintransation
}

func (bt *Bibletool) SetTranslations(t []string) {
	bt.config.Translations = t
}

func (bt *Bibletool) GetSelectedTranslations() []string {
	return bt.config.Translations
}

func (bt *Bibletool) FilteredTranslations() (list []string) {
	for _, translation := range bt.config.Translations {
		if translation == bt.config.Maintransation {
			continue
		}
		list = append(list, translation)
	}
	return
}

func (bt *Bibletool) SetSameDocument(b bool) {
	bt.config.SameDocument = b
}

func (bt *Bibletool) GetSameDocument() bool {
	return bt.config.SameDocument
}

func (bt *Bibletool) SetPastor(s string) {
	bt.config.Pastor = s
}

func (bt *Bibletool) GetPastor() string {
	return strings.TrimSpace(bt.config.Pastor)
}

func (bt *Bibletool) SetSermonTitle(s string) {
	bt.config.SermonTitle = s
}

func (bt *Bibletool) GetSermonTitle() string {
	return strings.TrimSpace(bt.config.SermonTitle)
}

func (bt *Bibletool) SetVerses(s string) {
	bt.config.Verses = s
}

func (bt *Bibletool) GetVerses() string {
	return bt.config.Verses
}
