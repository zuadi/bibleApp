package main

//This script is a bibletranslation that gives out a .txt file for using the text further and as well as in pdf.
// The maintranslation has to be choosen and further translation can be checked, the checkbox same document includes all text in one txt and pdf document
import (
	"bibletool/bibletool"
	"bibletool/bibletool/env"
	"bibletool/internal/ui"
	"bibletool/utils"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	// This function is a bibletranslation help that there can be a main translation choosen and the desired verses entered.
	// it will create a txt and pdf file of each translation or a combined file with all translation.
	// it checks if verse exists in main translation and if not give out the not found verses

	bt, err := bibletool.NewBibletool()
	if err != nil {
		panic(err)
	}

	// get all bible translation
	allTranslations, err := bt.GetAllTranslations()
	if err != nil {
		panic(err)
	}

	bt.DebugLog("main", "initiate new fyne app")

	// make new fyne app
	app := app.New()

	bt.DebugLog("main", "initiate fyne windows builder")
	wb := ui.NewWindowBuilder(allTranslations)
	bt.DebugLog("main", "set icon resource path")
	wb.SetIconResource(env.IconPath.GetValue())
	bt.DebugLog("main", "set last user entered value from config file")
	wb.SetMaintranslation(bt.GetMaintranslation())
	wb.SetTranslations(bt.GetSelectedTranslations())
	wb.SetSameDocument(bt.GetSameDocument())
	wb.SetPastorName(bt.GetPastor())
	wb.SetSermonTitle(bt.GetSermonTitle())
	wb.SetVerseEntries(bt.GetVerses())

	bt.DebugLog("main", "initiate all callback functions for saving values to config file")

	//call back for main translation
	wb.GetMainTranslation = func(set string) {
		bt.SetMaintranslation(set)
	}

	//call back for translations
	wb.GetTranslations = func(set []string) {
		bt.SetTranslations(set)
	}

	//call back for same document settings
	wb.GetSameDocument = func(set bool) {
		bt.SetSameDocument(set)
	}

	bt.DebugLog("main", "build main window")
	mainWindow := wb.BuildMainWindow(app, bt.AppName)

	var openWindow bool
	bt.DebugLog("main", "build warning dialog window")
	warningDialog := wb.BuildDialogWindow(app, bt.AppName)
	warningDialog.SetOnClosed(func() { openWindow = false })

	bt.DebugLog("main", "initiate fyne progress bars")
	docprogress := widget.NewProgressBar()
	progress := widget.NewProgressBar()
	//progressbar for pdf generating
	bt.DebugLog("main", "pdf progress window")
	pdfProgressWindow := wb.BuildPdfProgressWindow(app, bt.AppName)

	// start translating when 'translation' button pressed
	wb.Translate = func() {
		bt.DebugLog("main", "translating")

		mainTranslation := bt.GetMaintranslation()
		// warning error window still open
		if openWindow && mainTranslation == "" {
			warningDialog.RequestFocus()
			warningDialog.Show()
			return
		}

		go func() {
			defer wb.EnableButton()

			// check whether main translation is not set
			if mainTranslation == "" {
				openWindow = true
				fyne.Do(func() {
					warningDialog.SetContent(container.NewCenter(container.NewVBox(
						widget.NewLabel("No Maintranslation choosen"),
						widget.NewButton("Ok", func() { warningDialog.Hide() }),
					),
					),
					)
					warningDialog.Show()
				})
				return
			}

			/* start translating */

			// save sermon title
			bt.SetSermonTitle(wb.GetSermonTitle())

			//save pastor name
			bt.SetPastor(wb.GetPastorName())

			//get bible verses
			bibleVerses, notFound := bt.GetBibleVerses(wb.GetVerseEntries(), wb.GetMaintranslationIndex())
			if notFound.IsError {
				wb.BuildErrorWindow(mainWindow, notFound.Error)
				return
			} else if notFound.Error != nil {
				var errText string
				var labelText = "This Bibleverses were not found:"

				if notFound.Error.Error() != "no bibleverses entered" {
					errText = notFound.Error.Error()
				} else {
					labelText = "No bibleverses entered"
				}

				openWindow = true

				fyne.Do(func() {
					warningDialog.SetContent(container.NewCenter(container.NewVBox(
						widget.NewLabel(labelText),
						widget.NewLabel(errText),
						widget.NewButton("Ok", func() { warningDialog.Hide() }),
					),
					),
					)
					warningDialog.Show()
				})
				return
			}

			fyne.Do(func() {
				warningDialog.Close()
			})

			start := time.Now()
			wb.DisableButton()
			bt.SetVerses(wb.GetVerseEntries())
			//progressbar
			fyne.Do(func() {
				mainProgressWindow := wb.BuildMainProgressWindow(app, bt.AppName,
					widget.NewLabel("Translating for you:"),
					widget.NewLabel("In progress..."),
					docprogress,
					widget.NewLabel("Total progress"),
					progress)
				mainProgressWindow.Show()
			})

			progresspdf := widget.NewProgressBar()

			fyne.Do(func() {
				pdfProgressWindow.SetContent(container.NewCenter(container.NewVBox(
					widget.NewLabel("Make pdf's"),
					progresspdf,
				),
				),
				)
			})

			// set call back for total progress
			var total float64
			bt.TotalProgress = func(p float64) {
				total += p
				fyne.Do(func() {
					progress.SetValue(total)

				})
			}

			// set call back for file generating progress
			var document float64
			bt.DocumentProgress = func(title string, p float64) {
				document += p
				fyne.Do(func() {
					docprogress.SetValue(document)
				})
			}

			// set call back for pdf convertion progress
			bt.PdfProgress = func(p float64) {
				fyne.Do(func() {
					progresspdf.SetValue(p)
				})
			}

			//check if biletranslation folder exists otherwise create
			if err := utils.MkDirs(bt.OutputDir, "html", "txt"); err != nil {
				bt.LogError("create dir and subdirs "+bt.OutputDir, err)
				wb.BuildErrorWindow(mainWindow, err)
				return
			}

			mainVerses, err := bibleVerses.GetMainVerseText(mainTranslation)
			if err != nil {
				bt.LogError("get verse text in maintranslation", err)
				wb.BuildErrorWindow(mainWindow, err)
				return
			}

			translationVerses, err := bt.GetTranslationVerses(bibleVerses, bt.FilteredTranslations()...)
			if err != nil {
				bt.LogError("get translation verses", err)
				wb.BuildErrorWindow(mainWindow, err)
				return
			}

			// set max progress scale
			docprogress.Max = 2 * float64((mainVerses.GetVerseAmount() + translationVerses.GetVerseAmount()))

			if bt.GetSameDocument() {
				progress.Max = 2
				progresspdf.Max = float64(+len(*translationVerses))

				if err := bt.WriteTextFile(mainVerses, translationVerses); err != nil {
					wb.BuildErrorWindow(mainWindow, err)
					return
				}

				if err := bt.WriteHtmlfile(mainVerses, translationVerses, true); err != nil {
					wb.BuildErrorWindow(mainWindow, err)
					return
				}

				names := []string{mainVerses.Name}
				for _, translation := range *translationVerses {
					names = append(names, translation.Name)
				}

				if err := bt.ConvertToPdf(names...); err != nil {
					wb.BuildErrorWindow(mainWindow, err)
					return
				}
			} else {
				//write seperate files
				progress.Max = 1 + (float64(+len(*translationVerses)))*2
				progresspdf.Max = float64(+len(*translationVerses))

				// // write to file
				names := []string{mainVerses.GetTranslationName()}

				if err := bt.WriteTextFile(mainVerses, nil); err != nil {
					wb.BuildErrorWindow(mainWindow, err)
					return
				}

				if err := bt.WriteHtmlfile(mainVerses, nil, false); err != nil {
					wb.BuildErrorWindow(mainWindow, err)
					return
				}

				for _, t := range *translationVerses {
					if err := bt.WriteTextFile(t, nil); err != nil {
						wb.BuildErrorWindow(mainWindow, err)
						return
					}

					if err := bt.WriteHtmlfile(t, nil, false); err != nil {
						wb.BuildErrorWindow(mainWindow, err)
						return
					}

					names = append(names, t.GetTranslationName())
				}

				if err := bt.ConvertToPdf(names...); err != nil {
					wb.BuildErrorWindow(mainWindow, err)
					return
				}

				fyne.Do(func() {
					pdfProgressWindow.Show()
				})
			}

			// wait till all jobs are done
			bt.Wg.Wait()

			//close bibletool and clean up
			if err := bt.Close(); err != nil {
				wb.BuildErrorWindow(mainWindow, err)
				return
			}

			fyne.Do(func() {
				mainWindow.Close()
			})
			fmt.Println("used time:", time.Since(start))
		}()
	}

	// open window
	mainWindow.ShowAndRun()
}
