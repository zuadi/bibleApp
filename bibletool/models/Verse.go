package models

type Verse struct {
	RightToLeft bool
	CSVIndex    []int
	Number      int
	Text        string
}

func (v *Verse) AddCSVIndex(index int) {
	v.CSVIndex = append(v.CSVIndex, index)
}
