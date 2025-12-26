package bibletool

import (
	"bibletool/bibletool/env"
	"bibletool/bibletool/models"
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"gitea.tecamino.com/paadi/html2pdf"
	"gitea.tecamino.com/paadi/html2pdf/converter"
	pdfModels "gitea.tecamino.com/paadi/html2pdf/models"
	"gitea.tecamino.com/paadi/pdfmerge"
)

func (bt *Bibletool) WriteTextFile(maintranslation *models.Translation, translations *models.Translations) error {
	defer bt.TotalProgressAdd(1)

	documentName := maintranslation.GetTranslationName()
	bt.DebugLog("WriteTextFile", "write text file "+documentName)
	f, err := os.Create(filepath.Join(bt.OutputDir, "txt", documentName+".txt"))
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
	documentName := maintranslation.GetTranslationName()

	bt.DebugLog("WriteHtmlfile", "write tmpl html file "+documentName)
	err := bt.WriteHtml(filepath.Join(bt.OutputDir, "html", documentName+".html"), models.HtmlStruct{
		Name:                "Main " + maintranslation.Name,
		SermonTitle:         bt.GetSermonTitle(),
		PastorName:          bt.GetPastor(),
		RightToLeftDocument: maintranslation.RightToLeft,
		MainTranslation:     maintranslation,
		Translations:        translations,
		Date:                time.Now().Format("02-January-2006"),
		CurrentPath:         template.URL(filepath.ToSlash(os.Args[0])),
		ProgressFnc:         bt.DocumentProgress,
		SameDocument:        sameDocument,
	})

	bt.TotalProgressAdd(1)
	if err != nil {
		bt.LogError("htmlbuilder", err)
	}
	return nil
}

func (bt *Bibletool) ConvertToPdf(documentNames ...string) error {
	var err error
	var c *converter.Converter

	chromePath := env.ChromePath.GetValue()
	bt.DebugLog("ConvertToPdf", "open chrome headless shell from "+chromePath)
	bt.Wg.Go(func() {
		c, err = html2pdf.NewConverterInstance(chromePath)
		if err != nil {
			bt.LogError("html2pdf", err)
		}

		c.SetProgressCallback(bt.PdfProgressAdd)

		var files []pdfModels.File
		for _, name := range documentNames {
			bt.DebugLog("ConvertToPdf", "convert "+name)
			files = append(files, pdfModels.File{
				Input:  filepath.Join(bt.OutputDir, "html", name+".html"),
				Output: filepath.Join(bt.OutputDir, name+".pdf"),
			})
		}

		err = c.Convert(files...)
		if err != nil {
			bt.LogError("html2pdf", err)
		}
	})
	return err
}

func (bt *Bibletool) CombinePDF() error {

	// get files in directory
	files, err := os.ReadDir(bt.OutputDir)
	if err != nil {
		bt.LogError("combine pdf", err)
	}

	//get list of all pdf's in output folder
	pdflist := []string{}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".pdf" {
			pdflist = append(pdflist, filepath.Join(bt.OutputDir, file.Name()))
		}
	}

	output := filepath.Join(bt.OutputDir, "AllTranslation.pdf")
	bt.DebugLog("CombinePDF", output)

	// merge them in one file
	err = pdfmerge.Pdfmerge(pdflist, output)
	if err != nil {
		bt.LogError("combine pdf", err)
	}
	return err
}
