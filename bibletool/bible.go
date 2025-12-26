package bibletool

import (
	"bibletool/bibletool/models"
	"strconv"
	"strings"
)

// this function reads the entered text/bible verses
// the separation for bibleverses in every line is ", "
// after that : and , will be replaced with ., a - will include from till verses
func (bt *Bibletool) GetBibleVerses(inputtext string, mainSelectIndex int) (bibleVerses models.Bibleverses, notFound models.NotFound) {
	var err error
	bibleVerses, notFound = bt.GetEntries(inputtext)
	//exit if string is empty
	if notFound.Error != nil {
		return bibleVerses, notFound
	}

	// read bibleindex and bibletranslation
	bt.bibleIndex, err = bt.ReadBibleIndexes()
	if err != nil {
		bt.LogError("read csv", err)
		notFound.IsError = true
		notFound.Error = err
		return bibleVerses, notFound
	}

	err = bibleVerses.CheckVerses(mainSelectIndex, bt.bibleIndex)
	notFound.Error = err
	return bibleVerses, notFound
}

func (bt *Bibletool) GetTranslationVerses(mainTranslation models.Bibleverses, translation ...string) (*models.Translations, error) {
	var translations models.Translations
	bt.DebugLog("GetTranslationVerses", "")
	for _, trans := range translation {
		t := &models.Translation{}
		t.SetTranslationName(trans)
		translations = append(translations, t)

		//open sqlite file
		db, err := models.NewBibleDatabase(trans)
		if err != nil {
			bt.Logger.Error("NewBibleDatabase", err)
			return nil, err
		}
		defer db.Close()

		// look for book names and book number
		dbBooks, err := db.GetBooks()
		if err != nil {
			bt.Logger.Error("GetBooks", err)
			return nil, err
		}

		bibleIndex := bt.bibleIndex.GetByValue(trans)
		//iterate over biblevers compilation
		for _, verses := range mainTranslation.BibleVerses {
			var previousVerse int
			paragraph := t.AddParagraph()

			//iterate over verses of compilation
			for _, verse := range verses.Verse {
				//iterate over verses
				for _, rowIndex := range verse.CSVIndex {
					cellData := bibleIndex[rowIndex]

					var chapterI, VerseI int
					if strings.Contains(cellData, ".") {
						split := strings.Split(cellData, ".")
						tmp := split[len(split)-1]
						tmpVerse := strings.TrimSpace(tmp)
						VerseI, err = strconv.Atoi(tmpVerse)
						if err != nil {
							bt.LogError("atoi verse", err)
						}

						chapter := split[len(split)-2]
						split = strings.Split(chapter, " ")
						tmp = split[len(split)-1]
						chapter = strings.TrimSpace(tmp)
						chapterI, err = strconv.Atoi(chapter)
						if err != nil {
							bt.LogError("atoi chapter", err)
						}
					}

					if previousVerse != VerseI {
						var bookFound bool

						for _, dbBook := range dbBooks {
							book := cellData
							words := strings.Fields(cellData)
							if len(words) > 1 {
								book = strings.Join(words[:len(words)-1], " ")
							}

							//get the whole book name and compare it works only with this languages because they have no numbers in the book name
							if t.RightToLeft {
								bookFound = book == dbBook.Name
							} else {
								//get first two letters of bookname for checking if it is the right book
								bookFound = book[:2] == dbBook.GetBookNameTillIndex(2)
							}

							//Check that string contains name of book and check if first two letters are equal
							if strings.Contains(dbBook.GetTrimmedBookName(), book) && bookFound {
								versText, err := db.GetVerse(dbBook.Number, verses.Chapter, verse.Number)
								if err != nil {
									bt.Logger.Error("GetVerse", err)
									continue
								}
								paragraph.AddTitle(book, chapterI, VerseI)
								paragraph.AddVerse(verse.Number, versText)
							}
						}
					}
					previousVerse = VerseI
				}
			}
		}
	}
	return &translations, nil
}
