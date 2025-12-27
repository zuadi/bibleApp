package models

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	_ "modernc.org/sqlite"
)

type BibleDatabase struct {
	database *sql.DB
}

func NewBibleDatabase(filepath string) (bd *BibleDatabase, err error) {
	bd = &BibleDatabase{}
	//open sqlite file
	uri, err := getTranlationDBPath(filepath)
	if err != nil {
		return nil, err
	}
	bd.database, err = sql.Open("sqlite", uri+"?mode=ro&_journal_mode=DELETE")
	if err != nil {
		return nil, err
	}
	return bd, nil
}

func (bd *BibleDatabase) Close() error {
	return bd.database.Close()
}

func (bd *BibleDatabase) GetBooks() ([]*DBBook, error) {
	// look for book names and book number
	rows, err := bd.database.Query("SELECT book_number, long_name FROM books")
	if err != nil {
		return nil, err
	}

	dbBooks := []*DBBook{}

	for rows.Next() {
		b := &DBBook{}
		dbBooks = append(dbBooks, b)
		err = rows.Scan(&b.Number, &b.Name)
		if err != nil {
			return nil, err
		}
	}
	return dbBooks, nil
}

func (bd *BibleDatabase) GetVerse(bookNumber, chapterNumber, verseNumber int) (text string, err error) {
	var rows *sql.Rows
	rows, err = bd.database.Query("SELECT text FROM verses WHERE book_number = ? AND chapter = ? AND verse = ?", bookNumber, chapterNumber, verseNumber)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		err = rows.Scan(&text)
		if err != nil {
			return
		}

		// remove special characters
		tmpText := strings.NewReplacer("<pb/>", "", "<i>", "", "</i>", "", "<t>", "", "</t>", "", "<e>", "", "</e>", "", "<J>", "", "</J>", "", "[T›]", "", "<br/>", "").Replace(text)

		// remove footnotes and notes inside <> </>
	next:
		for {
			// look if footnote char is in text
			for _, value := range []string{"f", "n", "S"} {
				key := fmt.Sprintf("<%s>", value)
				if strings.Contains(tmpText, key) {
					tmpText = fmt.Sprintf("%s%s", tmpText[:strings.Index(tmpText, key)], tmpText[strings.Index(tmpText, fmt.Sprintf("</%s>", value))+len(fmt.Sprintf("</%s>", value)):])
					continue next
				}
			}
			//remove "  "
			return strings.ReplaceAll(tmpText, "  ", " "), nil
		}
	}
	return
}

func getTranlationDBPath(translation string) (string, error) {
	absPath, err := filepath.Abs(filepath.Join("bibles", translation+".SQLite3"))
	if err != nil {
		return "", err
	}

	slashPath := filepath.ToSlash(absPath)

	if runtime.GOOS == "windows" {
		slashPath = "/" + slashPath
	}

	u := &url.URL{
		Scheme:   "file",
		Path:     slashPath,
		RawQuery: "mode=ro&immutable=1",
	}

	return u.String(), nil
}
