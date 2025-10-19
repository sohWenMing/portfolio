package handlers

import (
	"net/http"

	"github.com/sohWenMing/portfolio/internal/views/templating"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := templating.GetTestHTMLTemplate()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, nil)
}
