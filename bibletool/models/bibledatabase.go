package models

import (
	"C"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type BibleDatabase struct {
	database *sql.DB
}

func NewBibleDatabase(filepath string) (bd *BibleDatabase, err error) {
	bd = &BibleDatabase{}
	//open sqlite file
	bd.database, err = sql.Open("sqlite", getTranlationDBPath(filepath)+"?mode=ro")
	if err != nil {
		return nil, err
	}
	// Force classic journal mode
	_, err = bd.database.Exec(`PRAGMA journal_mode=DELETE;`)
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

func getTranlationDBPath(translation string) string {
	return filepath.Join("bibles", translation+".SQLite3")
}
