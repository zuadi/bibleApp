package models

// type Bibleindex struct {
// 	Bibletranslation string
// 	Bibleverses      string
// }

type BibleIndex [][]string

func (bi *BibleIndex) GetByIndex(index int) (output []string) {
	for _, v := range *bi {
		output = append(output, v[index])
	}
	return
}

func (bi *BibleIndex) GetByValue(value string) (output []string) {
	var index int
	for i, v := range (*bi)[0] {
		if v == value {
			index = i
		}
	}
	for _, v := range *bi {
		output = append(output, v[index])
	}
	return
}
