package bibletool

import (
	"bibletool/bibletool/consts"
	"bibletool/bibletool/models"

	templatebuilder "gitea.tecamino.com/paadi/template-builder"
)

func (bt *Bibletool) WriteHtml(outputPath string, data models.HtmlStruct) error {
	tmplBuilder := templatebuilder.NewTemplateBuilder()
	err := tmplBuilder.Generate(consts.HtmlTemplatePath, outputPath, &data)
	if err != nil {
		return err
	}
	return nil
}
