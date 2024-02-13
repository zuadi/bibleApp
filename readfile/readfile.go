package readfile

import (
	"regexp"
	"strings"
	"unicode"
)

func Notinclude(text string, notinclude []string) (out bool) {
	for _, char := range notinclude {
		if strings.Count(text, char) > 0 {
			out = true
			break
		}
	}
	return out
}
func Read(text string) (outlist string) {
	var tmpwords = make([]string, 0, 2)

	// replace tabs with spaces
	text = strings.Replace(text, "\t", " ", -1)

	// seperate words by whitespace
	lst_words := strings.Split(text, " ")

	for _, item := range lst_words {
		switch {
		case strings.Contains(item, "\r\n"):
			tmpsplit := strings.Split(item, "\r\n")
			tmpwords = append(tmpwords, tmpsplit...)
		case strings.Contains(item, "\r"):
			tmpsplit := strings.Split(item, "\r")
			tmpwords = append(tmpwords, tmpsplit...)
		case strings.Contains(item, "\n"):
			tmpsplit := strings.Split(item, "\n")
			tmpwords = append(tmpwords, tmpsplit...)
		default:
			tmpwords = append(tmpwords, item)
		}
	}

	//iterate over words
	for i, word := range tmpwords {
		//look for verse separator
		for _, separator := range []string{".", ",", ":"} {
			var tmpout string

			//if word contains a separator, first char and last char is digit
			if strings.ContainsAny(word, ".,:") && unicode.IsDigit(rune(word[0])) && unicode.IsDigit(rune(word[len(word)-1])) {

				//check that no .,:,, are in word before bibleverse
				charlist := []string{".", ",", ":,", "!", ";", "?", "'", `"`}
				skip := Notinclude(tmpwords[i-1], charlist)

				//check if book number before book
				if 1 < i && strings.Count(word, separator) == 1 && len(tmpwords[i-1]) > 0 {
					//ckeck of book is written like 1Corinthians to  separate number, check as well if a word contains points if yes skip it
					if unicode.IsDigit(rune(tmpwords[i-1][0])) {
						tmpout = strings.Join([]string{tmpout, string(tmpwords[i-1][0]), " ", tmpwords[i-1][1:], " ", word}, "")
					} else if len(tmpwords[i-2]) > 0 && len(tmpwords[i-2]) <= 4 && !skip {

						if unicode.IsDigit(rune(tmpwords[i-2][0])) {

							tmpout = strings.Join([]string{tmpout, string(tmpwords[i-2]), " ", tmpwords[i-1], " ", word}, "")
							// check if first char is not number or alpabetical
						} else if regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(string(tmpwords[i-1][0])) {
							tmpout = strings.Join([]string{tmpout, tmpwords[i-1], " ", word}, "")
						} else {
							tmpout = strings.Join([]string{tmpout, tmpwords[i-1][1:], " ", word}, "")
						}
					} else if !skip {
						// check if first char is number

						// check if first char is not number or alpabetical
						if regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(string(tmpwords[i-1][0])) {
							tmpout = strings.Join([]string{tmpout, tmpwords[i-1], " ", word}, "")
						} else {
							tmpout = strings.Join([]string{tmpout, tmpwords[i-1][1:], " ", word}, "")
						}
					}
					//check if word before versenumber
				} else if len(tmpwords) > i && strings.Count(word, separator) == 1 && len(tmpwords[i-1]) > 0 {

					if unicode.IsDigit(rune(tmpwords[i-1][0])) {
						tmpout = strings.Join([]string{tmpout, string(tmpwords[i-1][0]), " ", tmpwords[i-1][1:], " ", word}, "")
					} else {
						// check if first char is number

						// check if first char is not number or alpabetical
						if regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(string(tmpwords[i-1][0])) {
							tmpout = strings.Join([]string{tmpout, tmpwords[i-1], " ", word}, "")
							tmpout = strings.Join([]string{tmpout, tmpwords[i-1][1:], " ", word}, "")
						}
					}
				}

				//check if appendix is - and numeric
				if i+2 < len(tmpwords) {

					if strings.Contains(tmpwords[i+1], "-") && unicode.IsDigit(rune(tmpwords[i+2][0])) {
						tmpout = strings.Join([]string{tmpout, tmpwords[i+1], tmpwords[i+2]}, "")
					}
				}
				//check if verse is found like ths Isaiah33.3
			} else if strings.Contains(word, separator) && len(word) > 6 && strings.Count(word, separator) == 1 {
				indexpoint := strings.LastIndex(word, separator)
				//check that word is one index longer than indexpoint
				if len(word) > indexpoint && 0 < indexpoint {

					//if number before and after point
					if unicode.IsDigit(rune(word[indexpoint-1])) && unicode.IsDigit(rune(word[indexpoint+1])) {
						tmpout = strings.Join([]string{tmpout, word}, "")
					}
				}
			}
			if len(tmpout) > 0 {
				outlist = strings.Join([]string{outlist, tmpout, "\n"}, "")
			}
		}

	}
	return outlist
}
