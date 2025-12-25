package env

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	Env              EnvKey = "ENV"
	Debug            EnvKey = "DEBUG"
	ChromePath       EnvKey = "CHROME_PATH"
	AppName          EnvKey = "APP_NAME"
	IconPath         EnvKey = "ICON_FILE"
	BibleIndexFile   EnvKey = "BIBLEINDEX_FILE"
	ConfigFile       EnvKey = "CONFIG_FILE"
	HtmlTemplateFile EnvKey = "HTML_TEMPLATE_FILE"
	OutputDir        EnvKey = "OUTPUT_DIR"
)

type EnvKey string

func Load(path string) error {
	if path == "" {
		path = ".env"
	}
	return godotenv.Load(path)
}

func (key EnvKey) GetValue() string {
	return os.Getenv(string(key))
}
