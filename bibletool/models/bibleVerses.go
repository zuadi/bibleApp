package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gitea.tecamino.com/paadi/tecamino-logger/logging"
)

type Bibleverses struct {
	Notfound     bool
	NotFoundList []string
	BibleVerses  []*Bibleverse
	Logger       *logging.Logger
}

func (bv *Bibleverses) AddNewBibleverse() *Bibleverse {
	v := &Bibleverse{}
	bv.BibleVerses = append(bv.BibleVerses, v)
	return v
}

func (bv *Bibleverses) AddBibleVerse(nbv *Bibleverse) {
	bv.BibleVerses = append(bv.BibleVerses, nbv)
}

func (bv *Bibleverses) CheckVerses(mainTranslationIndex int, data BibleIndex) (err error) {
	// This function checks if the entered verses are in the main translation bible, gives out a list of not found verses and, gives out a slice with
	// all found bibleverses of maintext and index of the config file to find translation text in csv
	bv.Logger.Debug("CheckVerses", "")
	// if list is empty give this message out
	if len(bv.BibleVerses) == 0 {
		return errors.New("- No Biblevers entered\n")
	}

	//add 1 to skip first colum that is UID of config file
	bibleIndex := data.GetByIndex(mainTranslationIndex + 1)

	//look for bibleverses
	for _, bibleVerses := range bv.BibleVerses {
		// variables
		var compareBook string
		var notFoundCount int
		for _, verse := range bibleVerses.Verse {
			// if not valid bibleverse continue
			if strings.Contains(bibleVerses.Book, "not a valid") {
				bv.AddNotFound(bibleVerses.Book + "\n")
				continue
			}

			var found bool

			// iterate over rows of csv config file
			for rowIndex, cellData := range bibleIndex {
				// iterate over colum of row in csv config file

				var chapterI, verseI int
				var splitedString []string
				// if biblename is found in column get full book, chapter and verse name
				if strings.Contains(cellData, ".") {

					splitedString = strings.Split(cellData, ".")
					tmp := splitedString[len(splitedString)-1]
					verseI, err = strconv.Atoi(strings.TrimSpace(tmp))
					if err != nil {
						bv.Logger.Error("verse atoi error", err)
					}
					tmpChapter := splitedString[len(splitedString)-2]

					splitedString = strings.Split(tmpChapter, " ")
					tmp = splitedString[len(splitedString)-1]
					tmpChapter = strings.TrimSpace(tmp)
					chapterI, err = strconv.Atoi(tmpChapter)
					if err != nil {
						bv.Logger.Error("chapter atoi error", err)
					}
				}
				// replace short names with real names
				bibleVerses.ReplaceBookAbbreviation()
				if strings.Contains(strings.ToLower(cellData), strings.ToLower(bibleVerses.Book)) &&
					chapterI == bibleVerses.Chapter && verseI == verse.Number &&
					strings.EqualFold(string(cellData[0]), strings.ToLower(string(bibleVerses.Book[0]))) {
					found = true

					//book
					//trim item
					book := cellData
					words := strings.Fields(cellData)
					if len(words) > 1 {
						// Join everything except the last element back together with a single space
						book = strings.Join(words[:len(words)-1], " ")
					}
					bibleVerses.BookName = book

					bibleVerses.GetVerse(verseI).AddCSVIndex(rowIndex)
					// check if more than one book fits
					if compareBook != bibleVerses.BookName && compareBook != "" {
						bv.AddNotFound(fmt.Sprintf("- %s %d.%d more than one book found with '%s', need more specific bookname\n", strings.TrimSpace(bibleVerses.Book), bibleVerses.Chapter, verse.Number, bibleVerses.Book))
					}
					compareBook = bibleVerses.BookName
				}
			}

			if notFoundCount > 1 {
				break
			} else if notFoundCount == 1 {
				bv.AddSecondNotFound(fmt.Sprintf("- %s %d.%d, %d ...\n", strings.TrimSpace(bibleVerses.Book), bibleVerses.Chapter, verse.Number-1, verse.Number))
				notFoundCount += 1
			} else if !found {
				bv.AddNotFound(fmt.Sprintf("- %s %d.%d\n", strings.TrimSpace(bibleVerses.Book), bibleVerses.Chapter, verse.Number))
				notFoundCount += 1
			}
		}
	}
	return bv.GetAllNotFound()
}

func (bv *Bibleverses) AddNotFound(verse string) {
	bv.Logger.Debug("AddNotFound", "")
	bv.Notfound = true
	bv.NotFoundList = append(bv.NotFoundList, verse)
}

func (bv *Bibleverses) AddSecondNotFound(verse string) {
	bv.NotFoundList[len(bv.NotFoundList)-1] = verse
}

func (bv *Bibleverses) GetAllNotFound() error {
	bv.Logger.Debug("GetAllNotFound", "")
	var stringBuilder strings.Builder
	for i := range bv.NotFoundList {
		stringBuilder.WriteString(bv.NotFoundList[i])
	}
	notFound := stringBuilder.String()
	if notFound != "" {
		return errors.New(notFound)
	}
	return nil
}

func (bv *Bibleverses) GetMainVerseText(filepath string) (translation *Translation, err error) {
	translation = &Translation{IsMain: true}
	translation.SetTranslationName(filepath)

	bv.Logger.Debug("GetMainVerseText", "open sqlite database")
	//open sqlite file
	db, err := NewBibleDatabase(filepath)
	if err != nil {
		bv.Logger.Error("NewBibleDatabase", err)
		return translation, err
	}
	defer db.Close()

	bv.Logger.Debug("GetMainVerseText", "get all books")
	// look for book names and book number
	dbBooks, err := db.GetBooks()
	if err != nil {
		bv.Logger.Error("GetBooks", err)
		return translation, err
	}

	bv.Logger.Debug("GetMainVerseText", "get all bible verses")
	//look for current biblevers book
	for _, verses := range bv.BibleVerses {
		var bookFound bool

		paragraph := translation.AddParagraph()
		for _, verse := range verses.Verse {
			for _, dbBook := range dbBooks {
				//get the whole book name and compare it works only with this languages because they have no numbers in the book name
				if translation.RightToLeft {
					bookFound = verses.Book == dbBook.Name
				} else {
					//get first two letters of bookname for checking if it is the right book
					bookFound = verses.GetBookNameTillIndex(2) == dbBook.GetBookNameTillIndex(2)
				}

				//Check that string contains name of book and check if first two letters are equal
				if strings.Contains(dbBook.GetTrimmedBookName(), verses.GetTrimmedBookName()) && bookFound {
					//write verse titel
					paragraph.AddTitle(dbBook.Name, verses.Chapter, verse.Number)

					versText, err := db.GetVerse(dbBook.Number, verses.Chapter, verse.Number)
					if err != nil {
						bv.Logger.Error("GetVerse", err)
						return translation, err
					}
					paragraph.AddVerse(verse.Number, versText)
				}
			}
		}
	}
	return translation, nil
}
