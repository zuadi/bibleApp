package Modules

type Text string

type Bibleindex struct {
	Bibletranslation string
	Bibleverses      string
}
type Bibleverses struct {
	Book    []string
	Chapter []string
	Verse   []string
}
type Checkverse struct {
	Notfound     bool
	Notfoundlist [][]string
	BookList     [][]string
	ChapterList  [][]string
	VersList     [][]string
	Versindex    [][]int
}

type OutputText []struct {
	Titel string
	Verse [][]string
}

type UserChoices struct {
	Maintransation string
	Checkboxes     []string
	SameDocument   bool
	SermonTitle    string
	Pastor         string
	Folderpath     string
}
