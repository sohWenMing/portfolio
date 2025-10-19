package handlers

import (
	"fmt"
	"net/http"

	"github.com/sohWenMing/portfolio/internal/contextkeys"
	"github.com/sohWenMing/portfolio/internal/views/templating"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	csrfToken := r.Context().Value(contextkeys.CSRFTokenKey)
	//TODO: replace println with actual logging function
	if csrfToken == "" {
		fmt.Println("csrfToken could not be found")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	} else {
		fmt.Println("token in TestHandler: ", csrfToken)
	}
	tpl, err := templating.GetTestHTMLTemplate()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, nil)
}
