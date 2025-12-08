package main

//This script is a bibletranslation that gives out a .txt file for using the text further and as well as in pdf.
// The maintranslation has to be choosen and further translation can be checked, the checkbox same document includes all text in one txt and pdf document

import (
	"C"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skratchdot/open-golang/open"
)

import (
	"bibletool/basic"
	"bibletool/biblecsvreader"
	"bibletool/biblefunc"
	"bibletool/config"
	"bibletool/modules"
	"bibletool/output"
	"sync"
)

func main() {
	// This function is a bibletranslation help that there can be a main translation choosen and the desired verses entered.
	// it will create a txt and pdf file of each translation or a combined file with all translation.
	// it checks if verse exists in main translation and if not give out the not found verses

	//variables
	var UserChoice modules.UserChoices
	var wg sync.WaitGroup

	//call basic function
	ospath := basic.Settings()

	//read config if file exists
	config.Load(&UserChoice, ospath)

	// read bibleindex and bibletranslation
	Bibleindex := biblecsvreader.ReadCSV(ospath)

	// make new window
	a := app.New()
	w := a.NewWindow("Bibletool")
	r, _ := fyne.LoadResourceFromPath(ospath.IconPath)
	w.SetIcon(r)

	// make title
	label1 := widget.NewLabel("Select Main Translation")

	// make checkboxes
	setofcheck := widget.NewCheckGroup(Bibleindex.Bibletrans, func(s []string) {
		UserChoice.Checkboxes = s
	})

	// show main translation option
	selection := widget.NewSelect(Bibleindex.Bibletrans, func(s string) {
		UserChoice.Maintransation = s
	})

	// same document checkbox
	checkbox_document := widget.NewCheck("All in one document", func(result bool) { UserChoice.SameDocument = result })

	//set last User Choices read from config file
	setofcheck.SetSelected(UserChoice.Checkboxes)
	selection.SetSelected(UserChoice.Maintransation)
	checkbox_document.SetChecked(UserChoice.SameDocument)

	//select all translations
	sel_alltrans := widget.NewCheck("Select all", func(b bool) {
		if b {
			setofcheck.SetSelected(Bibleindex.Bibletrans)
			setofcheck.Refresh()
		} else {
			setofcheck.SetSelected([]string{})
			setofcheck.Refresh()
		}
	})

	//get sermontitle
	sermonname := widget.NewEntry()
	sermonname.PlaceHolder = "Enter Sermon title"
	label3 := widget.NewLabel("Sermontitle:")
	sermonname.SetText(strings.TrimSpace(UserChoice.SermonTitle))

	//get pastor name
	pastorname := widget.NewEntry()
	pastorname.PlaceHolder = "Enter Name of pastor"
	label4 := widget.NewLabel("Name of Pastor:")
	pastorname.SetText(strings.TrimSpace(UserChoice.Pastor))

	label2 := widget.NewLabel("Select Translations")

	verse_entry := widget.NewMultiLineEntry()

	b1 := widget.NewButton("Translate", func() {

		if UserChoice.Maintransation == "" {
			w2 := a.NewWindow("Bibletool")
			w2.Resize(fyne.NewSize(200, 200))
			w2label := widget.NewLabel("No Maintranslation choosen")
			button := widget.NewButton("Ok", func() { w2.Close() })
			w2.SetContent(container.NewCenter(container.NewVBox(w2label, button)))
			w2.Show()

		} else {

			//remove maintranslation of  checkresult
			var templist []string

			for i := range UserChoice.Checkboxes {
				if UserChoice.Checkboxes[i] == UserChoice.Maintransation {
					continue
				}

				templist = append(templist, UserChoice.Checkboxes[i])
			}

			transl_result := templist

			//extract entered bible verses
			entrytext := verse_entry.Text

			// get sermon title
			UserChoice.SermonTitle = sermonname.Text

			//get pastor name
			UserChoice.Pastor = pastorname.Text

			//get bible verses
			verses := biblefunc.Getbibleverses(entrytext)

			VerseCheck := verses.Check_verses(selection.SelectedIndex(), Bibleindex.CSVData)

			var listoutput string

			if len(VerseCheck.Notfoundlist) > 0 {
				listoutput = ""
				for _, item := range VerseCheck.Notfoundlist {
					for _, element := range item {
						listoutput = listoutput + element
					}
				}
			}

			if VerseCheck.Notfound {

				w3 := a.NewWindow("Bibletool")
				w3.Resize(fyne.NewSize(200, 200))
				w3label := widget.NewLabel("This Bibleverses were not found:")

				w3label1 := widget.NewLabel(listoutput)
				button2 := widget.NewButton("Ok", func() { w3.Close() })
				w3.SetContent(container.NewCenter(container.NewVBox(w3label, w3label1, button2)))
				w3.Show()
			} else {

				var currentcount = 1
				var documentname string

				//progressbar
				w4 := a.NewWindow("Bibletool")
				w4.SetIcon(r)
				w4.Resize(fyne.NewSize(250, 250))
				w4label1 := widget.NewLabel("Translating for you:")
				w4document := widget.NewLabel("In progress...")
				docprogress := widget.NewProgressBar()
				w4label2 := widget.NewLabel("Total progress")
				progress := widget.NewProgressBar()
				w4.CenterOnScreen()
				w4.SetContent(container.NewCenter(container.NewVBox(w4label1, w4document, docprogress, w4label2, progress)))
				w4.Show()

				//progressbar for pdf generating
				w5 := a.NewWindow("Bibletool")
				w5.SetIcon(r)
				w5.Resize(fyne.NewSize(200, 100))
				w5label := widget.NewLabel("Make pdf's")
				progresspdf := widget.NewProgressBar()
				w5.CenterOnScreen()
				w5.SetContent(container.NewCenter(container.NewVBox(w5label, progresspdf)))

				//check if biletranslation folder exists otherwise create
				if _, err := os.Stat(ospath.Outputpath); os.IsNotExist(err) {
					err = os.Mkdir(ospath.Outputpath, 0777)
					basic.CheckErr(err, "Error could not create Bibletranslation folder")
				} else {
					err = os.RemoveAll(ospath.Outputpath)
					basic.CheckErr(err, "Error could not remove Bibletranslation folder")
					err = os.Mkdir(ospath.Outputpath, 0777)
					basic.CheckErr(err, "Error could not create Bibletranslation folder")
				}
				//create subfolders of Bibletranslation
				err := os.Mkdir(ospath.Outputpath+ospath.Pathseperator+"html", 0777)
				basic.CheckErr(err, "Error could not create Bibletranslation\\html folder")
				err = os.Mkdir(ospath.Outputpath+ospath.Pathseperator+"txt", 0777)
				basic.CheckErr(err, "Error could not create Bibletranslation\\txt folder")

				text_main := biblefunc.GetVersText(ospath.Currentdirectory+"bibles"+ospath.Pathseperator+UserChoice.Maintransation+".SQLite3", VerseCheck)

				var lst_translationtext = make([]modules.OutputText, 0, 40)
				var lst_translation = make([]string, 0, 40)

				if len(transl_result) > 0 {
					for _, translation := range transl_result {

						outtext := biblefunc.GetTranslationVerses(VerseCheck, translation, Bibleindex.CSVData)

						text_translation := biblefunc.GetVersText(ospath.Currentdirectory+"bibles"+ospath.Pathseperator+translation+".SQLite3", outtext)

						lst_translationtext = append(lst_translationtext, text_translation)
						lst_translation = append(lst_translation, translation)
					}
				}

				if UserChoice.SameDocument { //write in same file
					output.Writesamedoctext(ospath, text_main, lst_translationtext, "Main "+UserChoice.Maintransation, w4document, docprogress)
					output.Writesamehtmlfile(text_main, lst_translationtext, "Main "+UserChoice.Maintransation, UserChoice.SermonTitle, UserChoice.Pastor, ospath, w4document, docprogress)

					//processbar info
					w4document.SetText(documentname)
					progress.Max = 1
					currentcount += 1
					progress.SetValue(float64(currentcount))

				} else { //write seperate files
					// write to file
					output.Writetextfile(ospath, text_main, UserChoice.Maintransation, docprogress)
					output.Writehtmlfile(text_main, "Main", UserChoice.Maintransation, UserChoice.SermonTitle, UserChoice.Pastor, ospath, docprogress, &wg)

					//processbar info
					// translation count for process bar
					progress.Max = float64(len(lst_translation) + 1)
					w4document.SetText(UserChoice.Maintransation)
					progress.SetValue(float64(currentcount))

					for i, text := range lst_translationtext {
						wg.Wait()
						output.Writetextfile(ospath, text, lst_translation[i], docprogress)
						output.Writehtmlfile(text, "Translation", lst_translation[i], UserChoice.SermonTitle, UserChoice.Pastor, ospath, docprogress, &wg)

						//processbar info
						w4document.SetText(lst_translation[i])
						currentcount += 1
						progress.SetValue(float64(currentcount))
					}

					w5.Show()
					w4.Hide()

					wg.Add(len(lst_translation) + 1)
					go func() {
						var progressvalue float64
						for {

							progresspdf.Max = float64(len(lst_translationtext) + 1)
							progressvalue = progresspdf.Max - float64(basic.WaitingQue)
							if progressvalue > 0 {
								w5label.SetText("make pdf's")
								progresspdf.SetValue(float64(progressvalue))
							}
						}

					}()

					wg.Wait()

					if len(lst_translation) > 1 {
						//make one pdf with all translations
						output.CombinedPDF(ospath)
					}

				}

				// open file browser of translation
				err = open.Run(ospath.Outputpath)
				basic.CheckErr(err, "Error open file browser")

				w4.Close()
				w5.Close()
				w.Close()

				// save checkbox settings for next start
				post := UserChoice
				config.Store(post, ospath)
				//remove temporary folder
				defer basic.Deltemp(ospath.Tempdir)

			}
		}
	})

	verse_entry.SetPlaceHolder("Enter bible verses here, like:\nLuke 10.1\nJoh 3.1-14\nPsalm 2.3, 4.5-7\n")

	toprow := container.NewGridWithRows(2, container.NewGridWithColumns(3, label1, label2, checkbox_document), container.NewCenter(sel_alltrans))
	//colunm1 := container.NewVBox(selection, dialogin)
	colunm1 := container.NewVBox(selection)
	colunm2 := container.NewVScroll(setofcheck)
	colunm3 := container.NewVBox(container.NewGridWithRows(6, label3, sermonname, label4, pastorname), b1)

	secondelement := container.NewGridWithColumns(3, colunm1, colunm2, colunm3)
	firstelement := container.NewVBox(toprow, secondelement)
	content := container.NewGridWithColumns(1, firstelement, container.NewMax(verse_entry))

	w.SetContent(content)

	w.ShowAndRun()
}
