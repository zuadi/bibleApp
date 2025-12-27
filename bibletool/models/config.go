package models

type Config struct {
	Maintransation string   `json:"mainTranslation"`
	Translations   []string `json:"translation"`
	SameDocument   bool     `json:"sameDocument"`
	SermonTitle    string   `json:"sermonTitle"`
	Pastor         string   `json:"pastor"`
	Verses         string   `json:"verses"`
}
