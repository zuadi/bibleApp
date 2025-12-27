package bibletool

import (
	"bibletool/bibletool/env"
	"bibletool/bibletool/models"
	"bibletool/utils"

	templatebuilder "gitea.tecamino.com/paadi/template-builder"
)

func (bt *Bibletool) WriteHtml(outputPath string, data models.HtmlStruct) error {
	bt.DebugLog("WriteHtml", "start template builder")
	tmplBuilder := templatebuilder.NewTemplateBuilder()

	path := utils.GetDistOsPath(env.HtmlTemplateFile.GetValue())
	err := tmplBuilder.Generate(path, outputPath, &data)
	if err != nil {
		return err
	}
	return nil
}
