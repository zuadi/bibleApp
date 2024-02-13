package output

import (
	"C"
	"bufio"
	"os"
	"strings"

	"time"

	"bibletool/Modules"
	"bibletool/basic"
	"embed"
	"fmt"
	"runtime"
	"unicode"

	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)
import (
	"bibletool/pdfmerge"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
)

//go:embed pics/*
var picsfiles embed.FS

func Writetextfile(paths *basic.OSPaths, in_text Modules.OutputText, filename string, docprocess *widget.ProgressBar) {
	// Create a txt file for writing
	f, _ := os.Create(paths.Outputpath + paths.Pathseperator + "txt" + paths.Pathseperator + filename + ".txt")

	// Create a writer
	w := bufio.NewWriter(f)

	docprocess.Max = float64(len(in_text))

	for i := range in_text {
		docprocess.SetValue(float64(i) + 1)

		w.WriteString(in_text[i].Titel + "\n")

		for ii := range in_text[i].Verse {
			var nodublicate []string
			inResult := make(map[string]bool)
			// remove dublicates
			for _, str := range in_text[i].Verse[ii] {
				if _, ok := inResult[str]; !ok {
					inResult[str] = true
					nodublicate = append(nodublicate, str)
				}
			}
			// write verses
			for _, s := range nodublicate {
				w.WriteString(s + "\n")
			}

		}
		if i < len(in_text)-1 {
			// add return between new verse
			w.WriteString("\n\n")
		}
	}

	// Very important to invoke after writing a large number of lines
	w.Flush()

}

func Writesamedoctext(path *basic.OSPaths, maintext Modules.OutputText, lst_transtext []Modules.OutputText, filename string, w4document *widget.Label, docprogress *widget.ProgressBar) {

	docprogress.Max = float64(len(maintext)) - 1

	// Create a txt file for writing
	f, _ := os.Create(path.Outputpath + path.Pathseperator + "txt" + path.Pathseperator + filename + ".txt")

	// Create a writer
	w := bufio.NewWriter(f)

	for i := range maintext {

		w4document.SetText(maintext[i].Titel)
		docprogress.SetValue(float64(i) + 1)

		// write title/verse
		w.WriteString(maintext[i].Titel + "\n")

		// remove dublicates
		for ii := range maintext[i].Verse {
			var nodublicate []string
			inResult := make(map[string]bool)
			// remove dublicates
			for _, str := range maintext[i].Verse[ii] {
				if _, ok := inResult[str]; !ok {
					inResult[str] = true
					nodublicate = append(nodublicate, str)
				}
			}

			// write main verses
			for _, s := range nodublicate {
				w.WriteString(s + "\n")
			}
		}
		if len(lst_transtext) > 0 {
			// add return between translation
			w.WriteString("\n")
		}

		// write translation verses
		for _, translation := range lst_transtext {

			for numb, text := range translation {
				if numb > i {
					break
				} else if numb < i {
					continue
				} else if numb == i {
					for _, verse := range text.Verse {
						var nodublicate2 []string
						inResult := make(map[string]bool)
						// remove dublicates
						for _, str := range verse {
							if _, ok := inResult[str]; !ok {
								inResult[str] = true
								nodublicate2 = append(nodublicate2, str)
							}
						}
						for _, item := range nodublicate2 {
							w.WriteString(item + "\n")
						}
					}
				}
				// add return between translation
				w.WriteString("\n")
			}
		}
		if i < len(maintext)-1 {
			// add return between new verse
			w.WriteString("\n\n")
		}
	}

	// Very important to invoke after writing a large number of lines
	w.Flush()
}

