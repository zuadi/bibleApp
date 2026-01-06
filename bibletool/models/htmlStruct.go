package models

import "html/template"

type HtmlStruct struct {
	Name                string
	SermonTitle         string
	PastorName          string
	MainTranslation     *Translation
	Translations        *Translations
	Date                string
	IconBase64          template.URL
	ProgressFnc         func(title string, progress float64)
	RightToLeftDocument bool
	SameDocument        bool
}

// is used in tmpl
func (h *HtmlStruct) Progress(title string, progress float64) string {
	if h.ProgressFnc == nil {
		return ""
	}
	h.ProgressFnc(title, progress)
	return ""
}
