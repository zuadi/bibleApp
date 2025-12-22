package main

//This script is a bibletranslation that gives out a .txt file for using the text further and as well as in pdf.
// The maintranslation has to be choosen and further translation can be checked, the checkbox same document includes all text in one txt and pdf document

import (
	"C"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)
import (
	"bibletool/bibletool"
	"bibletool/bibletool/consts"
	"fmt"
	"path/filepath"
)

func main() {
	// This function is a bibletranslation help that there can be a main translation choosen and the desired verses entered.
	// it will create a txt and pdf file of each translation or a combined file with all translation.
	// it checks if verse exists in main translation and if not give out the not found verses

	bt, err := bibletool.NewBibletool()
	if err != nil {
		panic(err)
	}

	//variables
	// var wg sync.WaitGroup

	// all bible translation
	allTranslations, err := bt.GetAllTranslations()
	if err != nil {
		panic(err)
	}

	// make new fyne app
	app := app.New()
	mainWindow := app.NewWindow(consts.AppName)
	mainWindow.SetMaster()

	r, err := fyne.LoadResourceFromPath(consts.IconPath)
	if err != nil {
		bt.LogError("load icon resource to fyne", err)
		panic(err)
	}
	mainWindow.SetIcon(r)

	// make title
	label1 := widget.NewLabel("Select Main Translation")

	// make checkboxes
	setofcheck := widget.NewCheckGroup(allTranslations, func(s []string) {
		bt.SetTranslations(s)
	})

	// show main translation option
	selection := widget.NewSelect(allTranslations, func(s string) {
		bt.SetMaintranslation(s)
	})

	// same document checkbox
	checkbox_document := widget.NewCheck("All in one document", func(result bool) { bt.SetSameDocument(result) })

	//set last User Choices read from config file
	setofcheck.SetSelected(bt.GetSelectedTranslations())
	selection.SetSelected(bt.GetMaintranslation())
	checkbox_document.SetChecked(bt.GetSameDocument())

	//select all translations
	sel_alltrans := widget.NewCheck("Select all", func(b bool) {
		t := []string{}
		if b {
			t = allTranslations
		}
		setofcheck.SetSelected(t)
		setofcheck.Refresh()
	})

	//get sermontitle
	sermonname := widget.NewEntry()
	sermonname.PlaceHolder = "Enter Sermon title"
	label3 := widget.NewLabel("Sermontitle:")
	sermonname.SetText(bt.GetSermonTitle())

	//get pastor name
	pastorname := widget.NewEntry()
	pastorname.PlaceHolder = "Enter Name of pastor"
	label4 := widget.NewLabel("Name of Pastor:")
	pastorname.SetText(bt.GetPastor())

	label2 := widget.NewLabel("Select Translations")

	verse_entry := widget.NewMultiLineEntry()

	var openWindow bool
	w2 := app.NewWindow(consts.AppName)

	b1 := widget.NewButton("Translate", func() {
		// warning error window still open
		if openWindow {
			w2.RequestFocus()
			return
		}

		go func() {
			// check whether main translation is not set
			mainTranslation := bt.GetMaintranslation()
			if mainTranslation == "" {
				openWindow = true
				w2.Resize(fyne.NewSize(200, 200))
				w2.SetContent(container.NewCenter(container.NewVBox(
					widget.NewLabel("No Maintranslation choosen"),
					widget.NewButton("Ok", func() { w2.Close() }),
				),
				),
				)
				w2.Show()
				w2.SetOnClosed(func() { openWindow = false })
				return
			}

			/* start translating */

			// save sermon title
			bt.SetSermonTitle(sermonname.Text)

			//save pastor name
			bt.SetPastor(pastorname.Text)

			//get bible verses
			bibleVerses, err := bt.GetBibleVerses(verse_entry.Text, selection.SelectedIndex())
			if err != nil {
				var errText string
				var labelText = "This Bibleverses were not found:"

				if err.Error() != "no bibleverses entered" {
					errText = err.Error()
				} else {
					labelText = "No bibleverses entered"
				}

				openWindow = true
				w2 = app.NewWindow(consts.AppName)
				w2.Resize(fyne.NewSize(200, 200))

				w2.SetContent(container.NewCenter(container.NewVBox(
					widget.NewLabel(labelText),
					widget.NewLabel(errText),
					widget.NewButton("Ok", func() { w2.Close() }),
				),
				),
				)
				w2.Show()
				w2.SetOnClosed(func() { openWindow = false })
				return
			}

			bt.TotalProgress = func(progress float64) {
				fmt.Println("total", progress)
			}

			bt.PdfProgress = func(progress float64) {
				fmt.Println("pdf", progress)
			}

			// var currentcount = 0.0
			// var documentname string

			bt.SetVerses(verse_entry.Text)
			//progressbar
			w2 = app.NewWindow(consts.AppName)
			w2.SetIcon(r)
			w2.Resize(fyne.NewSize(250, 250))

			w4document := widget.NewLabel("In progress...")
			docprogress := widget.NewProgressBar()
			progress := widget.NewProgressBar()

			w2.CenterOnScreen()
			w2.SetContent(container.NewCenter(container.NewVBox(
				widget.NewLabel("Translating for you:"),
				w4document,
				docprogress,
				widget.NewLabel("Total progress"),
				progress,
			),
			),
			)
			w2.Show()

			//progressbar for pdf generating
			w5 := app.NewWindow(consts.AppName)
			w5.SetIcon(r)
			w5.Resize(fyne.NewSize(200, 100))
			w5label := widget.NewLabel("Make pdf's")
			progresspdf := widget.NewProgressBar()
			w5.CenterOnScreen()
			w5.SetContent(container.NewCenter(container.NewVBox(
				w5label,
				progresspdf,
			),
			),
			)

			//check if biletranslation folder exists otherwise create
			if _, err := os.Stat(bt.OsPaths.Outputpath); !os.IsNotExist(err) {
				if err := os.RemoveAll(bt.OsPaths.Outputpath); err != nil {
					bt.LogError("remove dir "+bt.OsPaths.Outputpath, err)
				}
			}

			if err := os.Mkdir(bt.OsPaths.Outputpath, 0777); err != nil {
				bt.LogError("create dir "+bt.OsPaths.Outputpath, err)
			}

			//create subfolders of Bibletranslation
			if err := os.Mkdir(filepath.Join(bt.OsPaths.Outputpath, "html"), 0777); err != nil {
				bt.LogError("create dir "+filepath.Join(bt.OsPaths.Outputpath, "html"), err)
			}

			if err := os.Mkdir(filepath.Join(bt.OsPaths.Outputpath, "txt"), 0777); err != nil {
				bt.LogError("create dir "+filepath.Join(bt.OsPaths.Outputpath, "txt"), err)
			}

			mainVerses, err := bibleVerses.GetMainVerseText(mainTranslation)
			if err != nil {
				bt.LogError("get verse text in maintranslation", err)
			}

			// var lst_translationtext = make(models.Paragraphs, 0, 40)
			// var lst_translation = make([]string, 0, 40)

			translationVerses := bt.GetTranslationVerses(bibleVerses, bt.FilteredTranslations()...)

			var i float64
			bt.DocumentProgress = func(title string, proc float64) {
				i += proc
				fmt.Println(345, docprogress.Max, title, docprogress.Value, i)
				fyne.Do(func() {
					docprogress.SetValue(i)

				})
			}

			if bt.GetSameDocument() {
				docprogress.Max = float64(mainVerses.GetVerseAmount() + translationVerses.GetVerseAmount())
				bt.WriteSameTextFile(mainVerses, translationVerses)
				// bt.Writesamehtmlfile(text_main, lst_translationtext, "Main "+bt.GetMaintranslation(), w4document, docprogress)

				// //processbar info
				// w4document.SetText(documentname)
				// progress.Max = 1
				// currentcount += 1
				// progress.SetValue(float64(currentcount))

			} else {
				//write seperate files
				// fmt.Println(1, len(text_main))
				// fmt.Println(2, len(lst_translationtext))
				// docprogress.Max = float64(len(text_main) + len(lst_translationtext))

				// // write to file
				// bt.WriteTextFile(bt.GetMaintranslation(), text_main)

				// bt.Writehtmlfile(text_main, "Main "+bt.GetMaintranslation(), &wg)

				// //processbar info
				// // translation count for process bar
				// progress.Max = float64(len(lst_translation) + 1)
				// w4document.SetText(bt.GetMaintranslation())
				// progress.SetValue(float64(currentcount))

				//TODO:
				// for i, text := range lst_translationtext {
				// 	wg.Wait()
				// 	bt.WriteTextFile(lst_translation[i], text...)

				// 	bt.Writehtmlfile(text, "Translation "+lst_translation[i], &wg)

				// 	//processbar info
				// 	w4document.SetText(lst_translation[i])
				// 	currentcount += 1
				// 	progress.SetValue(float64(currentcount))
				// }

				w5.Show()

				// wg.Add(len(lst_translation) + 1)
				// go func() {
				// 	var progressvalue float64
				// 	for {

				// 		progresspdf.Max = float64(len(lst_translationtext) + 1)
				// 		progressvalue = progresspdf.Max - float64(basic.WaitingQue)
				// 		if progressvalue > 0 {
				// 			fmt.Println(basic.WaitingQue)
				// 			w5label.SetText("make pdf's")
				// 			progresspdf.SetValue(float64(progressvalue))
				// 		}
				// 		break

				// 	}
				// }()

				// wg.Wait()

				// if len(lst_translation) > 1 {
				// 	//make one pdf with all translations
				// 	bt.CombinePDF()
				// }

			}

			//close bibletool and clean up
			if err := bt.Close(); err != nil {
				panic(err)
			}

			mainWindow.Close()
		}()
	})

	verse_entry.SetPlaceHolder("Enter bible verses here, like:\nLuke 10.1\nJoh 3.1-14\nPsalm 2.3, 4.5-7\n")
	if verses := bt.GetVerses(); verses != "" {
		verse_entry.SetText(verses)
	}

	// set fyne
	mainWindow.SetContent(
		//set window content
		container.NewGridWithColumns(1,
			//set main content above verse field
			container.NewVBox(

				container.NewGridWithRows(2,
					//set top row with label and checkboxes all in one document
					container.NewGridWithColumns(3,
						label1,
						label2,
						checkbox_document,
					),
					container.NewCenter(sel_alltrans),
				),
				// set 3rd column with
				container.NewGridWithColumns(3,
					// set 1st column with main translation
					container.NewVBox(selection),
					// set 2nd column translation checkboxes
					container.NewVScroll(setofcheck),
					// set 3rd column with sermon title, pastor name and translate button
					container.NewVBox(container.NewGridWithRows(6,
						label3,
						sermonname,
						label4,
						pastorname),
						b1,
					),
				),
			),
			//set verse text field
			container.NewStack(verse_entry),
		),
	)
	// open window
	mainWindow.ShowAndRun()
}