func Writehtmlfile(in_text Modules.OutputText, translationtype string, filename string, sermon_title string, pastor_name string, paths *basic.OSPaths, docprogress *widget.ProgressBar, wg *sync.WaitGroup) {

	//remove header and footer
	if _, err := os.Stat(paths.Tempdir + paths.Pathseperator + filename + "header.html"); err == nil {
		os.Remove(paths.Tempdir + paths.Pathseperator + filename + "header.html")
	}
	if _, err := os.Stat(paths.Tempdir + paths.Pathseperator + filename + "footer.html"); err == nil {
		os.Remove(paths.Tempdir + paths.Pathseperator + filename + "footer.html")
	}
	readdirection := "<html>\n"
	if strings.Contains(filename, "Arabic") || strings.Contains(filename, "Hebrew") || strings.Contains(filename, "Persian") || strings.Contains(filename, "Aramaic") {
		readdirection = "<html dir=\"rtl\">\n"
	}

	filecontent := "<!DOCTYPE html>\n" +
		readdirection +
		"<head>\n" +
		"<meta charset=\"UTF-8\">\n" +
		"<title>" + translationtype + " " + filename + "</title>\n" +
		"<style>\n" +
		"body {\n" +
		"\tpadding-top: 30px;\n" +
		"\tmargin-top: 10px;\n" +
		"\tfont-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"\tpadding-bottom: 10px;\n" +
		"}\n\n" +
		"h1 {\n" +
		"\tfont-size: 18pt;\n" +
		"\tpadding-right: 30px;\n" +
		"\ttext-align: right;\n" +
		"\tfloat: right;\n" +
		"\tcolor: rgb(88, 80, 99);\n" +
		"\tfont-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"}\n\n" +
		"h2 {\n" +
		"\tfont-size: 16pt;\n" +
		"\tcolor: rgb(88, 80, 99);\n" +
		"\tfloat: left;\n" +
		"}\n\n" +
		"h3 {\n" +
		"\tfont-size: 16pt;\n" +
		"\tcolor: rgb(199, 99, 99);\n" +
		"\tfont-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"\tmargin-left: 2px;\n" +
		"\tpage-break-after: avoid;\n" +
		"}\n\n" +
		"p {\n" +
		"\tcolor: black;\n" +
		"\tfont-size: 16pt;\n" +
		"\tfont-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"\tmargin-left: 12px;\n" +
		"\tmargin-top: -20px;\n" +
		"\tmargin-bottom: 20px;\n" +
		"\tfont-weight: normal;\n" +
		"\tpage-break-inside: avoid;\n" +
		"}\n" +
		"</style>\n" +
		"</head>\n" +
		"<body style=\"border:5; margin: 120;\" onload=\"subst()\">\n"

	//make header file
	headercontent := "<!doctype html>\n" +
		"<html>\n" +
		"<head>\n" +
		"	<meta charset=\"utf-8\">\n" +
		"	<script>\n" +
		"		function substitutePdfVariables() {\n\n" +
		"			function getParameterByName(name) {\n" +
		"				var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);\n" +
		"				return match && decodeURIComponent(match[1].replace(/\\+/g, ' '));\n" +
		"			}\n\n" +
		"			function substitute(name) {\n" +
		"				var value = getParameterByName(name);\n" +
		"				var elements = document.getElementsByClassName(name);\n\n" +
		"				for (var i = 0; elements && i < elements.length; i++) {\n" +
		"					elements[i].textContent = value;\n" +
		"				}\n" +
		"			}\n\n" +
		"			['frompage', 'topage', 'page', 'webpage', 'section', 'subsection', 'subsubsection']\n" +
		"				.forEach(function(param) {\n" +
		"					substitute(param);\n" +
		"				});\n" +
		"		}\n" +
		"	</script>\n" +
		"  <style>\n" +
		"	.container{\n" +
		"	padding-top: 15px;\n" +
		"	padding-bottom: 130px;\n" +
		"   margin-bottom: 130px;\n" +
		"	align-items: center;\n" +
		"	justify-content: center;\n" +
		"}\n" +
		"img {\n" +
		"	padding-right: 20px;\n" +
		"	float: left;\n" +
		"}\n\n" +
		"h1 {\n" +
		"	font-size: 18pt;\n" +
		"	padding-right: 30px;\n" +
		"	text-align: right;\n" +
		"	float: right;\n" +
		"	color: rgb(88, 80, 99);\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"}\n" +
		"h2 {\n" +
		"	font-size: 16pt;\n" +
		"	color: rgb(88, 80, 99);\n" +
		"	float: left;\n" +
		"	font-weight: normal;\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"}\n" +
		"	p {\n" +
		"	  text-align: right;\n" +
		"	  size: 16pt;\n" +
		"	  font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	}\n" +
		"  </style>\n" +
		"</head>\n" +
		"<body onload=\"substitutePdfVariables()\">\n" +
		"	<div class=\"container\">\n" +
		"		<div class=\"image\">\n" +
		`			<img src="` + paths.Currentdirectory + "pics" + paths.Pathseperator + "pottershouse.png" + `" style="width:75px;height:65px;">` + "\n" +
		"			<h2>" + sermon_title + "</h2>\n" +
		"			<h1>" + translationtype + " " + filename + "</h1>\n" +
		"		</div>\n" +
		"</body>\n" +
		"</html>"

	// save as header html file
	htmlfile, err := os.OpenFile(paths.Tempdir+paths.Pathseperator+filename+"header.html", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		basic.CheckErr(err, "Error creating header.html file")
	} else {
		htmlfile.WriteString(headercontent)
	}
	htmlfile.Close()

	//make header file
	footercontent := "<!doctype html>\n" +
		"<html>\n" +
		"<head>\n" +
		"	<meta charset=\"utf-8\">\n" +
		"	<script>\n" +
		"		function substitutePdfVariables() {\n\n" +
		"			function getParameterByName(name) {\n" +
		"				var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);\n" +
		"				return match && decodeURIComponent(match[1].replace(/\\+/g, ' '));\n" +
		"			}\n\n" +
		"			function substitute(name) {\n" +
		"				var value = getParameterByName(name);\n" +
		"				var elements = document.getElementsByClassName(name);\n\n" +
		"				for (var i = 0; elements && i < elements.length; i++) {\n" +
		"					elements[i].textContent = value;\n" +
		"				}\n" +
		"			}\n\n" +
		"			['frompage', 'topage', 'page', 'webpage', 'section', 'subsection', 'subsubsection']\n" +
		"				.forEach(function(param) {\n" +
		"					substitute(param);\n" +
		"				});\n" +
		"		}\n" +
		"	</script>\n" +
		"  <style>\n" +
		"	p {\n" +
		"	  text-align: right;\n" +
		"	  size: 16pt;\n" +
		"	  font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	}\n" +
		"#footer {\n" +
		"	width: 99%;\n" +
		"	height: 80px;\n" +
		"	white-space: nowrap;\n" +
		"	}\n" +
		"	.alignleft {\n" +
		"	display: inline-block;\n" +
		"	width: 33%;\n" +
		"	height: 80%;\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	padding-top: 10px;\n" +
		"	}\n" +
		"	.aligncenter {\n" +
		"	display: inline-block;\n" +
		"	text-align: center;\n" +
		"	width: 33%;\n" +
		"	height: 80%;\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	}\n" +
		"	.alignright {\n" +
		"	display: inline-block;\n" +
		"	text-align: right;\n" +
		"	width: 33%;\n" +
		"	height: 80%;\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	}\n" +
		"  </style>\n" +
		"</head>\n" +
		"<body onload=\"substitutePdfVariables()\">\n" +
		"	<div id=\"footer\">\n" +
		"		<div class=\"alignleft\">" + time.Now().Format("02-January-2006") + "</p></div>\n" +
		"		<div class=\"aligncenter\">" + pastor_name + "</p></div>\n" +
		"		<div class=\"alignright\"> <span class=\"page\"></span>/<span class=\"topage\"></span></p></div>\n" +
		"</div>\n" +
		"</body>\n" +
		"</html>"

	// save as footer html file
	htmlfile, err = os.OpenFile(paths.Tempdir+paths.Pathseperator+filename+"footer.html", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		basic.CheckErr(err, "Error creating footer.html file")
	} else {
		htmlfile.WriteString(footercontent)
	}
	htmlfile.Close()

	docprogress.Max = float64(len(in_text)) * 1.2

	for i := range in_text {

		docprogress.SetValue(float64(i) + 1)

		//write bibleverses remove first dublicates
		for ii := range in_text[i].Verse {
			var nodublicate []string
			inResult := make(map[string]bool)
			// remove dublicates
			for _, str := range in_text[i].Verse[ii] {

				if _, ok := inResult[str]; !ok {
					inResult[str] = true
					nodublicate = append(nodublicate, str)
				}
			}
			// write verses
			var lst_versetext []string
			for _, s := range nodublicate {

				if len(s) < 90 {
					lst_versetext = append(lst_versetext, s)
				} else {
					var first_prt string
					var second_prt string
					var third_prt string
					var fourth_prt string
					var fifth_prt string
					var sixth_prt string
					var seventh_prt string
					var eighth_prt string
					var ninth_prt string
					var tenth_prt string

					var length int
					split_temp := strings.Split(s, " ")

					//check for to long line and separate it
					for _, elem := range split_temp {
						//add length of element and plus 1 for whitespace
						length += len(elem) + 1
						if length < 90 {
							first_prt = first_prt + elem + " "
						} else if length < 180 {
							second_prt = second_prt + elem + " "
						} else if length < 270 {
							third_prt = third_prt + elem + " "
						} else if length < 360 {
							fourth_prt = fourth_prt + elem + " "
						} else if length < 450 {
							fifth_prt = fifth_prt + elem + " "
						} else if length < 540 {
							sixth_prt = sixth_prt + elem + " "
						} else if length < 630 {
							seventh_prt = seventh_prt + elem + " "
						} else if length < 720 {
							eighth_prt = eighth_prt + elem + " "
						} else if length < 810 {
							ninth_prt = ninth_prt + elem + " "
						} else {
							tenth_prt = tenth_prt + elem + " "
						}

					}

					//put lines in,list
					if len(first_prt) > 0 {
						lst_versetext = append(lst_versetext, first_prt)
					}
					if len(second_prt) > 0 {
						lst_versetext = append(lst_versetext, second_prt)
					}
					if len(third_prt) > 0 {
						lst_versetext = append(lst_versetext, third_prt)
					}
					if len(fourth_prt) > 0 {
						lst_versetext = append(lst_versetext, fourth_prt)
					}
					if len(fifth_prt) > 0 {
						lst_versetext = append(lst_versetext, fifth_prt)
					}
					if len(sixth_prt) > 0 {
						lst_versetext = append(lst_versetext, sixth_prt)
					}
					if len(seventh_prt) > 0 {
						lst_versetext = append(lst_versetext, seventh_prt)
					}
					if len(eighth_prt) > 0 {
						lst_versetext = append(lst_versetext, eighth_prt)
					}
					if len(ninth_prt) > 0 {
						lst_versetext = append(lst_versetext, ninth_prt)
					}
					if len(tenth_prt) > 0 {
						lst_versetext = append(lst_versetext, tenth_prt)
					}
				}
			}

			//write title of bibleverse
			filecontent += "<h3>" + in_text[i].Titel + "</h3>\n"

			//write lines to html
			for _, element := range lst_versetext {
				filecontent += "<p>" + element + "</p>\n"
			}
		}
	}
	filecontent += "</body>\n</html>"

	// save as html file
	htmlfile, err = os.OpenFile(paths.Outputpath+paths.Pathseperator+"html"+paths.Pathseperator+filename+".html", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		basic.CheckErr(err, "Error creating "+filename+" .html file")
	} else {
		htmlfile.WriteString(filecontent)
	}
	htmlfile.Close()

	var stop bool
	var counter float64
	go func() {
		for {
			if stop || docprogress.Max*0.8+counter == docprogress.Max {
				break
			}
			counter += 1
			docprogress.SetValue(docprogress.Max*0.8 + counter)
			time.Sleep(time.Second * 1)
		}
	}()

	var wkhtmltopdfstring string
	switch runtime.GOOS {
	case "windows":
		wkhtmltopdfstring = "wkhtmltopdf.exe"
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			wkhtmltopdfstring = "wkhtmltopdfdarwinarm64"
		case "amd64":
			wkhtmltopdfstring = "wkhtmltopdfdarwinamd64"
		}
	}
	//create pdf

	go basic.StartcommandwithArgsWG(paths.Currentdirectory+wkhtmltopdfstring, []string{"--header-html", paths.Tempdir + paths.Pathseperator + filename + "header.html", "--footer-html", paths.Tempdir + paths.Pathseperator + filename + "footer.html", paths.Outputpath + paths.Pathseperator + "html" + paths.Pathseperator + filename + ".html", paths.Outputpath + paths.Pathseperator + filename + ".pdf"}, wg)

	stop = true
	docprogress.SetValue(docprogress.Max)
}

