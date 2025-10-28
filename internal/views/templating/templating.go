package templating

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/sohWenMing/portfolio/web/templates"
)

var templatesToParse []string = []string{
	"create_user_form.gohtml",
}

type TemplateExecutor struct {
	template *template.Template
}

func (t *TemplateExecutor) ExecuteTemplateWithCSRF(
	w http.ResponseWriter,
	r *http.Request,
	csrfField template.HTML,
	baseTemplate string,
	data any) {
	cloned, err := t.template.Clone()
	if err != nil {
		//TODO: replace println with actual logging function
		fmt.Println("error occured when cloning: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	cloned = cloned.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrfField
			},
		},
	)
	cloned.ExecuteTemplate(w, baseTemplate, data)
}

func InitTemplateExecutor() (tplExecutor *TemplateExecutor, err error) {
	baseTemplateWithFuncs := template.New("base")
	/*
		Attach required functions to the template
		Note that the csrfField function needs to return a placeholder hidden input, because at load time the
		value for csrfField cannot yet be evaluated.
	*/
	baseTemplateWithFuncs = baseTemplateWithFuncs.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return `<input type=hidden />`
			},
		},
	)
	embeddedFS := templates.TemplateFS

	tpl, err := baseTemplateWithFuncs.ParseFS(embeddedFS, templatesToParse...)
	if err != nil {
		return nil, err
	}
	return &TemplateExecutor{
		template: tpl,
	}, nil
}

/*
Used to initially atttempt to load and parse all the templates at the beginning of starting up the server.
If an error occurs, will panic and shut down program.
*/
func InitTemplates() *template.Template {
	// start with a base empty template using template.New()
	baseTemplateWithFuncs := template.New("base")
	/*
		Attach required functions to the template
		Note that the csrfField function needs to return a placeholder hidden input, because at load time the
		value for csrfField cannot yet be evaluated.
	*/
	baseTemplateWithFuncs = baseTemplateWithFuncs.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return `<input type=hidden />`
			},
		},
	)
	embeddedFS := templates.TemplateFS

	tpl, err := baseTemplateWithFuncs.ParseFS(embeddedFS, templatesToParse...)
	return template.Must(tpl, err)
}
