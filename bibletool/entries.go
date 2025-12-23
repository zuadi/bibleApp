package bibletool

import (
	"bibletool/bibletool/models"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Entries []Entry

type Entry []string

func (bt *Bibletool) GetEntries(input string) (bibleVerses models.Bibleverses, err error) {
	var entries Entries
	bibleVerses.Logger = bt.Logger
	//exit if string is empty
	if input == "" {
		return bibleVerses, errors.New("no bibleverses entered")
	}

	//seperate lines of bibleverses with the newpage character \n
	entryLines := strings.SplitSeq(input, "\n")

	//seperate book and verses, get items in line to a slice
	for item := range entryLines {
		// seperate item by ", " and add them to slice
		entries = append(entries, strings.Split(item, ", "))

	}

	// iterate over lines in entered text
	for _, entry := range entries {
		var book string
		var chapter int

		//iterate over item in line
		for i, e := range entry {
			// change : or , to . in each entry of line
			e = strings.NewReplacer(":", ".", ",", ".", ". ", " ").Replace(e)

			// trim return, whitespace and new line
			e = strings.TrimSpace(e)

			var newBibleVerse models.Bibleverse
			//if item is empty string skip
			if e == "" {
				continue
			} else if strings.Count(e, ".") == 0 && strings.Count(e, "-") == 0 {
				// if no . or no - is found in item it is not recognized as a verse
				bibleVerses.AddNewBibleverse().AddVerse(fmt.Sprintf("- '%s' is not a valid Bibleverse", e), 0, 0)
				continue
			}

			// check if first item has number to reconize if it is a book with number like 1 Corinthians
			// split at whitespace
			tempstring := strings.Split(e, " ")
			if i == 0 {

				//look if first character is number
				if regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`).MatchString(tempstring[0][:1]) {

					//if 3 element are in list number, book, vers
					if len(tempstring) > 2 {
						book = strings.Join(tempstring[:len(tempstring)-1], " ")
						chapter, err = strconv.Atoi(strings.Split(tempstring[len(tempstring)-1], ".")[0])
						if err != nil {
							bt.LogError("chapter to number", err)
						}
						//if 2 element are in list number, book
					} else if len(tempstring) == 2 {
						book = tempstring[0]
						chapter, err = strconv.Atoi(strings.Split(tempstring[1], ".")[0])
						if err != nil {
							bt.LogError("chapter to number", err)
						}
						//if 1 element are in list number, book
					} else {
						book = tempstring[0]
					}

					// if no book number is found continue as normal book and check how many whitespaces/elemts are in item
					// 6 elements found
				} else if strings.Count(e, " ") > 1 {
					book = strings.Join(tempstring[:strings.Count(e, " ")-1], " ")
					chapter, err = strconv.Atoi(strings.Split(tempstring[len(tempstring)-1], ".")[0])
					if err != nil {
						bt.LogError("chapter to number", err)
					}
					// 5 elements found
				} else {
					book = tempstring[0]
					chapter, err = strconv.Atoi(strings.Split(tempstring[len(tempstring)-1], ".")[0])
					if err != nil {
						bt.LogError("chapter to number", err)
					}
				}

			}
			// if - found in item at all verses from (before -) until after - to slice
			if strings.Contains(e, "-") {

				//get chapter number
				if i > 0 && strings.Count(e, ".") == 1 {
					// if . found get chapter number before .
					chapter, err = strconv.Atoi(strings.Split(e, ".")[len(strings.Split(e, "."))-2])
					if err != nil {
						bt.LogError("chapter to number", err)
					}
				}
				//get first vers before -
				leftside := strings.Split(e, "-")[0]

				// else if no . found get element before -
				temp_leftside := leftside
				// check if . occures more than one time in case the book number is like 1. Mose
				if strings.Count(leftside, ".") > 0 && strings.Count(leftside, ".") < 3 {
					temp_leftside = strings.Split(leftside, ".")[strings.Count(leftside, ".")]
				}

				// trim string from any whitspace or return
				trim_item := strings.TrimSpace(e)
				trim_leftside := strings.TrimSpace(temp_leftside)

				// change right string number to integer end verse
				rightside, err1 := strconv.Atoi(strings.Split(trim_item, "-")[1])
				if err1 != nil {
					bt.LogError("getbiblevers atoi err1", err1)
				}

				// change start verse string number to integer
				startvalue, err2 := strconv.Atoi(trim_leftside)
				if err2 != nil {
					bt.LogError("getbiblevers atoi err2", err2)
				}
				switch {
				case startvalue > rightside:
					//add to slice for printing out not found bibleverses
					bv := bibleVerses.AddNewBibleverse()
					bv.AddVerse(fmt.Sprintf("- '%s' is not a valid Bibleverse", e), 0, 0)

					// if no error add book chapter and verse to slice
				case err1 == nil && err2 == nil:
					for i := startvalue; i <= rightside; i++ {
						newBibleVerse.AddVerse(book, chapter, i)

					}
				}

				// if item contains . look for chapter and verse
			} else if strings.Contains(e, ".") {
				var vers int

				// check if . occues more than one time
				switch {
				case i > 0 && strings.Count(e, ".") == 1:
					// add to slice new chapter and verse
					chapter, err = strconv.Atoi(strings.Split(e, ".")[len(strings.Split(e, "."))-2])
					if err != nil {
						bt.LogError("chapter to number", err)
					}
					vers, err = strconv.Atoi(strings.Split(e, ".")[len(strings.Split(e, "."))-1])
					if err != nil {
						bt.LogError("verse to number", err)
					}
				case strings.Count(e, ".") == 1:
					// add only verse because its first item including book and chapter
					vers, err = strconv.Atoi(strings.Split(e, ".")[len(strings.Split(e, "."))-1])
					if err != nil {
						bt.LogError("verse to number", err)
					}
				case strings.Count(e, ".") == 2:
					// else choose this for when two . are in item like 1. Mose 1.1
					vers, err = strconv.Atoi(strings.Split(e, ".")[len(strings.Split(e, "."))-1])
					if err != nil {
						bt.LogError("verse to number", err)
					}
				}

				// add book chapter and verse to temporary slice
				newBibleVerse.AddVerse(book, chapter, vers)
			}

			// add temp slice including from - to verses to final verse slice
			bibleVerses.AddBibleVerse(&newBibleVerse)
		}
	}
	return
}
