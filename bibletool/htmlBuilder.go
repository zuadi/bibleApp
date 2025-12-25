package bibletool

import (
	"bibletool/bibletool/env"
	"bibletool/bibletool/models"

	templatebuilder "gitea.tecamino.com/paadi/template-builder"
)

func (bt *Bibletool) WriteHtml(outputPath string, data models.HtmlStruct) error {
	tmplBuilder := templatebuilder.NewTemplateBuilder()
	err := tmplBuilder.Generate(env.HtmlTemplateFile.GetValue(), outputPath, &data)
	if err != nil {
		return err
	}
	return nil
}
