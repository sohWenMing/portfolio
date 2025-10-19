package templates

import "embed"

//go:embed *.gohtml
var TemplateFS embed.FS
