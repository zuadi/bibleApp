package bibletool

import (
	"bibletool/bibletool/models"
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"gitea.tecamino.com/paadi/html2pdf"
	"gitea.tecamino.com/paadi/pdfmerge"
)

func (bt *Bibletool) WriteTextFile(maintranslation *models.Translation, translations *models.Translations) error {
	defer bt.TotalProgressAdd(1)

	documentName := "Translation " + maintranslation.Name + ".txt"
	if maintranslation.IsMain {
		documentName = "Main " + maintranslation.Name + ".txt"
	}
	f, err := os.Create(filepath.Join(bt.OsPaths.Outputpath, "txt", documentName))
	if err != nil {
		bt.LogError("write text file", err)
		return err
	}
	defer f.Close()

	// Create a writer
	w := bufio.NewWriter(f)

	for i, paragraph := range maintranslation.Paragraphs {
		if i > 0 {
			// add return between new verse
			_, err = w.WriteString("\n\n")
			if err != nil {
				bt.LogError("write text file", err)
			}
		}

		// write title/verse
		_, err = w.WriteString(paragraph.Title + "\n")
		if err != nil {
			bt.LogError("write text file", err)
		}

		for _, verse := range paragraph.Verse {
			bt.DocumentProgressAdd(paragraph.Title, 1)

			_, err = fmt.Fprintf(w, "%d %s\n", verse.Number, verse.Text)
			if err != nil {
				bt.LogError("write text file", err)
			}
		}

		if translations == nil {
			continue
		}

		// add return between translation
		_, err = w.WriteString("\n")
		if err != nil {
			bt.LogError("write text file", err)
		}

		//write translation verses
		for _, translation := range *translations {
			for _, verse := range translation.Paragraphs[i].Verse {
				bt.DocumentProgressAdd(paragraph.Title, 1)

				_, err = fmt.Fprintf(w, "%d %s\n", verse.Number, verse.Text)
				if err != nil {
					bt.LogError("write text file", err)
				}

				// add return between translation
				_, err = w.WriteString("\n")
				if err != nil {
					bt.LogError("write text file", err)
				}
			}
		}
	}

	// Very important to invoke after writing a large number of lines
	err = w.Flush()
	if err != nil {
		bt.LogError("write text file", err)
	}
	return err
}

func (bt *Bibletool) WriteHtmlfile(maintranslation *models.Translation, translations *models.Translations, sameDocument bool) error {
	var documentName string
	if maintranslation.IsMain {
		documentName = "Main " + maintranslation.Name
	} else {
		documentName = "Translation " + maintranslation.Name

	}

	err := bt.WriteHtml(filepath.Join(bt.OsPaths.Outputpath, "html", documentName+".html"), models.HtmlStruct{
		Name:                "Main " + maintranslation.Name,
		SermonTitle:         bt.GetSermonTitle(),
		PastorName:          bt.GetPastor(),
		RightToLeftDocument: maintranslation.RightToLeft,
		MainTranslation:     maintranslation,
		Translations:        translations,
		Date:                time.Now().Format("02-January-2006"),
		CurrentPath:         template.URL(filepath.ToSlash(bt.OsPaths.Currentdirectory)),
		ProgressFnc:         bt.DocumentProgress,
		SameDocument:        sameDocument,
	})

	bt.TotalProgressAdd(1)
	if err != nil {
		bt.LogError("htmlbuilder", err)
	}

	bt.Wg.Go(func() {
		defer bt.PdfProgressAdd(1)
		bt.PdfProgressAdd(1)
		err = html2pdf.Convert("assets", filepath.Join(bt.OsPaths.Outputpath, "html", documentName+".html"), filepath.Join(bt.OsPaths.Outputpath, documentName+".pdf"))
		if err != nil {
			bt.LogError("html2pdf", err)
		}
	})
	return nil
}

func (bt *Bibletool) CombinePDF() error {

	// get files in directory
	files, err := os.ReadDir(bt.OsPaths.Outputpath)
	if err != nil {
		bt.LogError("combine pdf", err)
	}

	//get list of all pdf's in output folder
	var pdflist = make([]string, 0, 40)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".pdf" {
			pdflist = append(pdflist, filepath.Join(bt.OsPaths.Outputpath, file.Name()))
		}
	}

	// merge them in one file
	err = pdfmerge.Pdfmerge(pdflist, filepath.Join(bt.OsPaths.Outputpath, "AllTranslation.pdf"))
	if err != nil {
		bt.LogError("combine pdf", err)
	}
	return err
}
