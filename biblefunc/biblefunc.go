package biblefunc

import (
	"bibletool/Abbreviation"
	"bibletool/Modules"
	"bibletool/basic"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type VerseList []Modules.Bibleverses

var verselist VerseList

// this function reads the entered text/bible verses
// the separation for bibleverses in every line is ", "
// after that : and , will be replaced with ., a - will include from till verses
func Getbibleverses(inputtext string) *VerseList {

	verselist = VerseList{}

	//exit if string is empty
	if inputtext == "" {
		return &verselist
	}

	//seperate lines of bibleverses with the newpage character \n
	entrylines := strings.Split(inputtext, "\n")

	//seperate book and verses, get items in line to a slice
	var itemofline = make([][]string, 0, 70)
	for _, item := range entrylines {
		// seperate item by ", " and add them to slice
		itemofline = append(itemofline, strings.Split(item, ", "))
	}

	// iterate over lines in entered text
	for _, line := range itemofline {
		var book, chapter = "0", "0"

		//iterate over item in line
		for i, item := range line {
			// change : or , to . in each item of line
			r := strings.NewReplacer(":", ".", ",", ".", ". ", " ")
			item = r.Replace(item)

			// trim return, whitespace and new line
			item = strings.TrimSpace(item)

			var tmp_list Modules.Bibleverses
			//if item is empty string skip
			if item == "" {
				continue
			} else if strings.Count(item, ".") == 0 && strings.Count(item, "-") == 0 {
				// if no . or no - is found in item it is not recognized as a verse
				var tmplst Modules.Bibleverses
				//add to slice for printing out not found bibleverses
				tmplst.Book = append(tmplst.Book, strings.Join([]string{"- '", item, "' is not a valid Bibleverse"}, ""))
				tmplst.Chapter = append(tmplst.Chapter, "0")
				tmplst.Verse = append(tmplst.Verse, "0")
				verselist = append(verselist, tmplst)
				continue
			}

			// check if first item has number to reconize if it is a book with number like 1 Corinthians
			// split at whitespace
			tempstring := strings.Split(item, " ")
			if i == 0 {

				//look if first character is number
				if regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`).MatchString(tempstring[0][:1]) {

					//if 3 element are in list number, book, vers
					if len(tempstring) > 2 {
						book = strings.Join(tempstring[:len(tempstring)-1], " ")
						chapter = strings.Split(tempstring[len(tempstring)-1], ".")[0]
						//if 2 element are in list number, book
					} else if len(tempstring) == 2 {
						book = tempstring[0]
						chapter = strings.Split(tempstring[1], ".")[0]
						//if 1 element are in list number, book
					} else {
						book = tempstring[0]
					}

					// if no book number is found continue as normal book and check how many whitespaces/elemts are in item
					// 6 elements found
				} else if strings.Count(item, " ") > 1 {
					book = strings.Join(tempstring[:strings.Count(item, " ")-1], " ")
					chapter = strings.Split(tempstring[len(tempstring)-1], ".")[0]
					// 5 elements found
				} else {
					book = tempstring[0]
					chapter = strings.Split(tempstring[len(tempstring)-1], ".")[0]
				}

			}
			// if - found in item at all verses from (before -) until after - to slice
			if strings.Contains(item, "-") {

				//get chapter number
				if i > 0 && strings.Count(item, ".") == 1 {
					// if . found get chapter number before .
					chapter = strings.Split(item, ".")[len(strings.Split(item, "."))-2]
				}
				//get first vers before -
				leftside := strings.Split(item, "-")[0]

				// else if no . found get element before -
				temp_leftside := leftside
				// check if . occures more than one time in case the book number is like 1. Mose
				if strings.Count(leftside, ".") > 0 && strings.Count(leftside, ".") < 3 {
					temp_leftside = strings.Split(leftside, ".")[strings.Count(leftside, ".")]
				}

				// trim string from any whitspace or return
				trim_item := strings.TrimSpace(item)
				trim_leftside := strings.TrimSpace(temp_leftside)

				// change right string number to integer end verse
				rightside, err1 := strconv.Atoi(strings.Split(trim_item, "-")[1])
				// change start verse string number to integer
				startvalue, err2 := strconv.Atoi(trim_leftside)
				switch {
				case startvalue > rightside:
					var tmplst Modules.Bibleverses
					//add to slice for printing out not found bibleverses
					tmplst.Book = append(tmplst.Book, strings.Join([]string{"- '", item, "' is not a valid Bibleverse"}, ""))
					tmplst.Chapter = append(tmplst.Chapter, "0")
					tmplst.Verse = append(tmplst.Verse, "0")
					verselist = append(verselist, tmplst)
					// if no error add book chapter and verse to slice
				case err1 == nil && err2 == nil:
					for i := startvalue; i <= rightside; i++ {
						if strings.Contains(leftside, book) {
							tmp_list.Book = append(tmp_list.Book, book)
							tmp_list.Chapter = append(tmp_list.Chapter, chapter)
							tmp_list.Verse = append(tmp_list.Verse, strconv.Itoa(i))
						} else {
							tmp_list.Book = append(tmp_list.Book, book)
							tmp_list.Chapter = append(tmp_list.Chapter, chapter)
							tmp_list.Verse = append(tmp_list.Verse, strconv.Itoa(i))
						}
					}
				default:
					// if error integer print it out
					fmt.Print(err1, err2)
					fmt.Print("error entered value")
				}

				// if item contains . look for chapter and verse
			} else if strings.Contains(item, ".") {
				var vers string

				// check if . occues mor than one time
				switch {
				case i > 0 && strings.Count(item, ".") == 1:
					// add to slice new chapter and verse
					chapter = strings.Split(item, ".")[len(strings.Split(item, "."))-2]
					vers = strings.Split(item, ".")[len(strings.Split(item, "."))-1]
				case strings.Count(item, ".") == 1:
					// add only verse because its first item including book and chapter
					vers = strings.Split(item, ".")[len(strings.Split(item, "."))-1]
				case strings.Count(item, ".") == 2:
					// else choose this for when two . are in item like 1. Mose 1.1
					vers = strings.Split(item, ".")[len(strings.Split(item, "."))-1]
				}

				// add book chapter and verse to temporary slice
				tmp_list.Book = append(tmp_list.Book, book)
				tmp_list.Chapter = append(tmp_list.Chapter, chapter)
				tmp_list.Verse = append(tmp_list.Verse, vers)
			}

			// add temp slice including from - to verses to final verse slice
			if len(tmp_list.Book) > 0 {
				verselist = append(verselist, tmp_list)
			}
		}
	}
	//return slice of entered verses of main text
	return &verselist
}

func (in_list *VerseList) Check_verses(main_selectindex int, data [][]string) Modules.Checkverse {
	// This function checks if the entered verses are in the main translation bible, gives out a list of not found verses and, gives out a slice with
	// all found bibleverses of maintext and index of the config file to find translation text in csv
	var Checkstruct Modules.Checkverse

	// insert first empty line to not found slice
	Checkstruct.Notfoundlist = append(Checkstruct.VersList, []string{"\n"})

	// if list is empty give this message out
	if len(*in_list) == 0 {
		Checkstruct.Notfound = true
		Checkstruct.Notfoundlist = append(Checkstruct.VersList, []string{"- No Biblevers entered\n"})
		return Checkstruct
	}

	//add 1 to skip first colum that is UID of config file
	tmp_maintranslationindex := main_selectindex + 1

	//look for bibleverses
	for _, verscompilation := range *in_list {
		// variables
		var temp_lst1, temp_lst2, temp_lst3 []string
		var temp_lst4 []int
		var bookcompare, tmp_book string
		var identicalbooks bool
		not_found := true

		for i := range verscompilation.Book {
			// if not valid bibleverse continue
			if strings.Contains(verscompilation.Book[i], "not a valid") {
				Checkstruct.Notfoundlist = append(Checkstruct.Notfoundlist, []string{verscompilation.Book[i], "\n"})
				continue
			}

			not_found = true
			// iterate over rows of csv config file
			for ii, row := range data {
				// iterate over colum of row in csv config file
				for iii, item := range row {
					var tmp_chapter, tmp_vers string
					var tmp_split []string
					// if biblename is found in column get full book, chapter and verse name
					if iii == tmp_maintranslationindex {
						if strings.Contains(item, ".") {

							tmp_split = strings.Split(item, ".")
							tmp := tmp_split[len(tmp_split)-1]
							tmp_vers = strings.TrimSpace(tmp)

							tmp_chapter = tmp_split[len(tmp_split)-2]

							tmp_split = strings.Split(tmp_chapter, " ")
							tmp = tmp_split[len(tmp_split)-1]
							tmp_chapter = strings.TrimSpace(tmp)
						}
						// replace short names with real names
						verscompilation.Book[i] = Abbreviation.Replace(verscompilation.Book[i])
						if strings.Contains(strings.ToLower(item), strings.ToLower(verscompilation.Book[i])) && tmp_chapter == strings.TrimSpace(verscompilation.Chapter[i]) && tmp_vers == strings.TrimSpace(verscompilation.Verse[i]) && strings.ToLower(string(item[0])) == strings.ToLower(string(verscompilation.Book[i][0])) {
							not_found = false

							//book
							//trim item
							tmp_item := strings.TrimSpace(item)
							if strings.Count(tmp_item, " ") > 0 {
								tmp_item2 := strings.Split(tmp_item, " ")

								for i := 0; i+1 <= strings.Count(tmp_item, " "); i++ {
									if i == 0 {
										tmp_book = tmp_item2[i]
									} else {
										tmp_book = strings.Join([]string{tmp_book, " ", tmp_item2[i]}, "")
									}
								}
							}

							temp_lst1 = append(temp_lst1, tmp_book)
							//chapter
							temp_lst2 = append(temp_lst2, tmp_chapter)
							//vers
							temp_lst3 = append(temp_lst3, tmp_vers)
							//bibleversindex
							temp_lst4 = append(temp_lst4, ii)

							// check if more than one book fits
							if bookcompare != tmp_book && bookcompare != "" {
								not_found = true
								identicalbooks = true
							}
							bookcompare = tmp_book
						}
					}
				}
			}
			if not_found {
				tmpnotfound := strings.Join([]string{"- ", strings.TrimSpace(verscompilation.Book[i]), " ", strings.TrimSpace(verscompilation.Chapter[i]), ".", strings.TrimSpace(verscompilation.Verse[i]), "\n"}, "")
				if identicalbooks {
					tmpnotfound = strings.Replace(tmpnotfound, "\n", strings.Join([]string{" more than one book found with '", verscompilation.Book[i], "', need more specific bookname", "\n"}, ""), 1)
				}
				Checkstruct.Notfoundlist = append(Checkstruct.Notfoundlist, []string{tmpnotfound})
			}
		}
		Checkstruct.BookList = append(Checkstruct.BookList, temp_lst1)
		Checkstruct.ChapterList = append(Checkstruct.ChapterList, temp_lst2)
		Checkstruct.VersList = append(Checkstruct.VersList, temp_lst3)
		Checkstruct.Versindex = append(Checkstruct.Versindex, temp_lst4)

	}

	if len(Checkstruct.Notfoundlist) == 1 {
		Checkstruct.Notfound = false
		Checkstruct.Notfoundlist = nil
	} else if len(Checkstruct.Notfoundlist) > 1 {
		Checkstruct.Notfound = true
	}

	return Checkstruct
}

func GetVersText(filepath string, input Modules.Checkverse) (output_text Modules.OutputText) {
	var lst_books []struct {
		Name   string
		Number int
	}

	var temp_out struct {
		Titel string
		Verse [][]string
	}

	//check if file exists

	//open sqlite file
	db, err := sql.Open("sqlite3", filepath)
	basic.CheckErr(err, "Error open, SQLite3 file for checking bible verses:")

	defer db.Close()

	// look for book names and book number
	rows, err := db.Query("SELECT book_number, long_name FROM books")
	basic.CheckErr(err, strings.Join([]string{"Error looking up data of ", filepath, " file"}, ""))

	var book_number int
	var long_name string

	for rows.Next() {
		err = rows.Scan(&book_number, &long_name)
		lst_books = append(lst_books, struct {
			Name   string
			Number int
		}{Name: long_name, Number: book_number})
		basic.CheckErr(err, strings.Join([]string{"Error reading data in rows of ", filepath, " file"}, ""))

	}

	var titel string
	var lst_vers []string

	//look for current biblevers book
	for i, elem := range input.VersList {

		lst_vers = []string{}
		temp_out.Titel, titel = "", ""
		temp_out.Verse = [][]string{}
		var beginningnumber, savedVersnumber = "0", "0"
		var Biblebookname, Biblebooknameiterate string

		for ii := range elem {

			for _, value := range lst_books {

				//get the whole book name and compare it works only with this languages because they have no numbers in the book name
				if strings.Contains(filepath, "Arabic") || strings.Contains(filepath, "Hebrew") || strings.Contains(filepath, "Persian") || strings.Contains(filepath, "Aramaic") {
					Biblebookname = input.BookList[i][ii]
					Biblebooknameiterate = value.Name
				} else {
					//get first two letters of bookname for checking if it is the right book
					Biblebookname = input.BookList[i][ii][0:2]
					Biblebooknameiterate = value.Name[0:2]
				}

				//Check taht string contains name of book and check if first two letters are equal
				if strings.Contains(strings.TrimSpace(value.Name), strings.TrimSpace(input.BookList[i][ii])) && Biblebooknameiterate == Biblebookname {

					rows, err = db.Query("SELECT text FROM verses WHERE book_number = ? AND chapter = ? AND verse = ?", value.Number, input.ChapterList[i][ii], input.VersList[i][ii])
					basic.CheckErr(err, strings.Join([]string{"Error at get verse in file ", filepath}, ""))

					var text string

					VersnumebrforOutput := input.VersList[i][ii]
					//write vers titel
					if savedVersnumber == "0" {
						titel = strings.Join([]string{value.Name, " ", input.ChapterList[i][ii], ".", input.VersList[i][ii]}, "")
						beginningnumber = input.VersList[i][ii]
					} else {
						titel = strings.Join([]string{value.Name, " ", input.ChapterList[i][ii], ".", beginningnumber, "-", input.VersList[i][ii]}, "")
					}

					savedVersnumber = input.VersList[i][ii]

					for rows.Next() {
						err = rows.Scan(&text)
						basic.CheckErr(err, strings.Join([]string{"Error at scan rows in file ", filepath}, ""))

						// remove special characters
						temp_text := text

						r := strings.NewReplacer("<pb/>", "", "<i>", "", "</i>", "", "<t>", "", "</t>", "", "<e>", "", "</e>", "", "<J>", "", "</J>", "", "[T›]", "", "<br/>", "")
						temp_text = r.Replace(temp_text)

						// remove footnotes and notes inside <> </>
						footnoteschars := []string{"f", "n", "S"}
						for {
							breakloop := true

							// look if footnote char is in text
							for _, value := range footnoteschars {
								if strings.Contains(temp_text, strings.Join([]string{"<", value, ">"}, "")) {
									breakloop = false
									temp_text = strings.Join([]string{temp_text[:strings.Index(temp_text, strings.Join([]string{"<", value, ">"}, ""))], temp_text[strings.Index(temp_text, strings.Join([]string{"</", value, ">"}, ""))+len(strings.Join([]string{"</", value, ">"}, "")):]}, "")
								}
							}

							// quit the loop if no footnote char is found
							if breakloop {
								lst_vers = append(lst_vers, strings.Join([]string{VersnumebrforOutput, " ", temp_text}, ""))
								break
							}
						}
					}
				}
			}
		}
		temp_out.Titel = titel
		temp_out.Verse = append(temp_out.Verse, lst_vers)

		output_text = append(output_text, temp_out)
	}

	return output_text
}

func GetTranslationVerses(input Modules.Checkverse, tranls_selection string, data [][]string) (outtext Modules.Checkverse) {
	var translation_column = -1

	//iterate over biblevers compilation
	for _, verscompilation := range input.Versindex {

		var tmp_book string
		var temp_lst1, temp_lst2, temp_lst3 []string
		var temp_lst4 []int

		//iterate over verses of compilation
		for _, versind := range verscompilation {

			//iterate over csv sheet
			for ii, row := range data {

				// skip if not translation column or not vers row
				if ii != versind && translation_column != -1 {
					continue
				}

				for iii, item := range row {
					var tmp_chapter, tmp_vers string
					var tmp_split []string

					if ii == 0 && data[ii][iii] == tranls_selection {
						//save column of translation
						translation_column = iii
						continue
					} else if translation_column != iii {
						continue
					}

					if iii == translation_column && versind == ii {

						if strings.Contains(item, ".") {

							tmp_split = strings.Split(item, ".")
							tmp := tmp_split[len(tmp_split)-1]
							tmp_vers = strings.TrimSpace(tmp)

							tmp_chapter = tmp_split[len(tmp_split)-2]
							tmp_split = strings.Split(tmp_chapter, " ")
							tmp = tmp_split[len(tmp_split)-1]
							tmp_chapter = strings.TrimSpace(tmp)
						}
						//book
						//trim item
						tmp_item := strings.TrimSpace(item)
						if strings.Count(tmp_item, " ") > 0 {
							tmp_item2 := strings.Split(tmp_item, " ")

							for i := 0; i+1 <= strings.Count(tmp_item, " "); i++ {
								if i == 0 {
									tmp_book = tmp_item2[i]
								} else {
									tmp_book = strings.Join([]string{tmp_book, " ", tmp_item2[i]}, "")
								}
							}
						}

						temp_lst1 = append(temp_lst1, tmp_book)
						//chapter
						temp_lst2 = append(temp_lst2, tmp_chapter)
						//vers
						temp_lst3 = append(temp_lst3, tmp_vers)
						//bibleversindex
						temp_lst4 = append(temp_lst4, ii)

						break
					}
				}
			}
		}
		outtext.BookList = append(outtext.BookList, temp_lst1)
		outtext.ChapterList = append(outtext.ChapterList, temp_lst2)
		outtext.VersList = append(outtext.VersList, temp_lst3)
		outtext.Versindex = append(outtext.Versindex, temp_lst4)
	}

	return outtext
}
