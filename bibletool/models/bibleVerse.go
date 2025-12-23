package models

import (
	"strings"
)

type Bibleverse struct {
	Book     string
	BookName string
	Chapter  int
	Verse    []*Verse
}

func (bv *Bibleverse) GetBookNameTillIndex(index int) string {
	return bv.BookName[:index]
}

func (bv *Bibleverse) GetTrimmedBookName() string {
	return strings.TrimSpace(bv.BookName)
}
func (bv *Bibleverse) AddVerse(book string, chapter, versNumber int) {
	bv.Book = book
	bv.Chapter = chapter
	bv.Verse = append(bv.Verse, &Verse{Number: versNumber})
}

func (bv *Bibleverse) ReplaceBookAbbreviation() {
	bv.Book = ReplaceAbbreviation(bv.Book)
}

func (bv *Bibleverse) GetVerse(number int) *Verse {
	for _, v := range bv.Verse {
		if v.Number == number {
			return v
		}
	}
	return nil
}
