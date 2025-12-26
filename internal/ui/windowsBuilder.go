package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WindowBuilder struct {
	GetMainTranslation func(set string)
	mainTranslation    *widget.Select
	GetTranslations    func(set []string)
	translations       *widget.CheckGroup
	allTranslations    *widget.Check
	GetSameDocument    func(set bool)
	sameDocument       *widget.Check
	sermonTitle        *widget.Entry
	pastorName         *widget.Entry
	Translate          func()
	translateButton    *widget.Button
	verseField         *widget.Entry
	icon               fyne.Resource
}

func NewWindowBuilder(translations []string) *WindowBuilder {
	wb := WindowBuilder{}

	// main translation option
	wb.mainTranslation = widget.NewSelect(translations, func(set string) {
		if wb.GetSameDocument != nil {
			wb.GetMainTranslation(set)
		}
	})

	// translation select checkboxes
	wb.translations = widget.NewCheckGroup(translations, func(set []string) {
		if wb.GetTranslations != nil {
			wb.GetTranslations(set)
		}
	})

	//select all translations
	wb.allTranslations = widget.NewCheck("Select all", func(b bool) {
		t := []string{}
		if b {
			t = translations
		}
		wb.translations.SetSelected(t)
		wb.translations.Refresh()
	})

	// same document checkbox
	wb.sameDocument = widget.NewCheck("All in one document", func(set bool) {
		if wb.GetSameDocument != nil {
			wb.GetSameDocument(set)
		}
	})

	// sermontitle
	wb.sermonTitle = widget.NewEntry()
	wb.sermonTitle.PlaceHolder = "Enter Sermon title"

	// pastor name
	wb.pastorName = widget.NewEntry()
	wb.pastorName.PlaceHolder = "Enter Name of pastor"

	//translate button
	wb.translateButton = widget.NewButton("Translate", func() {
		if wb.Translate != nil {
			wb.Translate()
		}
	})

	// verse entry field
	wb.verseField = widget.NewMultiLineEntry()
	wb.verseField.SetPlaceHolder("Enter bible verses here, like:\nLuke 10.1\nJoh 3.1-14\nPsalm 2.3, 4.5-7\n")

	return &wb
}

func (wb *WindowBuilder) SetIconResource(path string) (err error) {
	wb.icon, err = fyne.LoadResourceFromPath(path)
	return
}

func (wb *WindowBuilder) SetMaintranslation(value string) {
	wb.mainTranslation.SetSelected(value)
}

func (wb *WindowBuilder) GetMaintranslationIndex() int {
	return wb.mainTranslation.SelectedIndex()
}

func (wb *WindowBuilder) SetTranslations(value []string) {
	wb.translations.SetSelected(value)
}

func (wb *WindowBuilder) SetSameDocument(value bool) {
	wb.sameDocument.SetChecked(value)
}

func (wb *WindowBuilder) SetSermonTitle(value string) {
	wb.sermonTitle.SetText(value)
}

func (wb *WindowBuilder) GetSermonTitle() string {
	return wb.sermonTitle.Text
}

func (wb *WindowBuilder) SetPastorName(value string) {
	wb.pastorName.SetText(value)
}

func (wb *WindowBuilder) GetPastorName() string {
	return wb.pastorName.Text
}

func (wb *WindowBuilder) SetVerseEntries(value string) {
	if value == "" {
		return
	}
	wb.verseField.SetText(value)
}

func (wb *WindowBuilder) GetVerseEntries() string {
	return wb.verseField.Text
}

func (wb *WindowBuilder) DisableButton() {
	wb.translateButton.Disable()
}

func (wb *WindowBuilder) ButtonIsDisabled() bool {
	return wb.translateButton.Disabled()
}

func (wb *WindowBuilder) EnableButton() {
	wb.translateButton.Enable()
}

func (wb *WindowBuilder) BuildMainWindow(a fyne.App, appName string) fyne.Window {
	w := a.NewWindow(appName)
	w.SetMaster()
	w.SetIcon(wb.icon)

	w.SetContent(
		//set window content
		container.NewGridWithColumns(1,
			//set main content above verse field
			container.NewVBox(
				container.NewGridWithRows(2,
					//set top row with label and checkboxes all in one document
					container.NewGridWithColumns(3,
						widget.NewLabel("Select Main Translation"),
						widget.NewLabel("Select Translations"),
						wb.sameDocument,
					),
					container.NewCenter(wb.allTranslations),
				),
				// set 3rd column with
				container.NewGridWithColumns(3,
					// set 1st column with main translation
					container.NewVBox(wb.mainTranslation),
					// set 2nd column translation checkboxes
					container.NewVScroll(wb.translations),
					// set 3rd column with sermon title, pastor name and translate button
					container.NewVBox(container.NewGridWithRows(6,
						widget.NewLabel("Sermontitle:"),
						wb.sermonTitle,
						widget.NewLabel("Name of Pastor:"),
						wb.pastorName),
						wb.translateButton,
					),
				),
			),
			//set verse text field
			container.NewStack(wb.verseField),
		),
	)
	return w
}

func (wb *WindowBuilder) BuildDialogWindow(a fyne.App, appName string) fyne.Window {
	w := a.NewWindow(appName)
	w.SetIcon(wb.icon)
	w.Resize(fyne.NewSize(200, 200))
	w.SetCloseIntercept(func() { w.Hide() })
	w.CenterOnScreen()
	return w
}

func (wb *WindowBuilder) BuildMainProgressWindow(a fyne.App, appName string, content ...fyne.CanvasObject) fyne.Window {
	w := a.NewWindow(appName)
	w.SetIcon(wb.icon)
	w.Resize(fyne.NewSize(250, 250))
	w.CenterOnScreen()
	w.SetContent(
		container.NewCenter(
			container.NewVBox(content...),
		),
	)
	return w
}

func (wb *WindowBuilder) BuildPdfProgressWindow(a fyne.App, appName string) fyne.Window {
	w := a.NewWindow(appName)
	w.SetIcon(wb.icon)
	w.Resize(fyne.NewSize(200, 100))
	w.CenterOnScreen()
	return w
}