func Writesamehtmlfile(maintext Modules.OutputText, lst_transtext []Modules.OutputText, filename string, sermon_title string, pastor_name string, paths *basic.OSPaths, w4document *widget.Label, docprocess *widget.ProgressBar) {

	filecontent := "<!DOCTYPE html>\n" +
		"<html>\n" +
		"<head>\n" +
		"<meta charset=\"UTF-8\">\n" +
		"<title>" + filename + "</title>\n" +
		"<style>\n" +
		"body {\n" +
		"\tpadding-top: 30px;\n" +
		"\tmargin-top: 10px;\n" +
		"\tfont-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"\tpadding-bottom: 10px;\n" +
		"}\n\n" +
		"h1 {\n" +
		"\tfont-size: 18pt;\n" +
		"\tpadding-right: 30px;\n" +
		"\ttext-align: right;\n" +
		"\tfloat: right;\n" +
		"\tcolor: rgb(88, 80, 99);\n" +
		"\tfont-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"}\n\n" +
		"h2 {\n" +
		"\tfont-size: 16pt;\n" +
		"\tcolor: rgb(88, 80, 99);\n" +
		"\tfloat: left;\n" +
		"}\n\n" +
		"h3 {\n" +
		"\tfont-size: 16pt;\n" +
		"\tcolor: rgb(199, 99, 99);\n" +
		"\tfont-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"\tmargin-left: 2px;\n" +
		"\tpage-break-after: avoid;\n" +
		"}\n\n" +
		"p {\n" +
		"\tcolor: black;\n" +
		"\tfont-size: 16pt;\n" +
		"\tfont-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"\tmargin-left: 12px;\n" +
		"\tmargin-top: -20px;\n" +
		"\tmargin-bottom: 20px;\n" +
		"\tfont-weight: normal;\n" +
		"\tpage-break-inside: avoid;\n" +
		"}\n" +
		"</style>\n" +
		"</head>\n" +
		"<body style=\"border:5; margin: 120;\" onload=\"subst()\">\n"

	//make header file
	headercontent := "<!doctype html>\n" +
		"<html>\n" +
		"<head>\n" +
		"	<meta charset=\"utf-8\">\n" +
		"	<script>\n" +
		"		function substitutePdfVariables() {\n\n" +
		"			function getParameterByName(name) {\n" +
		"				var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);\n" +
		"				return match && decodeURIComponent(match[1].replace(/\\+/g, ' '));\n" +
		"			}\n\n" +
		"			function substitute(name) {\n" +
		"				var value = getParameterByName(name);\n" +
		"				var elements = document.getElementsByClassName(name);\n\n" +
		"				for (var i = 0; elements && i < elements.length; i++) {\n" +
		"					elements[i].textContent = value;\n" +
		"				}\n" +
		"			}\n\n" +
		"			['frompage', 'topage', 'page', 'webpage', 'section', 'subsection', 'subsubsection']\n" +
		"				.forEach(function(param) {\n" +
		"					substitute(param);\n" +
		"				});\n" +
		"		}\n" +
		"	</script>\n" +
		"  <style>\n" +
		"	.container{\n" +
		"	padding-top: 15px;\n" +
		"	padding-bottom: 130px;\n" +
		"   margin-bottom: 130px;\n" +
		"	align-items: center;\n" +
		"	justify-content: center;\n" +
		"}\n" +
		"img {\n" +
		"	padding-right: 20px;\n" +
		"	float: left;\n" +
		"}\n\n" +
		"h1 {\n" +
		"	font-size: 18pt;\n" +
		"	padding-right: 30px;\n" +
		"	text-align: right;\n" +
		"	float: right;\n" +
		"	color: rgb(88, 80, 99);\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"}\n" +
		"h2 {\n" +
		"	font-size: 16pt;\n" +
		"	color: rgb(88, 80, 99);\n" +
		"	float: left;\n" +
		"	font-weight: normal;\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"}\n" +
		"	p {\n" +
		"	  text-align: right;\n" +
		"	  size: 16pt;\n" +
		"	  font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	}\n" +
		"  </style>\n" +
		"</head>\n" +
		"<body onload=\"substitutePdfVariables()\">\n" +
		"	<div class=\"container\">\n" +
		"		<div class=\"image\">\n" +
		`			<img src="` + paths.Currentdirectory + "pics" + paths.Pathseperator + "pottershouse.png" + `" style="width:75px;height:65px;">` + "\n" +
		"			<h2>" + sermon_title + "</h2>\n" +
		"			<h1>" + filename + "</h1>\n" +
		"		</div>\n" +
		"</body>\n" +
		"</html>"

	// save as header html file
	htmlfile, err := os.OpenFile(paths.Tempdir+"header.html", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		basic.CheckErr(err, "Error creating header.html file")
	} else {
		htmlfile.WriteString(headercontent)
	}
	htmlfile.Close()

	//make header file
	footercontent := "<!doctype html>\n" +
		"<html>\n" +
		"<head>\n" +
		"	<meta charset=\"utf-8\">\n" +
		"	<script>\n" +
		"		function substitutePdfVariables() {\n\n" +
		"			function getParameterByName(name) {\n" +
		"				var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);\n" +
		"				return match && decodeURIComponent(match[1].replace(/\\+/g, ' '));\n" +
		"			}\n\n" +
		"			function substitute(name) {\n" +
		"				var value = getParameterByName(name);\n" +
		"				var elements = document.getElementsByClassName(name);\n\n" +
		"				for (var i = 0; elements && i < elements.length; i++) {\n" +
		"					elements[i].textContent = value;\n" +
		"				}\n" +
		"			}\n\n" +
		"			['frompage', 'topage', 'page', 'webpage', 'section', 'subsection', 'subsubsection']\n" +
		"				.forEach(function(param) {\n" +
		"					substitute(param);\n" +
		"				});\n" +
		"		}\n" +
		"	</script>\n" +
		"  <style>\n" +
		"	p {\n" +
		"	  text-align: right;\n" +
		"	  size: 16pt;\n" +
		"	  font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	}\n" +
		"#footer {\n" +
		"	width: 99%;\n" +
		"	height: 80px;\n" +
		"	white-space: nowrap;\n" +
		"	}\n" +
		"	.alignleft {\n" +
		"	display: inline-block;\n" +
		"	width: 33%;\n" +
		"	height: 80%;\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	padding-top: 10px;\n" +
		"	}\n" +
		"	.aligncenter {\n" +
		"	display: inline-block;\n" +
		"	text-align: center;\n" +
		"	width: 33%;\n" +
		"	height: 80%;\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	}\n" +
		"	.alignright {\n" +
		"	display: inline-block;\n" +
		"	text-align: right;\n" +
		"	width: 33%;\n" +
		"	height: 80%;\n" +
		"	font-family: Arial, Helvetica, sans-serif," + `"Nirmala UI"` + ";\n" +
		"	}\n" +
		"  </style>\n" +
		"</head>\n" +
		"<body onload=\"substitutePdfVariables()\">\n" +
		"	<div id=\"footer\">\n" +
		"		<div class=\"alignleft\">" + time.Now().Format("02-January-2006") + "</p></div>\n" +
		"		<div class=\"aligncenter\">" + pastor_name + "</p></div>\n" +
		"		<div class=\"alignright\"> <span class=\"page\"></span>/<span class=\"topage\"></span></p></div>\n" +
		"</div>\n" +
		"</body>\n" +
		"</html>"

	// save as footer html file
	htmlfile, err = os.OpenFile(paths.Tempdir+"footer.html", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		basic.CheckErr(err, "Error creating footer.html file")
	} else {
		htmlfile.WriteString(footercontent)
	}
	htmlfile.Close()

	docprocess.Max = float64(len(maintext))

	for i := range maintext {

		w4document.SetText(maintext[i].Titel)
		docprocess.SetValue(float64(i) + 1)

		//write bibleverses remove first dublicates
		for ii := range maintext[i].Verse {
			var nodublicate = make([]string, 0, 50)
			inResult := make(map[string]bool)
			// remove dublicates
			for _, str := range maintext[i].Verse[ii] {
				if _, ok := inResult[str]; !ok {
					inResult[str] = true
					nodublicate = append(nodublicate, str)
				}
			}
			// write verses maintext
			var lst_versetext = make([]string, 0, 70)
			for _, s := range nodublicate {

				if len(s) < 90 {
					lst_versetext = append(lst_versetext, s)
				} else {
					var first_prt, second_prt, third_prt, fourth_prt, fifth_prt string
					var sixth_prt, seventh_prt, eighth_prt, ninth_prt, tenth_prt string
					var length int
					split_temp := strings.Split(s, " ")

					//check for to long line and separate it
					for _, elem := range split_temp {
						//add length of element and plus 1 for whitespace
						length += len(elem) + 1
						if length < 90 {
							first_prt = first_prt + elem + " "
						} else if length < 180 {
							second_prt = second_prt + elem + " "
						} else if length < 270 {
							third_prt = third_prt + elem + " "
						} else if length < 360 {
							fourth_prt = fourth_prt + elem + " "
						} else if length < 450 {
							fifth_prt = fifth_prt + elem + " "
						} else if length < 540 {
							sixth_prt = sixth_prt + elem + " "
						} else if length < 630 {
							seventh_prt = seventh_prt + elem + " "
						} else if length < 720 {
							eighth_prt = eighth_prt + elem + " "
						} else if length < 810 {
							ninth_prt = ninth_prt + elem + " "
						} else {
							tenth_prt = tenth_prt + elem + " "
						}
					}

					//put lines in,list
					if len(first_prt) > 0 {
						lst_versetext = append(lst_versetext, first_prt)
					}
					if len(second_prt) > 0 {
						lst_versetext = append(lst_versetext, second_prt)
					}
					if len(third_prt) > 0 {
						lst_versetext = append(lst_versetext, third_prt)
					}
					if len(fourth_prt) > 0 {
						lst_versetext = append(lst_versetext, fourth_prt)
					}
					if len(fifth_prt) > 0 {
						lst_versetext = append(lst_versetext, fifth_prt)
					}
					if len(sixth_prt) > 0 {
						lst_versetext = append(lst_versetext, sixth_prt)
					}
					if len(seventh_prt) > 0 {
						lst_versetext = append(lst_versetext, seventh_prt)
					}
					if len(eighth_prt) > 0 {
						lst_versetext = append(lst_versetext, eighth_prt)
					}
					if len(ninth_prt) > 0 {
						lst_versetext = append(lst_versetext, ninth_prt)
					}
					if len(tenth_prt) > 0 {
						lst_versetext = append(lst_versetext, tenth_prt)
					}
				}
			}

			// write empty line between main verses and translation verses
			for k := 0; k < 1; k++ {
				lst_versetext = append(lst_versetext, "<br>")
			}

			// add translation text
			for _, translation := range lst_transtext {

				for numb, text := range translation {
					if numb > i {
						break
					} else if numb < i {
						continue
					} else if numb == i {
						for _, verse := range text.Verse {
							var nodublicate2 = make([]string, 0, 70)
							inResult := make(map[string]bool)
							// remove dublicates
							for _, str := range verse {
								if _, ok := inResult[str]; !ok {
									inResult[str] = true
									nodublicate2 = append(nodublicate2, str)
								}
							}

							for _, item := range nodublicate2 {

								if len(item) < 90 {
									lst_versetext = append(lst_versetext, item)
								} else {
									var first_prt, second_prt, third_prt, fourth_prt string
									var length int
									split_temp := strings.Split(item, " ")

									//check for to long line and separate it
									for _, elem := range split_temp {
										//add length of element and plus 1 for whitespace
										length += len(elem) + 1
										if length < 90 {
											first_prt = first_prt + elem + " "
										} else if length < 180 {
											second_prt = second_prt + elem + " "
										} else if length < 270 {
											third_prt = third_prt + elem + " "
										} else {
											fourth_prt = fourth_prt + elem + " "
										}
									}
									//put lines in,list
									if len(first_prt) > 0 {
										lst_versetext = append(lst_versetext, first_prt)
									}
									if len(second_prt) > 0 {
										lst_versetext = append(lst_versetext, second_prt)
									}
									if len(third_prt) > 0 {
										lst_versetext = append(lst_versetext, third_prt)
									}
									if len(fourth_prt) > 0 {
										lst_versetext = append(lst_versetext, fourth_prt)
									}
								}
							}
							// write empty line between the verses
							for j := 0; j < 1; j++ {
								lst_versetext = append(lst_versetext, "<br>")
							}
						}
					}
				}
			}

			//write title of bibleverse
			filecontent += "<h3>" + maintext[i].Titel + "</h3>\n"

			//write lines to html
			for _, element := range lst_versetext {

				//check if arabic string
				isarabic := checkifArabic(element)
				fmt.Println(33, isarabic)

				if isarabic {
					filecontent += "<p><span dir=\"rtl\">" + element + "</span></p>\n"
				} else {
					filecontent += "<p>" + element + "</p>\n"
				}

			}

			// write empty space between the verses
			filecontent += "<br>"
		}
	}

	filecontent += "</body>\n</html>"

	// save as html file
	htmlfile, err = os.OpenFile(paths.Outputpath+paths.Pathseperator+"html"+paths.Pathseperator+filename+".html", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		basic.CheckErr(err, "Error creating "+filename+" .html file")
	} else {
		htmlfile.WriteString(filecontent)
	}
	htmlfile.Close()

	var wkhtmltopdfstring string
	switch runtime.GOOS {
	case "windows":
		wkhtmltopdfstring = "wkhtmltopdf.exe"
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			wkhtmltopdfstring = "wkhtmltopdfdarwinarm64"
		case "amd64":
			wkhtmltopdfstring = "wkhtmltopdfdarwinamd64"
		}

	}
	//create pdf
	basic.StartcommandwithArgs(paths.Currentdirectory+wkhtmltopdfstring, []string{"--header-html", paths.Tempdir + "header.html", "--footer-html", paths.Tempdir + "footer.html", paths.Outputpath + paths.Pathseperator + "html" + paths.Pathseperator + filename + ".html", paths.Outputpath + paths.Pathseperator + filename + ".pdf"})

}

func checkifArabic(input string) bool {
	//check if string has at least 5 characters in arabic then count it as true
	var count = 0
	var isArabic = false

	for _, v := range input {
		if unicode.In(v, unicode.Arabic) {
			isArabic = true
			count += 1
		} else {
			isArabic = false
		}
	}

	isArabic = count >= 5

	return isArabic
}

func CombinedPDF(paths *basic.OSPaths) {

	// get files in directory
	files, err := ioutil.ReadDir(paths.Outputpath)
	if err != nil {
		log.Fatal(err)
	}

	//get list of all pdf's in output folder
	var pdflist = make([]string, 0, 40)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".pdf" {
			pdflist = append(pdflist, paths.Outputpath+paths.Pathseperator+file.Name())
		}
	}

	// merge them in one file
	pdfmerge.Pdfmerge(pdflist, paths.Outputpath+"/AllTranslation.pdf")
}
