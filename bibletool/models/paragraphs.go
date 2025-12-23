package models

import "fmt"

type Paragraphs []*Paragraph

type Paragraph struct {
	Title       string
	Verse       []Verse
	startVerse  int
	RightToLeft bool
}

func (p *Paragraph) AddTitle(bookName string, chapterNumber, verseNumber int) {
	if p.Title == "" {
		p.Title = fmt.Sprintf("%s %d.%d", bookName, chapterNumber, verseNumber)
		p.startVerse = verseNumber
	} else {
		p.Title = fmt.Sprintf("%s %d.%d-%d", bookName, chapterNumber, p.startVerse, verseNumber)
	}
}

func (p *Paragraph) AddVerse(number int, text string) {
	v := Verse{
		Number: number,
		Text:   text,
	}
	p.Verse = append(p.Verse, v)
}

func (p *Paragraph) GetVerseAmount() int {
	return len(p.Verse)
}
