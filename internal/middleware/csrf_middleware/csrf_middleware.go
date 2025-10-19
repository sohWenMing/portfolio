package csrfmiddleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/sohWenMing/portfolio/internal/contextkeys"
)

/*
Middleware gets token which should have been attached to protected routes. Token to be retrieved and passed on to
next handler, retrieve by using ctx.Value(GetCsrfTokenKey())
*/

func CSRFMWGetToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := csrf.Token(r)
		if token == "" {
			//TODO: change fmt.Println to actual logging func
			fmt.Println("csrf token could not be retrieved")
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		} else {
			fmt.Println("csrf token - ", token)
		}
		currCtx := r.Context()
		key := contextkeys.CSRFTokenKey
		ctx := context.WithValue(currCtx, key, token)
		r = r.WithContext(ctx)
		fmt.Println("token: ", token)
		next.ServeHTTP(w, r)
	})
}
