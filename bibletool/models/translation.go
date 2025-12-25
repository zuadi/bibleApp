package models

import (
	"fmt"
	"regexp"
)

type Translations []*Translation

func (t *Translations) GetParagraphAmount() int {
	var amount int
	for _, tr := range *t {
		amount += tr.GetPargraphAmount()
	}
	return amount
}

func (t *Translations) GetVerseAmount() int {
	var amount int
	for _, tr := range *t {
		amount += tr.GetVerseAmount()
	}
	return amount
}

type Translation struct {
	Name        string
	IsMain      bool
	RightToLeft bool
	Paragraphs  Paragraphs
}

func (t *Translation) SetTranslationName(translation string) {
	t.Name = translation
	// The | character acts as "OR"
	re := regexp.MustCompile(`Arabic|Hebrew|Persian|Aramaic`)
	t.RightToLeft = re.MatchString(translation)
}

func (t *Translation) AddParagraph() *Paragraph {
	p := &Paragraph{}
	t.Paragraphs = append(t.Paragraphs, p)
	return p
}

func (t *Translation) GetTranslationName() string {
	if t.IsMain {
		return fmt.Sprintf("Main %s", t.Name)
	}
	return fmt.Sprintf("Translation %s", t.Name)
}

func (t *Translation) GetPargraphAmount() int {
	return len(t.Paragraphs)
}

func (t *Translation) GetVerseAmount() int {
	var amount int
	for _, p := range t.Paragraphs {
		amount += p.GetVerseAmount()
	}
	return amount
}

func (t *Translation) GetParagraphByIndex(index int) *Paragraph {
	return t.Paragraphs[index]
}
