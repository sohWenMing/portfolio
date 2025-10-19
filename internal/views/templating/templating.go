package templating

import (
	"fmt"
	"html/template"
	"io/fs"

	"github.com/sohWenMing/portfolio/web/templates"
)

func GetTestHTMLTemplate() (*template.Template, error) {
	embeddedFS := templates.TemplateFS

	tpl, err := template.ParseFS(embeddedFS, "main_template.gohtml")
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func ReadFile() {
	embeddedTemplate := templates.TemplateFS
	readFile, _ := fs.ReadFile(embeddedTemplate, "main_template.gohtml")
	fmt.Println("readFile: ", string(readFile))
}
