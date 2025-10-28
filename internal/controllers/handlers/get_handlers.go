package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
)

type TemplateExecutor interface {
	ExecuteTemplateWithCSRF(w http.ResponseWriter, r *http.Request, csrfField template.HTML, baseTemplate string, data any)
}

func TestHandler(tplExecutor TemplateExecutor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templateField := csrf.TemplateField(r)
		fmt.Println("templateField: ", templateField)
		// csrfToken := r.Context().Value(contextkeys.CSRFTokenKey)
		// csrfField := fmt.Sprintf("<input type=hidden name=csrfField value=%s", csrfToken)
		// csrfFieldHTML := template.HTML(csrfField)
		//TODO: replace println with actual logging function

		// if csrfToken == "" {
		// 	fmt.Println("csrfToken could not be found")
		// 	http.Error(w, "Internal Error", http.StatusInternalServerError)
		// 	return
		// } else {
		// 	fmt.Println("token in TestHandler: ", csrfToken)
		// }
		tplExecutor.ExecuteTemplateWithCSRF(w, r, csrf.TemplateField(r), "create_user_form.gohtml", nil)
	}
}

func TestPing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ping successful"))
}
